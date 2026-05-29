package t

import (
	"testing"

	"github.com/starter-go/buckets/modules/buckets"
	"github.com/starter-go/buckets/src/test/golang/unit"
	"github.com/starter-go/units"
)

func TestUseMock(t *testing.T) {

	args := []string{}
	ctx := units.NewContext()
	props := map[string]string{
		"debug.enabled":        "1",
		"debug.log-properties": "1",
	}

	ctx.Arguments = args
	ctx.Module = buckets.ModuleTest()
	ctx.Properties = props
	ctx.UsePanic = true
	ctx.Selector = "#" + unit.TheMockUnit
	ctx.T = t

	units.Run(ctx)

}
