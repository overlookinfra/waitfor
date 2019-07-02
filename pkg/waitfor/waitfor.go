// Copyright 2019 by the contributors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package waitfor

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	defaultRetries = 10               // defaultRetries specifies how many times to retry a check
	defaultTimeout = 10 * time.Second // defaultTimeout specifies the time to sleep between retries.
)

// Check is function that determines if a dependency is available
// it returns an error if the dependency is not ready.
type Check func() error

// Options defines the retry and timeout parameters when
// performing checks
type Options struct {
	Retries int
	Timeout time.Duration
}

// Handler is a dependency Handler.  It maintains a list
// of dependencies to check
type Handler struct {
	checksMutex  sync.RWMutex
	dependencies map[string]Check
}

// Add a new dependency check
// e.g
// dependancies = waitfor.NewDependencies()
// dependancies.Add("rest-endpoint", waitfor.ServiceListening("0.0.0.0:8443", 10*time.Second))
func (h *Handler) Add(name string, waitFor Check) {
	h.checksMutex.Lock()
	defer h.checksMutex.Unlock()
	h.dependencies[name] = waitFor
}

// Wait executes the check for all dependencies
// If any dependency fails after the given number of retries and timeouts it will return error.
// All dependencies are checked in parallel.
func (h *Handler) Wait(options ...Options) error {

	option := &Options{
		defaultRetries,
		defaultTimeout,
	}

	if len(options) >= 1 {
		option = &options[0]
	}

	var wg sync.WaitGroup
	errorMessages := make(chan error, len(h.dependencies)+1)

	for name, check := range h.dependencies {
		wg.Add(1)
		go func(n string, c Check) {
			log.Printf("Waiting for %s\n", n)
			if err := performCheck(n, c, *option); err != nil {
				errorMessages <- err
			} else {
				logrus.Debugf("%s is ready\n", n)
			}
			wg.Done()
		}(name, check)
	}

	wg.Wait()
	close(errorMessages)

	return toError(errorMessages)
}

func toError(messages chan error) (err error) {
	var errMsg string
	for e := range messages {
		errMsg += fmt.Sprintf("%s\n", e.Error())
	}
	return fmt.Errorf(errMsg)
}

// NewDependencies returns a dependancy handler.
// This can be used to add multiple dependencies for an application.
// Calling wait on the handler will cause the handler to wait until all dependencies are ready, or until
// it timesout.
func NewDependencies() *Handler {
	h := &Handler{
		dependencies: make(map[string]Check),
	}
	return h
}

// ServiceListening verifies if the endpoint specified by addr
// is listening.  The address string may contain a port e.g.
// 192.168.0.1:1234
func ServiceListening(addr string, timeout time.Duration) Check {
	return func() error {
		conn, err := net.DialTimeout("tcp", addr, timeout)
		if err != nil {
			return err
		}
		return conn.Close()
	}
}

// DatabaseReady determines if a database can be connected to.
// It takes a driver such as mysql or postgres-db and a datasource string.  e.g.
// to connect to a postgres database call the method with these parameters
// DatabaseReady("postgres", "host=localhost port=5432 user=postgres password=secret dbname=postgres sslmode=disable")
func DatabaseReady(driver, datasource string) Check {
	return func() error {
		conn, err := sql.Open(driver, datasource)
		if err != nil {
			return err
		}

		err = conn.Ping()
		if err != nil {
			return err
		}

		return conn.Close()
	}

}

// performCheck executes the check function.   It will retry the number of times
// specified by the options.   After each check it will sleep for the timeout interval.
func performCheck(name string, check Check, option Options) error {
	for i := 0; i < option.Retries; i++ {
		err := check()
		if err == nil {
			break
		}
		if i+1 == option.Retries {
			return fmt.Errorf("Timeout waiting for %s because [%s]", name, err)
		}
		time.Sleep(option.Timeout)
	}

	return nil
}
