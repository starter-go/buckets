package test4buckets
import (
    p262c04a06 "github.com/starter-go/buckets"
    pfde5cc127 "github.com/starter-go/buckets/src/test/golang/unit/com"
     "github.com/starter-go/application"
)

// type pfde5cc127.LocalFileUnit in package:github.com/starter-go/buckets/src/test/golang/unit/com
//
// id:com-fde5cc127c251d8d-com-LocalFileUnit
// class:class-0dc072ed44b3563882bff4e657a52e62-Units
// alias:
// scope:singleton
//
type pfde5cc127c_com_LocalFileUnit struct {
}

func (inst* pfde5cc127c_com_LocalFileUnit) register(cr application.ComponentRegistry) error {
	r := cr.NewRegistration()
	r.ID = "com-fde5cc127c251d8d-com-LocalFileUnit"
	r.Classes = "class-0dc072ed44b3563882bff4e657a52e62-Units"
	r.Aliases = ""
	r.Scope = "singleton"
	r.NewFunc = inst.new
	r.InjectFunc = inst.inject
	return r.Commit()
}

func (inst* pfde5cc127c_com_LocalFileUnit) new() any {
    return &pfde5cc127.LocalFileUnit{}
}

func (inst* pfde5cc127c_com_LocalFileUnit) inject(injext application.InjectionExt, instance any) error {
	ie := injext
	com := instance.(*pfde5cc127.LocalFileUnit)
	nop(ie, com)

	
    com.Service = inst.getService(ie)


    return nil
}


func (inst*pfde5cc127c_com_LocalFileUnit) getService(ie application.InjectionExt)p262c04a06.Service{
    return ie.GetComponent("#alias-262c04a06c32904104382e2b8d56c279-Service").(p262c04a06.Service)
}



// type pfde5cc127.MockUnit in package:github.com/starter-go/buckets/src/test/golang/unit/com
//
// id:com-fde5cc127c251d8d-com-MockUnit
// class:class-0dc072ed44b3563882bff4e657a52e62-Units
// alias:
// scope:singleton
//
type pfde5cc127c_com_MockUnit struct {
}

func (inst* pfde5cc127c_com_MockUnit) register(cr application.ComponentRegistry) error {
	r := cr.NewRegistration()
	r.ID = "com-fde5cc127c251d8d-com-MockUnit"
	r.Classes = "class-0dc072ed44b3563882bff4e657a52e62-Units"
	r.Aliases = ""
	r.Scope = "singleton"
	r.NewFunc = inst.new
	r.InjectFunc = inst.inject
	return r.Commit()
}

func (inst* pfde5cc127c_com_MockUnit) new() any {
    return &pfde5cc127.MockUnit{}
}

func (inst* pfde5cc127c_com_MockUnit) inject(injext application.InjectionExt, instance any) error {
	ie := injext
	com := instance.(*pfde5cc127.MockUnit)
	nop(ie, com)

	
    com.Service = inst.getService(ie)


    return nil
}


func (inst*pfde5cc127c_com_MockUnit) getService(ie application.InjectionExt)p262c04a06.Service{
    return ie.GetComponent("#alias-262c04a06c32904104382e2b8d56c279-Service").(p262c04a06.Service)
}


