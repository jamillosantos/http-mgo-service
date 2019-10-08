package mgosrv

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/globalsign/mgo"
	"github.com/lab259/go-rscsrv"
)

// MgoServiceConfiguration describes the `MgoService` configuration.
type MgoServiceConfiguration struct {
	Addresses []string        `yaml:"addresses"`
	Database  string          `yaml:"database"`
	Username  string          `yaml:"username"`
	Password  string          `yaml:"password"`
	PoolSize  int             `yaml:"pool_size"`
	Timeout   int             `yaml:"timeout"`
	Mode      *MgoServiceMode `yaml:"mode"`
	UseTLS    bool            `yaml:"use_tls"`
}

// MgoServiceMode is an alias for the `mgo.Mode` that implements
// Unmarshaling from the YAML.
type MgoServiceMode mgo.Mode

// UnmarshalYAML implements the marshaling a `mgo.Mode` to string.
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

// MgoService implements the mgo service itself.
type MgoService struct {
	serviceState
	session       *mgo.Session
	Configuration MgoServiceConfiguration
}

// LoadConfiguration is an abstract method that should be overwritten on the
// actual usage of this service.
func (service *MgoService) LoadConfiguration() (interface{}, error) {
	return nil, errors.New("not implemented")
}

// ApplyConfiguration implements the type verification of the given
// `configuration` and applies it to the service.
func (service *MgoService) ApplyConfiguration(configuration interface{}) error {
	switch c := configuration.(type) {
	case MgoServiceConfiguration:
		service.Configuration = c
	case *MgoServiceConfiguration:
		service.Configuration = *c
	default:
		return rscsrv.ErrWrongConfigurationInformed
	}

	// If the configuration mode is not defined, get the default behavior.
	// From the MongoDB documentation.
	if service.Configuration.Mode == nil {
		defaultMode := MgoServiceMode(mgo.Primary)
		service.Configuration.Mode = &defaultMode
	}

	return nil
}

// Restart restarts the service.
func (service *MgoService) Restart() error {
	if service.isRunning() {
		err := service.Stop()
		if err != nil {
			return err
		}
	}
	return service.Start()
}

// Start initialize the mongo connection and saves the session.
func (service *MgoService) Start() error {
	if !service.isRunning() {
		var err error

		dialInfo := &mgo.DialInfo{
			Addrs:    service.Configuration.Addresses,
			Timeout:  time.Duration(service.Configuration.Timeout) * time.Second,
			Database: service.Configuration.Database,
			Username: service.Configuration.Username,
			Password: service.Configuration.Password,
		}

		if service.Configuration.UseTLS {
			dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
				tlsConfig := &tls.Config{}
				return tls.Dial("tcp", addr.String(), tlsConfig)
			}
		}

		service.session, err = mgo.DialWithInfo(dialInfo)
		if err != nil {
			return err
		}

		// Applies the mode
		service.session.SetMode(mgo.Mode(*service.Configuration.Mode), true)

		// Sets the pool size
		if service.Configuration.PoolSize > 0 {
			service.session.SetPoolLimit(service.Configuration.PoolSize)
		}

		// Pings the session to ensure it is working
		err = service.session.Ping()
		if err != nil {
			return err
		}

		service.setRunning(true)
	}
	return nil
}

// newSession suppose to pull a new session instance from the pool.
func (service *MgoService) newSession() *mgo.Session {
	return service.session.Clone()
}

// Stop stops the service.
func (service *MgoService) Stop() error {
	if service.isRunning() {
		service.session.Close()
		service.setRunning(false)
	}
	return nil
}

// RunWithSession runs a handler passing a new instance of the a session.
func (service *MgoService) RunWithSession(handler func(session *mgo.Session) error) error {
	if service.isRunning() {
		newSession := service.newSession()
		defer newSession.Close()
		return handler(newSession)
	}
	return rscsrv.ErrServiceNotRunning
}
