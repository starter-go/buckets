package com

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"io"
	"sort"
	"strconv"

	"github.com/starter-go/afs"
	"github.com/starter-go/base/lang"
	"github.com/starter-go/buckets"
	"github.com/starter-go/buckets/src/test/golang/unit"
	"github.com/starter-go/vlog"

	"github.com/starter-go/units"
)

type CommonCrudUnits struct {

	//starter:component

	_as func(units.Unit) //starter:as(".")

	Service    buckets.Service  //starter:inject("#")
	DirManager units.DirManager //starter:inject("#")

	BucketName string //starter:inject("${units.common-crud.bucket-name}")

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

	u7 := &units.Registration{
		Name:     unit.TheCommonCrudUnits + "sum-api",
		Enabled:  true,
		Priority: 13,
		Do:       inst.runTrySumAPI,
	}

	u8 := &units.Registration{
		Name:     unit.TheCommonCrudUnits + "meta-io",
		Enabled:  true,
		Priority: 12,
		Do:       inst.runTryMetaIO,
	}

	list = append(list, u1)
	list = append(list, u2)
	list = append(list, u3)
	list = append(list, u4)
	list = append(list, u5)
	list = append(list, u6)
	list = append(list, u7)
	list = append(list, u8)

	return list
}

func (inst *CommonCrudUnits) prepareTesting(cc context.Context) (*innerCommonCrudTesting, error) {

	t := new(innerCommonCrudTesting)
	sel := inst.BucketName
	ser := inst.Service

	bkt, err := ser.GetBucket(cc, sel)
	if err != nil {
		return nil, err
	}

	t.cc = cc
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

func (inst *CommonCrudUnits) runInsert(cc context.Context) error {

	t, err := inst.prepareTesting(cc)
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

func (inst *CommonCrudUnits) runUpdate(cc context.Context) error {

	t, err := inst.prepareTesting(cc)
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

func (inst *CommonCrudUnits) runFetch(cc context.Context) error {

	t, err := inst.prepareTesting(cc)
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

func (inst *CommonCrudUnits) runRemove(cc context.Context) error {

	t, err := inst.prepareTesting(cc)
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

func (inst *CommonCrudUnits) runFileRead(cc context.Context) error {

	t, err := inst.prepareTesting(cc)
	if err != nil {
		return err
	}
	steps := t.steps

	steps = append(steps, t.doInit)
	steps = append(steps, t.doFileRead)
	steps = append(steps, t.doComputeAllSum)
	steps = append(steps, t.doCheck)
	steps = append(steps, t.doLog)

	t.steps = steps
	return inst.runSteps(t)
}

func (inst *CommonCrudUnits) runFileWrite(cc context.Context) error {

	t, err := inst.prepareTesting(cc)
	if err != nil {
		return err
	}
	steps := t.steps

	steps = append(steps, t.doInit)
	steps = append(steps, t.doFileWrite)
	steps = append(steps, t.doCheck)
	steps = append(steps, t.doLog)

	t.steps = steps
	return inst.runSteps(t)
}

func (inst *CommonCrudUnits) runTryMetaIO(cc context.Context) error {

	t, err := inst.prepareTesting(cc)
	if err != nil {
		return err
	}
	steps := t.steps

	steps = append(steps, t.doInit)
	steps = append(steps, t.doInsert)
	steps = append(steps, t.doTryMetaIO)

	// steps = append(steps, t.doCheck)
	// steps = append(steps, t.doLog)

	t.steps = steps
	return inst.runSteps(t)
}

func (inst *CommonCrudUnits) runTrySumAPI(cc context.Context) error {

	t, err := inst.prepareTesting(cc)
	if err != nil {
		return err
	}
	steps := t.steps

	steps = append(steps, t.doInit)
	steps = append(steps, t.doInsert)
	steps = append(steps, t.doTrySumAPI)

	// steps = append(steps, t.doCheck)
	// steps = append(steps, t.doLog)

	t.steps = steps
	return inst.runSteps(t)
}

////////////////////////////////////////////////////////////////////////////////

type innerCommonCrudTesting struct {
	cc context.Context

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

func (inst *innerCommonCrudTesting) doFileRead(u *CommonCrudUnits) error {

	cc := inst.cc
	bucket := inst.bucket
	fapi := bucket.ForFiles()

	data1 := inst.body1data
	sum1 := inst.innerComputeSum(data1)
	name1 := buckets.ObjectName("objects/" + sum1 + "/data")

	oName := name1
	inst.oName = name1

	// file & dir

	hDir := &units.DirHolder{
		Context: cc,
		Key:     units.DirKeyTemp,
		Scope:   units.DirScopeRuntime,
	}

	hDir, err := u.DirManager.GetDir(hDir)
	if err != nil {
		return err
	}

	dir := hDir.Path
	file := dir.GetChild("file-to-read.demo")

	// objects

	data1r := bytes.NewReader(data1)
	data1rc := io.NopCloser(data1r)
	o1 := &buckets.Object{
		Context: cc,
		Name:    oName,
		Data:    data1rc,
	}

	o2 := &buckets.ObjectFile{
		Path: file,
	}
	o2.Name = oName

	// write (as mem)

	o1, err = bucket.Put(o1)
	if err != nil {
		return err
	}

	// read (as file)

	o2, err = fapi.FetchFile(o2)
	if err != nil {
		return err
	}

	// load body2

	bin3, err := file.GetIO().ReadBinary(nil)
	if err != nil {
		return err
	}
	sum3 := inst.innerComputeSum(bin3)
	inst.body2data = bin3
	inst.body2sum = sum3

	return nil
}

func (inst *innerCommonCrudTesting) doFileWrite(u *CommonCrudUnits) error {

	cc := inst.cc
	bucket := inst.bucket
	fapi := bucket.ForFiles()

	data1 := inst.body1data
	sum1 := inst.innerComputeSum(data1)
	name1 := buckets.ObjectName("objects/" + sum1 + "/data")

	inst.oName = name1

	// file & dir

	hDir := &units.DirHolder{
		Context: cc,
		Key:     units.DirKeyTemp,
		Scope:   units.DirScopeRuntime,
	}

	hDir, err := u.DirManager.GetDir(hDir)
	if err != nil {
		return err
	}

	dir := hDir.Path
	file := dir.GetChild("file-to-write.demo")

	if !dir.Exists() {
		om := new(afs.OptionsMaker)
		// om.Create().WriteOnly()
		om.SetMode(7, 5, 5)
		opt := om.Options()
		dir.Mkdirs(&opt)
	}

	// prepare demo file

	om := new(afs.OptionsMaker)
	om.Create().WriteOnly()
	om.SetMode(6, 4, 4)
	opt := om.Options()

	err = file.GetIO().WriteBinary(data1, &opt)
	if err != nil {
		return err
	}

	// write (as file)

	fo1 := new(buckets.ObjectFile)
	fo1.Name = name1
	fo1.Context = cc
	fo1.Path = file
	fo1.Type = "application/x-bin-data"

	fo1, err = fapi.PutFile(fo1)
	if err != nil {
		return err
	}

	// read (as mem)

	o2 := new(buckets.Object)
	o2.Name = name1
	o2.Context = cc

	o2, err = bucket.Fetch(o2)
	if err != nil {
		return err
	}

	o2reader := o2.Data
	if o2reader == nil {
		return fmt.Errorf("o2reader is nil")
	}
	defer o2reader.Close()

	o2data, err := io.ReadAll(o2reader)
	if err != nil {
		return err
	}

	inst.body2data = o2data

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

func (inst *innerCommonCrudTesting) doTryMetaIO(u *CommonCrudUnits) error {

	cc := inst.cc
	bucket := inst.bucket
	name := inst.oName

	o1 := bucket.GetObject(name)
	o2 := bucket.GetObject(name)

	o1.Context = cc
	o2.Context = cc

	// set meta

	meta := o1.Meta
	meta = meta.Init()
	o1.Meta = meta

	meta.Put("x", "foo")
	meta.Put("y", "bar")
	meta.Put("z", "hello")

	o1, err := bucket.SetMeta(o1)
	if err != nil {
		return err
	}

	// get meta

	o2, err = bucket.GetMeta(o2)
	if err != nil {
		return err
	}

	vlog.Info("Meta:")
	meta = o2.Meta
	namelist := make([]string, 0)

	for k := range meta {
		namelist = append(namelist, k.String())
	}

	sort.Strings(namelist)

	for _, name := range namelist {
		value := meta.Get(buckets.MetaName(name))
		vlog.Info("    %s = %s", name, value)
	}

	return nil
}

func (inst *innerCommonCrudTesting) doTrySumAPI(u *CommonCrudUnits) error {

	cc := inst.cc
	bucket := inst.bucket
	name := inst.oName
	sumapi := bucket.ForSum()

	o1 := new(buckets.Object)
	o1.Context = cc
	o1.Name = name

	o1, err := bucket.Fetch(o1)
	if err != nil {
		return err
	}

	o1reader := o1.Data
	defer o1reader.Close()

	o1raw, err := io.ReadAll(o1reader)
	if err != nil {
		return err
	}

	alg := sumapi.Algorithm()
	ha := sumapi.Hash()
	ha.Write(o1raw)
	sum := []byte{}
	sum = ha.Sum(sum)
	hex := lang.HexFromBytes(sum)
	size := len(o1raw)

	etag := o1.Sum

	vlog.Info("[object name:'%s'      size:%d   etag:'%s']", name, size, etag.String())
	vlog.Info("[object name:'%s' algorithm:'%s'  sum:'%s']", name, alg, hex)

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
