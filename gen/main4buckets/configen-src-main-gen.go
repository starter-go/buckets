package main4buckets
import (
    p0d2a11d16 "github.com/starter-go/afs"
    p0ef6f2938 "github.com/starter-go/application"
    p262c04a06 "github.com/starter-go/buckets"
    pf6935ed6a "github.com/starter-go/buckets/src/main/golang/core"
    pac41a78b8 "github.com/starter-go/buckets/src/main/golang/localfiles"
    p019b1834a "github.com/starter-go/buckets/src/main/golang/mock"
     "github.com/starter-go/application"
)

// type pf6935ed6a.BucketDriverManagerImpl in package:github.com/starter-go/buckets/src/main/golang/core
//
// id:com-f6935ed6a99cef59-core-BucketDriverManagerImpl
// class:
// alias:alias-262c04a06c32904104382e2b8d56c279-DriverManager
// scope:singleton
//
type pf6935ed6a9_core_BucketDriverManagerImpl struct {
}

func (inst* pf6935ed6a9_core_BucketDriverManagerImpl) register(cr application.ComponentRegistry) error {
	r := cr.NewRegistration()
	r.ID = "com-f6935ed6a99cef59-core-BucketDriverManagerImpl"
	r.Classes = ""
	r.Aliases = "alias-262c04a06c32904104382e2b8d56c279-DriverManager"
	r.Scope = "singleton"
	r.NewFunc = inst.new
	r.InjectFunc = inst.inject
	return r.Commit()
}

func (inst* pf6935ed6a9_core_BucketDriverManagerImpl) new() any {
    return &pf6935ed6a.BucketDriverManagerImpl{}
}

func (inst* pf6935ed6a9_core_BucketDriverManagerImpl) inject(injext application.InjectionExt, instance any) error {
	ie := injext
	com := instance.(*pf6935ed6a.BucketDriverManagerImpl)
	nop(ie, com)

	
    com.RawDriverList = inst.getRawDriverList(ie)


    return nil
}


func (inst*pf6935ed6a9_core_BucketDriverManagerImpl) getRawDriverList(ie application.InjectionExt)[]p262c04a06.DriverRegistry{
    dst := make([]p262c04a06.DriverRegistry, 0)
    src := ie.ListComponents(".class-262c04a06c32904104382e2b8d56c279-DriverRegistry")
    for _, item1 := range src {
        item2 := item1.(p262c04a06.DriverRegistry)
        dst = append(dst, item2)
    }
    return dst
}



// type pf6935ed6a.BucketServiceImpl in package:github.com/starter-go/buckets/src/main/golang/core
//
// id:com-f6935ed6a99cef59-core-BucketServiceImpl
// class:
// alias:alias-262c04a06c32904104382e2b8d56c279-Service
// scope:singleton
//
type pf6935ed6a9_core_BucketServiceImpl struct {
}

func (inst* pf6935ed6a9_core_BucketServiceImpl) register(cr application.ComponentRegistry) error {
	r := cr.NewRegistration()
	r.ID = "com-f6935ed6a99cef59-core-BucketServiceImpl"
	r.Classes = ""
	r.Aliases = "alias-262c04a06c32904104382e2b8d56c279-Service"
	r.Scope = "singleton"
	r.NewFunc = inst.new
	r.InjectFunc = inst.inject
	return r.Commit()
}

func (inst* pf6935ed6a9_core_BucketServiceImpl) new() any {
    return &pf6935ed6a.BucketServiceImpl{}
}

func (inst* pf6935ed6a9_core_BucketServiceImpl) inject(injext application.InjectionExt, instance any) error {
	ie := injext
	com := instance.(*pf6935ed6a.BucketServiceImpl)
	nop(ie, com)

	
    com.Drivers = inst.getDrivers(ie)
    com.AC = inst.getAC(ie)


    return nil
}


func (inst*pf6935ed6a9_core_BucketServiceImpl) getDrivers(ie application.InjectionExt)p262c04a06.DriverManager{
    return ie.GetComponent("#alias-262c04a06c32904104382e2b8d56c279-DriverManager").(p262c04a06.DriverManager)
}


func (inst*pf6935ed6a9_core_BucketServiceImpl) getAC(ie application.InjectionExt)p0ef6f2938.Context{
    return ie.GetContext()
}



// type pac41a78b8.Driver in package:github.com/starter-go/buckets/src/main/golang/localfiles
//
// id:com-ac41a78b8d91e909-localfiles-Driver
// class:class-262c04a06c32904104382e2b8d56c279-DriverRegistry
// alias:
// scope:singleton
//
type pac41a78b8d_localfiles_Driver struct {
}

func (inst* pac41a78b8d_localfiles_Driver) register(cr application.ComponentRegistry) error {
	r := cr.NewRegistration()
	r.ID = "com-ac41a78b8d91e909-localfiles-Driver"
	r.Classes = "class-262c04a06c32904104382e2b8d56c279-DriverRegistry"
	r.Aliases = ""
	r.Scope = "singleton"
	r.NewFunc = inst.new
	r.InjectFunc = inst.inject
	return r.Commit()
}

func (inst* pac41a78b8d_localfiles_Driver) new() any {
    return &pac41a78b8.Driver{}
}

func (inst* pac41a78b8d_localfiles_Driver) inject(injext application.InjectionExt, instance any) error {
	ie := injext
	com := instance.(*pac41a78b8.Driver)
	nop(ie, com)

	
    com.Enabled = inst.getEnabled(ie)
    com.Priority = inst.getPriority(ie)
    com.AFS = inst.getAFS(ie)


    return nil
}


func (inst*pac41a78b8d_localfiles_Driver) getEnabled(ie application.InjectionExt)bool{
    return ie.GetBool("${buckets-driver.file.enabled}")
}


func (inst*pac41a78b8d_localfiles_Driver) getPriority(ie application.InjectionExt)int{
    return ie.GetInt("${buckets-driver.file.priority}")
}


func (inst*pac41a78b8d_localfiles_Driver) getAFS(ie application.InjectionExt)p0d2a11d16.FS{
    return ie.GetComponent("#alias-0d2a11d163e349503a64168a1cdf48a2-FS").(p0d2a11d16.FS)
}



// type p019b1834a.Driver in package:github.com/starter-go/buckets/src/main/golang/mock
//
// id:com-019b1834ac6506d8-mock-Driver
// class:class-262c04a06c32904104382e2b8d56c279-DriverRegistry
// alias:
// scope:singleton
//
type p019b1834ac_mock_Driver struct {
}

func (inst* p019b1834ac_mock_Driver) register(cr application.ComponentRegistry) error {
	r := cr.NewRegistration()
	r.ID = "com-019b1834ac6506d8-mock-Driver"
	r.Classes = "class-262c04a06c32904104382e2b8d56c279-DriverRegistry"
	r.Aliases = ""
	r.Scope = "singleton"
	r.NewFunc = inst.new
	r.InjectFunc = inst.inject
	return r.Commit()
}

func (inst* p019b1834ac_mock_Driver) new() any {
    return &p019b1834a.Driver{}
}

func (inst* p019b1834ac_mock_Driver) inject(injext application.InjectionExt, instance any) error {
	ie := injext
	com := instance.(*p019b1834a.Driver)
	nop(ie, com)

	
    com.Enabled = inst.getEnabled(ie)
    com.Priority = inst.getPriority(ie)


    return nil
}


func (inst*p019b1834ac_mock_Driver) getEnabled(ie application.InjectionExt)bool{
    return ie.GetBool("${buckets-driver.mock.enabled}")
}


func (inst*p019b1834ac_mock_Driver) getPriority(ie application.InjectionExt)int{
    return ie.GetInt("${buckets-driver.mock.priority}")
}


