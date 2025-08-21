package com

import (
	"bytes"
	"io"

	"github.com/starter-go/buckets"
	"github.com/starter-go/buckets/src/test/golang/unit"
	"github.com/starter-go/vlog"

	"github.com/starter-go/units"
)

type MockUnit struct {

	//starter:component

	_as func(units.Units) //starter:as(".")

	Service buckets.Service //starter:inject("#")

}

func (inst *MockUnit) _impl() units.Units {
	return inst
}

func (inst *MockUnit) Units(list []*units.Registration) []*units.Registration {

	u1 := &units.Registration{
		Name:     unit.TheMockUnit,
		Enabled:  true,
		Priority: 1,
		Test:     inst.runTest,
	}

	list = append(list, u1)
	return list
}

func (inst *MockUnit) runTest() error {

	holder := buckets.BucketHolder{}
	ser := inst.Service
	name := "mock.demo"

	err := holder.SetName(name).SetService(ser).Init()
	if err != nil {
		return err
	}

	// bucket
	bucket, err := holder.GetBucket()
	if err != nil {
		return err
	}

	// object
	o1 := bucket.GetObject("/a/b/c/d")

	// data
	data1 := "hello,bucket"
	data2 := bytes.NewBufferString(data1)

	// write
	o1.Data = io.NopCloser(data2)
	o1.Type = "application/x-bin"
	_, err = bucket.Put(o1)
	if err != nil {
		return err
	}

	// read
	o2 := bucket.GetObject(o1.Name)
	o2, err = bucket.Fetch(o2)
	if err != nil {
		return err
	}

	data3 := o2.Data
	data4, err := io.ReadAll(data3)
	if err != nil {
		return err
	}
	vlog.Info("read object data:")
	vlog.Info("    name = %v", o2.Name)
	vlog.Info("     len = %v", len(data4))
	vlog.Info("     sum = %v", o2.Sum)
	vlog.Info("    type = %v", o2.Type)

	// todo ...

	// meta
	// todo ...

	return nil
}
