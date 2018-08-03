package mgosrv

import (
	"log"
	"testing"

	"github.com/jamillosantos/macchiato"
	"github.com/lab259/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/mgo.v2"
	"gopkg.in/yaml.v2"
)

func TestService(t *testing.T) {
	log.SetOutput(GinkgoWriter)
	RegisterFailHandler(Fail)
	macchiato.RunSpecs(t, "Mgo Test Suite")
}

func pingSession(session *mgo.Session) error {
	return session.Ping()
}

var _ = Describe("MgoService", func() {
	It("should fail loading a configuration", func() {
		var service MgoService
		configuration, err := service.LoadConfiguration()
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("not implemented"))
		Expect(configuration).To(BeNil())
	})

	It("should fail applying configuration", func() {
		var service MgoService
		err := service.ApplyConfiguration(map[string]interface{}{
			"address": "localhost",
		})
		Expect(err).To(Equal(http.ErrWrongConfigurationInformed))
	})

	It("should apply the configuration using a pointer", func() {
		var service MgoService
		err := service.ApplyConfiguration(&MgoServiceConfiguration{
			Addresses: []string{"addresses"},
			Username:  "username",
			Password:  "password",
			Database:  "database",
			PoolSize:  1,
			Timeout:   60,
		})
		Expect(err).To(BeNil())
		Expect(service.Configuration.Addresses).To(Equal([]string{"addresses"}))
		Expect(service.Configuration.Username).To(Equal("username"))
		Expect(service.Configuration.Password).To(Equal("password"))
		Expect(service.Configuration.Database).To(Equal("database"))
		Expect(service.Configuration.PoolSize).To(Equal(1))
		Expect(service.Configuration.Timeout).To(Equal(60))
	})

	It("should apply the configuration using a copy", func() {
		var service MgoService
		err := service.ApplyConfiguration(MgoServiceConfiguration{
			Addresses: []string{"addresses"},
			Username:  "username",
			Password:  "password",
			Database:  "database",
			PoolSize:  1,
			Timeout:   60,
		})
		Expect(err).To(BeNil())
		Expect(service.Configuration.Addresses).To(Equal([]string{"addresses"}))
		Expect(service.Configuration.Username).To(Equal("username"))
		Expect(service.Configuration.Password).To(Equal("password"))
		Expect(service.Configuration.Database).To(Equal("database"))
		Expect(service.Configuration.PoolSize).To(Equal(1))
		Expect(service.Configuration.Timeout).To(Equal(60))
	})

	It("should start the service", func() {
		var service MgoService
		Expect(service.ApplyConfiguration(MgoServiceConfiguration{
			Addresses: []string{"localhost"},
			Username:  "",
			Password:  "",
			Database:  "test-service-database",
			PoolSize:  1,
			Timeout:   60,
		})).To(BeNil())
		Expect(service.Start()).To(BeNil())
		defer service.Stop()
		Expect(service.RunWithSession(pingSession)).To(BeNil())
	})

	It("should stop the service", func() {
		var service MgoService
		Expect(service.ApplyConfiguration(MgoServiceConfiguration{
			Addresses: []string{"localhost"},
			Username:  "",
			Password:  "",
			Database:  "test-service-database",
			PoolSize:  1,
			Timeout:   60,
		})).To(BeNil())
		Expect(service.Start()).To(BeNil())
		Expect(service.Stop()).To(BeNil())
		Expect(service.RunWithSession(func(session *mgo.Session) error {
			return nil
		})).To(Equal(http.ErrServiceNotRunning))
	})

	It("should restart the service", func() {
		var service MgoService
		Expect(service.ApplyConfiguration(MgoServiceConfiguration{
			Addresses: []string{"localhost"},
			Username:  "",
			Password:  "",
			Database:  "test-service-database",
			PoolSize:  1,
			Timeout:   60,
		})).To(BeNil())
		Expect(service.Start()).To(BeNil())
		Expect(service.Restart()).To(BeNil())
		Expect(service.RunWithSession(pingSession)).To(BeNil())
	})

	Describe("Mode YAML", func() {
		It("should fill with default mode", func() {
			var service MgoService
			var configuration MgoServiceConfiguration

			var yamlData = `
addresses:
  - "localhost"
database: "database"
username: ""
password: ""
timeout: 60
`

			err := yaml.Unmarshal([]byte(yamlData), &configuration)

			Expect(err).To(BeNil())

			err = service.ApplyConfiguration(configuration)
			Expect(err).To(BeNil())
			Expect(*service.Configuration.Mode).To(Equal(MgoServiceMode(mgo.Primary)))
		})

		It("should fail if mode is a wrong value", func() {
			var configuration MgoServiceConfiguration

			var yamlData = `
addresses:
  - "localhost"
database: "database"
username: ""
password: ""
timeout: 60
mode: "wrong_mode"
`

			err := yaml.Unmarshal([]byte(yamlData), &configuration)

			Expect(err.Error()).To(ContainSubstring("invalid mode"))
		})

		It("should pass changing the mode", func() {
			var service MgoService
			var configuration MgoServiceConfiguration

			var yamlData = `
addresses:
  - "localhost"
database: "database"
username: ""
password: ""
timeout: 60
mode: "nearest"
`

			err := yaml.Unmarshal([]byte(yamlData), &configuration)

			Expect(err).To(BeNil())

			err = service.ApplyConfiguration(configuration)
			Expect(err).To(BeNil())
			Expect(*service.Configuration.Mode).To(Equal(MgoServiceMode(mgo.Nearest)))
		})
	})
})
