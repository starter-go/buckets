package main

import (
	"os"

	"github.com/starter-go/buckets/modules/buckets"
	"github.com/starter-go/units"
)

func main() {

	a := os.Args
	m := buckets.ModuleTest()

	ctx := &units.Context{
		Arguments: a,
		Module:    m,
		UsePanic:  true,
	}

	units.Run(ctx)

}
