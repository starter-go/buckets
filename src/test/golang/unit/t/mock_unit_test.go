package t

import (
	"testing"

	"github.com/starter-go/buckets/modules/buckets"
	"github.com/starter-go/buckets/src/test/golang/unit"
	"github.com/starter-go/units"
)

func TestUseMock(t *testing.T) {

	args := []string{}

	props := map[string]string{
		"debug.enabled":        "1",
		"debug.log-properties": "1",
	}

	units.Run(&units.Config{
		Args:       args,
		Cases:      unit.TheMockUnit,
		Module:     buckets.ModuleTest(),
		T:          t,
		Properties: props,
		UsePanic:   false,
	})

}
