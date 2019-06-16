[![CircleCI](https://circleci.com/gh/lab259/go-rscsrv-mgo.svg?style=shield)](https://circleci.com/gh/lab259/go-rscsrv-mgo) [![codecov](https://codecov.io/gh/lab259/go-rscsrv-mgo/branch/master/graph/badge.svg)](https://codecov.io/gh/lab259/go-rscsrv-mgo)

# go-rscsrv-mgo

The go-rscsrv-mgo is the MongoDB resource service that wraps [globalsign/mgo](//github.com/globalsign/mgo) library.

## Dependencies

It depends on the [lab259/go-rscsrv](//github.com/lab259/go-rscsrv) (and its dependencies,
of course) itself and the [globalsign/mgo](//github.com/globalsign/mgo) library.

## Installation

We recommend you to use `dep`.

Or, you may fetch the library directly using `go get`.

```bash
go get github.com/lab259/go-rscsrv-mgo
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
