package main4buckets

import "github.com/starter-go/application"

func nop(a ... any) {    
}

func registerComponents(cr application.ComponentRegistry) error {
    ac:=&autoRegistrar{}
    ac.init(cr)
    return ac.addAll()
}

type comFactory interface {
    register(cr application.ComponentRegistry) error
}

type autoRegistrar struct {
    cr application.ComponentRegistry
}

func (inst *autoRegistrar) init(cr application.ComponentRegistry) {
	inst.cr = cr
}

func (inst *autoRegistrar) register(factory comFactory) error {
	return factory.register(inst.cr)
}

func (inst*autoRegistrar) addAll() error {

    
    inst.register(&p019b1834ac_mock_Driver{})
    inst.register(&pac41a78b8d_localfiles_Driver{})
    inst.register(&pf6935ed6a9_core_BucketDriverManagerImpl{})
    inst.register(&pf6935ed6a9_core_BucketServiceImpl{})


    return nil
}
