# waitfor
[![Go Report Card](https://goreportcard.com/badge/github.com/davejohnston/waitfor)](https://goreportcard.com/report/github.com/davejohnston/waitfor)
[![GoDoc](https://godoc.org/github.com/davejohnston/waitfor?status.svg)](https://godoc.org/github.com/davejohnston/waitfor)

WaitFor is a simple library for implementing dependency checks in your Go application.

## Summary
Waitfor provides a simple API to wait for your applications dependancies to become available. For example if your application relies on an external database wait for can check that the service is up and listening.
This can be useful in container environments where your services DNS name is not registered until the remote service becomes healthy.

## Usage

See the [GoDoc examples](https://godoc.org/github.com/davejohnston/waitfor) for more detail.

 - Install with `go get`: `go get -u github.com/davejohnston/waitfor`
 
  - Create a dependancy handler:
   ```go
   var dependancies = waitfor.NewDependencies()()
   ```
  - Add a dependancy:
  ```go
  dependancies.Add("rest-api", waitfor.ServiceListening("0.0.0.0:9999", 10*time.Second))
  dependancies.Add("postgres-db", waitfor.DatabaseReady("postgres", "host=localhost port=5423 user=postgres password=secret dbname=postgres sslmode=disable")
  ```
  - Wait for the dependancies:
  ```go
  if err := dependancies.Wait(); err != nil {
      log.fatal(err)
  }
  ```

