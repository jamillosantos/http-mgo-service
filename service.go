package http_mgo_service

import (
	"errors"
	"time"
	"gopkg.in/mgo.v2"
	"github.com/jamillosantos/http"
	"fmt"
)

type MgoServiceConfiguration struct {
	Addresses []string       `yaml:"addresses"`
	Database  string         `yaml:"database"`
	Username  string         `yaml:"username"`
	Password  string         `yaml:"password"`
	PoolSize  int            `yaml:"pool_size"`
	Timeout   int            `yaml:"timeout"`
	Mode      *MgoServiceMode `yaml:"mode"`
}

type MgoServiceMode mgo.Mode

func (mode *MgoServiceMode) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var modeStr string

	err := unmarshal(&modeStr)
	if err != nil {
		return err
	}

	switch modeStr {
	case "primary":
		*mode = MgoServiceMode(mgo.Primary)
		break
	case "primarypreferred":
		*mode = MgoServiceMode(mgo.PrimaryPreferred)
		break
	case "secondary":
		*mode = MgoServiceMode(mgo.Secondary)
		break
	case "secondarypreferred":
		*mode = MgoServiceMode(mgo.SecondaryPreferred)
		break
	case "nearest":
		*mode = MgoServiceMode(mgo.Nearest)
		break
	default:
		return fmt.Errorf("invalid mode: %s", modeStr)
	}
	return nil
}

type MgoService struct {
	running       bool
	session       *mgo.Session
	Configuration MgoServiceConfiguration
}
type MgoServiceSessionHandler func(session *mgo.Session) error

func (service *MgoService) LoadConfiguration() (interface{}, error) {
	return nil, errors.New("not implemented")
}

func (service *MgoService) ApplyConfiguration(configuration interface{}) error {
	switch c := configuration.(type) {
	case MgoServiceConfiguration:
		service.Configuration = c
	case *MgoServiceConfiguration:
		service.Configuration = *c
	default:
		return http.ErrWrongConfigurationInformed
	}

	if service.Configuration.Mode == nil {
		defaultMode := MgoServiceMode(mgo.Primary)
		service.Configuration.Mode = &defaultMode
	}

	return nil
}

func (service *MgoService) Restart() error {
	if service.running {
		err := service.Stop()
		if err != nil {
			return err
		}
	}
	return service.Start()
}

func (service *MgoService) Start() error {
	if !service.running {
		var err error

		if service.Configuration.Username == "" && service.Configuration.Password == "" {
			service.session, err = mgo.Dial(fmt.Sprintf("%s/%s", service.Configuration.Addresses[0], service.Configuration.Database))
		} else {
			dialInfo := &mgo.DialInfo{
				Addrs:    service.Configuration.Addresses,
				Timeout:  time.Duration(service.Configuration.Timeout) * time.Second,
				Database: service.Configuration.Database,
				Username: service.Configuration.Username,
				Password: service.Configuration.Password,
			}
			service.session, err = mgo.DialWithInfo(dialInfo)
		}

		if err != nil {
			return err
		}

		service.session.SetMode(mgo.Mode(*service.Configuration.Mode), true)
		if service.Configuration.PoolSize > 0 {
			service.session.SetPoolLimit(service.Configuration.PoolSize)
		}

		_, err = service.session.DB("").CollectionNames()
		if err != nil {
			return err
		}

		service.running = true
	}
	return nil
}

func (service *MgoService) newSession() *mgo.Session {
	return service.session.Clone()
}

func (service *MgoService) Stop() error {
	if service.running {
		service.session.Close()
		service.running = false
	}
	return nil
}

func (service *MgoService) RunWithSession(handler MgoServiceSessionHandler) error {
	if !service.running {
		return http.ErrServiceNotRunning
	}
	newSession := service.newSession()
	defer newSession.Close()
	return handler(newSession)
}
