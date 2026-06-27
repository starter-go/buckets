package localfiles

import (
	"context"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"strings"

	"github.com/starter-go/afs"
	"github.com/starter-go/base/lang"
	"github.com/starter-go/buckets"
)

type layout struct {
	configFile      afs.Path // @{workspace}/.bucket/config
	dotBucketDir    afs.Path // @{workspace}/.bucket
	workspaceFolder afs.Path // @{workspace}
}

////////////////////////////////////////////////////////////////////////////////

type innerLocalBucket struct {
	config1 buckets.Configuration
	oo      buckets.OpenOptions
	layout  layout
	context context.Context
}

// SetMeta implements buckets.Bucket.
func (inst *innerLocalBucket) SetMeta(o1 *buckets.Object) (*buckets.Object, error) {

	h, err := inst.innerGetObjectHolder(o1)
	if err != nil {
		return nil, err
	}

	buf := new(innerMetaBuffer)
	name := h.object.Name

	// read
	err = h.loadMeta(o1, buf)
	if err != nil {
		return nil, err
	}

	// append
	src := o1.Meta
	for k, v := range src {
		buf.set(k, v)
	}

	// write
	err = h.writeMeta(o1, buf)
	if err != nil {
		return nil, err
	}

	// make result
	o2 := inst.GetObject(name)
	o2.Bucket = inst
	o2.Context = o1.Context
	o2.Name = name
	o2.Meta = buf.table
	o2.Existed = true

	return o2, err
}

func (inst *innerLocalBucket) _impl() (buckets.Bucket, buckets.BucketFileAPI, buckets.BucketNativeSumAPI) {
	return inst, inst, inst
}

// ForFiles implements buckets.Bucket.
func (inst *innerLocalBucket) ForFiles() buckets.BucketFileAPI {
	return inst
}

// ForSum implements buckets.Bucket.
func (inst *innerLocalBucket) ForSum() buckets.BucketNativeSumAPI {
	return inst
}

// Algorithm implements buckets.BucketNativeSumAPI.
func (inst *innerLocalBucket) Algorithm() buckets.CheckSumAlgorithm {
	return buckets.AlgorithmSHA256
}

// Hash implements buckets.BucketNativeSumAPI.
func (inst *innerLocalBucket) Hash() hash.Hash {
	return sha256.New()
}

// Bucket implements buckets.BucketFileAPI.
func (inst *innerLocalBucket) Bucket() buckets.Bucket {
	return inst
}

// FetchFile implements buckets.BucketFileAPI.
func (inst *innerLocalBucket) FetchFile(o1 *buckets.ObjectFile) (*buckets.ObjectFile, error) {

	if o1 == nil {
		return nil, fmt.Errorf("param: 'object(want)' is nil")
	}
	holder, err := inst.innerGetObjectHolder(&o1.Object)
	if err != nil {
		return nil, err
	}

	srcFile := holder.dataFile
	dstFile := o1.Path

	/// io

	ctrl := new(innerTempFileCtrl)
	ctrl.init(dstFile)
	defer ctrl.finish()

	srcReader, err := srcFile.GetIO().OpenReader(nil)
	if err != nil {
		return nil, err
	}
	defer srcReader.Close()

	tmpFile, opt, err := ctrl.prepareToWrite()
	if err != nil {
		return nil, err
	}

	tmpWriter, err := tmpFile.GetIO().OpenWriter(opt)
	if err != nil {
		return nil, err
	}
	defer tmpWriter.Close()

	cnt, err := io.Copy(tmpWriter, srcReader)
	if err != nil {
		return nil, err
	}

	err = ctrl.commit()
	if err != nil {
		return nil, err
	}

	// make result

	o2 := new(buckets.ObjectFile)
	o2.Context = o1.Context
	o2.Bucket = inst
	o2.Name = o1.Name
	o2.Size = cnt

	// load meta
	metaBuff := new(innerMetaBuffer)
	holder.loadMeta(&o2.Object, metaBuff)
	o2.Meta = metaBuff.table

	return o2, nil
}

// PutFile implements buckets.BucketFileAPI.
func (inst *innerLocalBucket) PutFile(o1 *buckets.ObjectFile) (*buckets.ObjectFile, error) {

	if o1 == nil {
		return nil, fmt.Errorf("param: 'object(want)' is nil")
	}

	holder, err := inst.innerGetObjectHolder(&o1.Object)
	if err != nil {
		return nil, err
	}

	ctrl := new(innerTempFileCtrl)
	dst := holder.dataFile
	src := o1.Path
	hasOlder := dst.IsFile()

	err = ctrl.init(dst)
	if err != nil {
		return nil, err
	}

	tmp3, opt3, err := ctrl.prepareToWrite()
	if err != nil {
		return nil, err
	}

	defer ctrl.finish()

	// open reader

	om := new(afs.OptionsMaker)
	om.ReadOnly()
	opt1 := om.Options()
	srcReader, err := src.GetIO().OpenReader(&opt1)
	if err != nil {
		return nil, err
	}
	defer srcReader.Close()

	// open writer

	tmpWriter, err := tmp3.GetIO().OpenWriter(opt3)
	if err != nil {
		return nil, err
	}
	defer tmpWriter.Close()

	// copy

	count, err := io.Copy(tmpWriter, srcReader)
	if err != nil {
		return nil, err
	}

	tmpWriter.Close()

	// commit

	err = ctrl.commit()
	if err != nil {
		return nil, err
	}

	//  compute  meta

	now := lang.Now()
	metaBuf := new(innerMetaBuffer)
	metaSett := metaBuf.setter()
	metaGett := metaBuf.getter()

	if hasOlder {
		holder.loadMeta(&o1.Object, metaBuf)
	} else {
		metaSett.setTimeI64(MetaNameCreatedAt, now)
		metaSett.setTimeStr(MetaNameCreatedAtTime, now)
	}
	holder.computeMeta(&o1.Object, metaBuf)

	// write meta

	holder.writeMeta(&o1.Object, metaBuf)

	// make result

	sum := metaGett.getSum(MetaNameSum)
	res := new(buckets.ObjectFile)

	if sum != nil {
		res.Sum = *sum
	}

	res.Name = o1.Name
	res.Bucket = inst
	res.Size = count
	res.Type = o1.Type
	res.Meta = metaBuf.table
	res.Context = o1.Context
	res.Path = o1.Path
	res.Existed = true

	return res, nil
}

// Delete implements buckets.Bucket.
func (inst *innerLocalBucket) Delete(o1 *buckets.Object) error {

	h, err := inst.innerGetObjectHolder(o1)
	if err != nil {
		return err
	}

	f1 := h.dataFile
	f2 := h.metaFile
	flist := []afs.Path{f1, f2}
	count := 0

	for _, p := range flist {
		if p == nil {
			continue
		}
		if !p.Exists() {
			continue
		}
		err := p.Delete()
		if err != nil {
			return err
		}
		count++
	}

	if count == 0 {
		name := o1.Name
		return fmt.Errorf("no object with name: %s", name)
	}
	return nil
}

func (inst *innerLocalBucket) isChildOf(child, parent afs.Path) bool {
	sep := child.GetFS().Separator()
	path1 := parent.GetPath()
	path2 := child.GetPath()
	if !strings.HasSuffix(path1, sep) {
		path1 = path1 + sep
	}
	return strings.HasPrefix(path2, path1)
}

func (inst *innerLocalBucket) computeObjectPath(name buckets.ObjectName) afs.Path {
	dot := inst.layout.dotBucketDir
	wks := inst.layout.workspaceFolder
	child := wks.GetChild(name.String())
	if !inst.isChildOf(child, dot) && inst.isChildOf(child, wks) {
		return child
	}
	panic("bad path of object: " + name.String())
}

func (inst *innerLocalBucket) innerGetObjectHolder(o1 *buckets.Object) (*innerObjectHolder, error) {

	if o1 == nil {
		return nil, fmt.Errorf("object is nil")
	}

	path := inst.computeObjectPath(o1.Name)

	h := new(innerObjectHolder)
	h.dataFile = path.GetChild("object.data")
	h.metaFile = path.GetChild("object.meta")
	h.object = *o1
	h.object.Bucket = inst

	return h, nil
}

func (inst *innerLocalBucket) SetContext(ctx context.Context) buckets.Bucket {

	if ctx == nil {
		return inst
	}

	inst.context = ctx
	return inst
}

func (inst *innerLocalBucket) GetContext() context.Context {

	ctx := inst.context

	if ctx == nil {
		ctx = context.Background()
		inst.context = ctx
	}

	return ctx
}

func (inst *innerLocalBucket) GetObject(name buckets.ObjectName) *buckets.Object {
	o := new(buckets.Object)
	o.Bucket = inst
	o.Name = name
	o.Size = -1
	return o
}

func (inst *innerLocalBucket) Fetch(o1 *buckets.Object) (*buckets.Object, error) {

	h, err := inst.innerGetObjectHolder(o1)
	if err != nil {
		return nil, err
	}

	data, err := h.openDataReader()
	if err != nil {
		return nil, err
	}

	o2 := inst.GetObject(o1.Name)
	err = h.readMeta(o2)
	if err != nil {
		return nil, err
	}

	// load meta
	metaBuff := new(innerMetaBuffer)
	h.loadMeta(o2, metaBuff)

	o2.Meta = metaBuff.table
	o2.Data = data
	return o2, nil
}

func (inst *innerLocalBucket) Put(o1 *buckets.Object) (*buckets.Object, error) {

	h, err := inst.innerGetObjectHolder(o1)
	if err != nil {
		return nil, err
	}

	o2 := inst.GetObject(h.object.Name)
	meta := new(innerMetaBuffer)

	wtr, err := h.openDataWriter()
	if err != nil {
		return nil, err
	}
	defer wtr.Close()

	var total int64 = 0
	rdr := o1.Data
	if rdr != nil {
		cb, err := io.Copy(wtr, rdr)
		if err == io.EOF {
			total = cb
		} else if err != nil {
			return nil, err
		}
	}

	o2.Type = o1.Type

	h.loadMeta(o2, meta)

	err = h.computeMeta(o2, meta)
	if err != nil {
		return nil, err
	}
	o2.Size = total

	err = h.writeMeta(o2, meta)
	if err != nil {
		return nil, err
	}

	// make result
	err = h.readMeta(o2)
	return o2, err
}

func (inst *innerLocalBucket) GetMeta(o1 *buckets.Object) (*buckets.Object, error) {

	h, err := inst.innerGetObjectHolder(o1)
	if err != nil {
		return nil, err
	}

	name := h.object.Name
	o2 := inst.GetObject(name)
	err = h.readMeta(o2)
	return o2, err
}

func (inst *innerLocalBucket) Exists(o1 *buckets.Object) (bool, error) {
	h, err := inst.innerGetObjectHolder(o1)
	if err != nil {
		return false, err
	}
	ex1 := h.metaFile.IsFile()
	ex2 := h.dataFile.IsFile()
	return (ex1 && ex2), nil
}

////////////////////////////////////////////////////////////////////////////////

type innerFileCopier struct {
	force bool
}

func (inst *innerFileCopier) copy(src, dst afs.Path) error {
	if !inst.force {
		if inst.innerIsEqual(src, dst) {
			return nil // skip
		}
	}
	return inst.innerDoCopy(src, dst)
}

func (inst *innerFileCopier) innerDoCopy(src, dst afs.Path) error {

	if src == nil || dst == nil {
		return fmt.Errorf("todo")
	}

	path1 := src.GetPath()
	path2 := dst.GetPath()
	if path1 == path2 {
		return nil // ignore
	}

	info := src.GetInfo()
	optionR := inst.innerGetFileOptionR()
	optionW := inst.innerGetFileOptionW()
	tmp := inst.innerNewTempFile(dst)
	defer inst.innerClearTempFile(tmp)

	rdr, err := src.GetIO().OpenReader(optionR)
	if err != nil {
		return err
	}
	defer rdr.Close()

	wtr, err := tmp.GetIO().OpenWriter(optionW)
	if err != nil {
		return err
	}
	defer wtr.Close()

	count, err := io.Copy(wtr, rdr)
	if err != nil {
		return err
	}

	rdr.Close()
	wtr.Close()

	if count != info.Length() {
		return fmt.Errorf("todo")
	}

	if dst.IsFile() {
		dst.Delete()
	}

	return tmp.MoveTo(dst, optionW)
}

func (inst *innerFileCopier) innerNewTempFile(src afs.Path) afs.Path {
	name := src.GetName()
	parent := src.GetParent()
	fsys := src.GetFS()
	tmp, err := fsys.CreateTempFile(name, ".tmp~", parent)
	if err != nil {
		panic(err)
	}
	return tmp
}

func (inst *innerFileCopier) innerClearTempFile(tmp afs.Path) {
	if tmp.IsFile() {
		tmp.Delete()
	}
}

func (inst *innerFileCopier) innerIsEqual(src, dst afs.Path) bool {

	if src == nil || dst == nil {
		return false
	}

	len1 := src.GetInfo().Length()
	len2 := dst.GetInfo().Length()
	if len1 != len2 {
		return false
	}

	sum1, _ := inst.innerComputeSum(src)
	sum2, _ := inst.innerComputeSum(dst)

	return sum1 == sum2
}

func (inst *innerFileCopier) innerComputeSum(file afs.Path) (lang.Hex, error) {

	opt := inst.innerGetFileOptionR()
	src, err := file.GetIO().OpenReader(opt)
	if err != nil {
		return "", err
	}
	defer src.Close()

	sha := sha256.New()
	io.Copy(sha, src)
	sum := sha.Sum(nil)

	hex := lang.HexFromBytes(sum[:])
	return hex, nil
}

func (inst *innerFileCopier) innerGetFileOptionR() *afs.Options {
	return afs.ToReadFile()
}

func (inst *innerFileCopier) innerGetFileOptionW() *afs.Options {
	return afs.ToCreateFile()
}

////////////////////////////////////////////////////////////////////////////////
