package com

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"io"
	"strconv"

	"github.com/starter-go/base/lang"
	"github.com/starter-go/buckets"
	"github.com/starter-go/buckets/src/test/golang/unit"
	"github.com/starter-go/vlog"

	"github.com/starter-go/units"
)

type CommonCrudUnits struct {

	//starter:component

	_as func(units.Unit) //starter:as(".")

	Service    buckets.Service //starter:inject("#")
	BucketName string          //starter:inject("${units.common-crud.bucket-name}")

}

func (inst *CommonCrudUnits) _impl() units.Unit {
	return inst
}

func (inst *CommonCrudUnits) ListRegistrations(list []*units.Registration) []*units.Registration {

	u1 := &units.Registration{
		Name:     unit.TheCommonCrudUnits + "fetch",
		Enabled:  true,
		Priority: 19,
		Do:       inst.runFetch,
	}

	u2 := &units.Registration{
		Name:     unit.TheCommonCrudUnits + "insert",
		Enabled:  true,
		Priority: 18,
		Do:       inst.runInsert,
	}

	u3 := &units.Registration{
		Name:     unit.TheCommonCrudUnits + "update",
		Enabled:  true,
		Priority: 17,
		Do:       inst.runUpdate,
	}

	u4 := &units.Registration{
		Name:     unit.TheCommonCrudUnits + "remove",
		Enabled:  true,
		Priority: 16,
		Do:       inst.runRemove,
	}

	u5 := &units.Registration{
		Name:     unit.TheCommonCrudUnits + "fread",
		Enabled:  true,
		Priority: 15,
		Do:       inst.runFileRead,
	}

	u6 := &units.Registration{
		Name:     unit.TheCommonCrudUnits + "fwrite",
		Enabled:  true,
		Priority: 14,
		Do:       inst.runFileWrite,
	}

	list = append(list, u1)
	list = append(list, u2)
	list = append(list, u3)
	list = append(list, u4)
	list = append(list, u5)
	list = append(list, u6)

	return list
}

// func (inst *CommonCrudUnits) runTest() error {

// 	holder := buckets.BucketHolder{}
// 	ser := inst.Service
// 	name := "mock.demo"

// 	err := holder.SetName(name).SetService(ser).Init()
// 	if err != nil {
// 		return err
// 	}

// 	// bucket
// 	bucket, err := holder.GetBucket()
// 	if err != nil {
// 		return err
// 	}

// 	// object
// 	o1 := bucket.GetObject("/a/b/c/d")

// 	// data
// 	data1 := "hello,bucket"
// 	data2 := bytes.NewBufferString(data1)

// 	// write
// 	o1.Data = io.NopCloser(data2)
// 	o1.Type = "application/x-bin"
// 	_, err = bucket.Put(o1)
// 	if err != nil {
// 		return err
// 	}

// 	// read
// 	o2 := bucket.GetObject(o1.Name)
// 	o2, err = bucket.Fetch(o2)
// 	if err != nil {
// 		return err
// 	}

// 	data3 := o2.Data
// 	data4, err := io.ReadAll(data3)
// 	if err != nil {
// 		return err
// 	}
// 	vlog.Info("read object data:")
// 	vlog.Info("    name = %v", o2.Name)
// 	vlog.Info("     len = %v", len(data4))
// 	vlog.Info("     sum = %v", o2.Sum)
// 	vlog.Info("    type = %v", o2.Type)

// 	// todo ...

// 	// meta
// 	// todo ...

// 	return nil
// }

func (inst *CommonCrudUnits) prepareTesting() (*innerCommonCrudTesting, error) {

	t := new(innerCommonCrudTesting)
	ctx := context.Background()
	sel := inst.BucketName
	ser := inst.Service

	bkt, err := ser.GetBucket(ctx, sel)
	if err != nil {
		return nil, err
	}

	t.bucket = bkt
	return t, nil
}

func (inst *CommonCrudUnits) runSteps(t *innerCommonCrudTesting) error {
	all := t.steps
	for _, st := range all {
		err := st(inst)
		if err != nil {
			return err
		}
	}
	return nil
}

func (inst *CommonCrudUnits) runInsert() error {

	t, err := inst.prepareTesting()
	if err != nil {
		return err
	}
	steps := t.steps

	steps = append(steps, t.doInit)

	steps = append(steps, t.doInsert)
	steps = append(steps, t.doLog)

	steps = append(steps, t.doFetch)
	steps = append(steps, t.doLog)

	steps = append(steps, t.doCheck)

	t.steps = steps
	return inst.runSteps(t)
}

func (inst *CommonCrudUnits) runUpdate() error {
	t, err := inst.prepareTesting()
	if err != nil {
		return err
	}
	steps := t.steps

	steps = append(steps, t.doInit)

	steps = append(steps, t.doInsert)
	steps = append(steps, t.doFetch)
	steps = append(steps, t.doLog)
	steps = append(steps, t.doCheck)

	steps = append(steps, t.doRebuildBody1)

	steps = append(steps, t.doUpdate)
	steps = append(steps, t.doFetch)
	steps = append(steps, t.doLog)
	steps = append(steps, t.doCheck)

	t.steps = steps
	return inst.runSteps(t)
}

func (inst *CommonCrudUnits) runFetch() error {
	t, err := inst.prepareTesting()
	if err != nil {
		return err
	}
	steps := t.steps

	steps = append(steps, t.doInit)

	steps = append(steps, t.doInsert)
	steps = append(steps, t.doLog)

	steps = append(steps, t.doFetch)
	steps = append(steps, t.doLog)

	steps = append(steps, t.doCheck)

	t.steps = steps
	return inst.runSteps(t)
}

func (inst *CommonCrudUnits) runRemove() error {
	t, err := inst.prepareTesting()
	if err != nil {
		return err
	}
	steps := t.steps

	steps = append(steps, t.doInit)

	steps = append(steps, t.doInsert)
	steps = append(steps, t.doFetch)
	steps = append(steps, t.doExists)
	steps = append(steps, t.doLog)

	steps = append(steps, t.doCheck)

	steps = append(steps, t.doRemove)
	steps = append(steps, t.doExists)
	steps = append(steps, t.doLog)

	t.steps = steps
	return inst.runSteps(t)
}

func (inst *CommonCrudUnits) runFileRead() error {
	t, err := inst.prepareTesting()
	if err != nil {
		return err
	}
	steps := t.steps

	steps = append(steps, t.doInit)
	steps = append(steps, t.doInsert)
	steps = append(steps, t.doCheck)
	steps = append(steps, t.doFetch)
	steps = append(steps, t.doCheck)
	steps = append(steps, t.doLog)

	t.steps = steps
	return inst.runSteps(t)
}

func (inst *CommonCrudUnits) runFileWrite() error {
	t, err := inst.prepareTesting()
	if err != nil {
		return err
	}
	steps := t.steps

	steps = append(steps, t.doInit)
	steps = append(steps, t.doInsert)
	steps = append(steps, t.doCheck)
	steps = append(steps, t.doFetch)
	steps = append(steps, t.doCheck)
	steps = append(steps, t.doLog)

	t.steps = steps
	return inst.runSteps(t)
}

////////////////////////////////////////////////////////////////////////////////

type innerCommonCrudTesting struct {
	steps []func(*CommonCrudUnits) error

	oName buckets.ObjectName

	// body1 : to write
	body1data []byte
	body1sum  lang.Hex

	// body2 : to read
	body2data []byte
	body2sum  lang.Hex

	bucket buckets.Bucket
}

func (inst *innerCommonCrudTesting) doInit(u *CommonCrudUnits) error {

	return inst.doRebuildBody1(u)
}

func (inst *innerCommonCrudTesting) doRemove(u *CommonCrudUnits) error {

	bucket := inst.bucket
	name := inst.oName

	it1 := new(buckets.Object)
	it1.Name = name

	vlog.Info("do remove object '%s'", name)

	err := bucket.Delete(it1)
	if err != nil {
		return err
	}

	inst.body2data = nil
	inst.body2sum = ""
	return nil
}

func (inst *innerCommonCrudTesting) doFetch(u *CommonCrudUnits) error {

	name := inst.oName
	it1 := new(buckets.Object)
	bucket := inst.bucket

	vlog.Info("do fetch from '%s'", name)

	it1.Name = name
	it2, err := bucket.Fetch(it1)
	if err != nil {
		return err
	}
	src := it2.Data
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	inst.body2data = data
	return nil
}

func (inst *innerCommonCrudTesting) doExists(u *CommonCrudUnits) error {

	name := inst.oName
	it1 := new(buckets.Object)
	bucket := inst.bucket

	vlog.Info("do exists of object '%s'", name)

	it1.Name = name
	yes, err := bucket.Exists(it1)
	if err != nil {
		return err
	}

	vlog.Info("    exist = %v", yes)

	return nil
}

func (inst *innerCommonCrudTesting) doInsert(u *CommonCrudUnits) error {

	bucket := inst.bucket
	it1 := new(buckets.Object)

	src, err := inst.innerGetBody1stream()
	if err != nil {
		return err
	}
	defer src.Close()

	inst.doComputeAllSum(u)
	sum := inst.body1sum
	name1 := buckets.ObjectName("objects/" + sum + "/data")

	it1.Name = name1
	it1.Data = src
	it1.Type = "application/x-bin-stream"

	vlog.Info("do insert to '%s'", name1)

	it2, err := bucket.Put(it1)
	if err != nil {
		return err
	}

	name2 := it2.Name
	if name1 != name2 {
		return fmt.Errorf("name1 != name2")
	}

	inst.oName = name2
	return nil
}

func (inst *innerCommonCrudTesting) doUpdate(u *CommonCrudUnits) error {

	bucket := inst.bucket
	it1 := new(buckets.Object)

	src, err := inst.innerGetBody1stream()
	if err != nil {
		return err
	}
	defer src.Close()

	name1 := inst.oName
	it1.Name = name1
	it1.Data = src
	it1.Type = "text/plain"

	vlog.Info("do update to '%s'", name1)

	it2, err := bucket.Put(it1)
	if err != nil {
		return err
	}

	name2 := it2.Name
	if name1 != name2 {
		return fmt.Errorf("name1 != name2")
	}

	return nil
}

func (inst *innerCommonCrudTesting) doCheck(u *CommonCrudUnits) error {

	// check size

	size1 := len(inst.body1data)
	size2 := len(inst.body2data)
	if size1 != size2 {
		return fmt.Errorf("body.1.size != body.2.size")
	}

	// check sum

	sum1 := inst.body1sum
	sum2 := inst.body2sum
	if sum1 != sum2 {
		return fmt.Errorf("body.1.sum != body.2.sum")
	}

	return nil
}

func (inst *innerCommonCrudTesting) doLog(u *CommonCrudUnits) error {

	inst.doComputeAllSum(u)

	size1 := len(inst.body1data)
	size2 := len(inst.body2data)

	sum1 := inst.body1sum
	sum2 := inst.body2sum

	vlog.Info("[body id:'1-to-wr' size:%d sum:'%s' ]", size1, sum1)
	vlog.Info("[body id:'2-to-rd' size:%d sum:'%s' ]", size2, sum2)

	return nil
}

func (inst *innerCommonCrudTesting) doComputeAllSum(u *CommonCrudUnits) error {

	sum1 := inst.innerComputeSum(inst.body1data)
	sum2 := inst.innerComputeSum(inst.body2data)

	inst.body1sum = sum1
	inst.body2sum = sum2

	return nil
}

func (inst *innerCommonCrudTesting) innerGetBody1stream() (io.ReadCloser, error) {

	data := inst.body1data
	rdr := bytes.NewReader(data)
	rc := io.NopCloser(rdr)

	return rc, nil
}

func (inst *innerCommonCrudTesting) innerComputeSum(d []byte) lang.Hex {

	sum := sha1.Sum(d)
	hex := lang.HexFromBytes(sum[:])

	return hex
}

func (inst *innerCommonCrudTesting) doRebuildBody1(u *CommonCrudUnits) error {

	const (
		nl = "\n"
	)

	size := 128
	raw := make([]byte, size)
	now := lang.Now()
	buf := bytes.NewBufferString("[random-testing-data]")

	rand.Read(raw)

	buf.WriteString(nl + "time:")
	buf.WriteString(now.String())
	buf.WriteString(nl + "content-length:")
	buf.WriteString(strconv.Itoa(size))

	buf.WriteString(nl)
	buf.WriteRune(0)
	buf.Write(raw)

	inst.body1data = buf.Bytes()
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// EOF
