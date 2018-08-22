[![CircleCI](https://circleci.com/gh/lab259/http-mgo-service.svg?style=shield)](https://circleci.com/gh/lab259/http-mgo-service) [![codecov](https://codecov.io/gh/lab259/http-mgo-service/branch/master/graph/badge.svg)](https://codecov.io/gh/lab259/http-mgo-service)

# http-mgo-service

The http-mgo-service is the [lab259/http](//github.com/lab259/http) service
implementation for the [go-mgo/mgo](//github.com/go-mgo/mgo) library.

## Dependencies

It depends on the [lab259/http](//github.com/lab259/http) (and its dependencies,
of course) itself and the [go-mgo/mgo](//github.com/go-mgo/mgo) library.

## Installation

First, fetch the library to the repository.

```bash
go get github.com/lab259/http-mgo-service
```

## Usage

Applying configuration and starting service

```go
// Create MgoService instance
var mgoService MgoService

// Applying configuration
err := mgoService.ApplyConfiguration(MgoServiceConfiguration{
    Addresses: []string{"localhost"},
	Username:  "username",
	Password:  "password",
	Database:  "my-db",
	PoolSize:  1,
    Timeout:   60,
})

if err != nil {
    panic(err)
}
        
// Starting service
err := mgoService.Start()

if err != nil {
    panic(err)
}

// Create a custom object
var object MyModel

// Executing something using a *mgo.Session
err := mgoService.RunWithSession(func(session *mgo.Session) error {
    // Retrieving an object from the MongoDB
    return session.DB("my-db").C("my-collection").FindId("my-object-id").One(&object)
})

if err != nil {
    panic(err)
}

id := object.Id // "my-object-id"
```

