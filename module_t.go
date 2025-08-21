package buckets

import (
	"embed"

	"github.com/starter-go/application"
	"github.com/starter-go/starter"
)

const (
	theModuleName     = "github.com/starter-go/buckets"
	theModuleVersion  = "0.0.0"
	theModuleRevision = 0
)

////////////////////////////////////////////////////////////////////////////////

const (
	theTestModuleResPath = "src/test/resources"
	theMainModuleResPath = "src/main/resources"
)

//go:embed "src/test/resources"
var theTestModuleResFS embed.FS

//go:embed "src/main/resources"
var theMainModuleResFS embed.FS

////////////////////////////////////////////////////////////////////////////////

func NewModuleLib() *application.ModuleBuilder {
	builder := new(application.ModuleBuilder)
	builder.Name(theModuleName + "#lib")
	builder.Version(theModuleVersion)
	builder.Revision(theModuleRevision)

	builder.EmbedResources(theMainModuleResFS, theMainModuleResPath)

	builder.Depend(starter.Module())

	return builder
}

func NewModuleTest() *application.ModuleBuilder {
	builder := new(application.ModuleBuilder)
	builder.Name(theModuleName + "#test")
	builder.Version(theModuleVersion)
	builder.Revision(theModuleRevision)

	builder.EmbedResources(theTestModuleResFS, theTestModuleResPath)

	builder.Depend(starter.Module())

	return builder
}
