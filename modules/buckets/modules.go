package buckets

import (
	"github.com/starter-go/application"
	"github.com/starter-go/buckets"
	"github.com/starter-go/buckets/gen/main4buckets"
	"github.com/starter-go/buckets/gen/test4buckets"
	"github.com/starter-go/mimetypes/modules/mimetypes"
	"github.com/starter-go/units/modules/units"
	"github.com/starter-go/v0/libvlog"
)

func ModuleLib() application.Module {
	mb := buckets.NewModuleLib()

	mb.Components(main4buckets.ExportComponents)

	mb.Depend(mimetypes.Module())
	mb.Depend(libvlog.Module())

	return mb.Create()
}

func ModuleTest() application.Module {
	mb := buckets.NewModuleTest()

	mb.Components(test4buckets.ExportComponents)

	mb.Depend(ModuleLib())
	mb.Depend(units.Module())

	return mb.Create()
}
