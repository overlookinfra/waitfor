package waitfor_test

import (
	"log"
	"testing"
	"time"

	"github.com/davejohnston/waitfor/pkg/waitfor"
	"github.com/stretchr/testify/assert"
)

func NewDependencies(t *testing.T) {
	dependencies := waitfor.NewDependencies()
	assert.NotNil(t, dependencies)
}

func ExampleNewDependencies() {
	dependencies := waitfor.NewDependencies()
	dependencies.Add("rest-api", waitfor.ServiceListening("0.0.0.0:9999", 10*time.Second))
	err := dependencies.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
