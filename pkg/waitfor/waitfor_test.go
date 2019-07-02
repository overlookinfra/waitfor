package waitfor_test

import (
	"log"
	"net"
	"testing"
	"time"

	"github.com/davejohnston/waitfor/pkg/waitfor"
	"github.com/stretchr/testify/assert"
)

func TestNewDependencies(t *testing.T) {
	dependencies := waitfor.NewDependencies()
	assert.NotNil(t, dependencies)
}

func TestAddDependenciesWithError(t *testing.T) {
	dependencies := waitfor.NewDependencies()
	dependencies.Add("rest-api", waitfor.ServiceListening("0.0.0.0:1234", 5*time.Second))
	err := dependencies.Wait(waitfor.Options{1, 1})
	assert.Error(t, err)
}

func TestAddDependencies(t *testing.T) {

	l, err := net.Listen("tcp", "0.0.0.0:9999")
	defer l.Close()
	assert.NoError(t, err)

	dependencies := waitfor.NewDependencies()
	dependencies.Add("rest-api", waitfor.ServiceListening("0.0.0.0:9999", 5*time.Second))
	err = dependencies.Wait(waitfor.Options{1, 1})
	assert.NoError(t, err)
}

func ExampleNewDependencies() {
	dependencies := waitfor.NewDependencies()
	dependencies.Add("rest-api", waitfor.ServiceListening("0.0.0.0:9999", 10*time.Second))
	err := dependencies.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
