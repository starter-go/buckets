package localfiles

import (
	"crypto/sha256"

	"fmt"
	"io"

	"github.com/starter-go/afs"
	"github.com/starter-go/application/properties"
	"github.com/starter-go/base/lang"
	"github.com/starter-go/buckets"
)

type innerObjectHolder struct {
	object   buckets.Object
	dataFile afs.Path
	metaFile afs.Path
}

func (inst *innerObjectHolder) openDataReader() (io.ReadCloser, error) {
	rdr := new(innerDataReader)
	rdr.dataFile = inst.dataFile
	return rdr.as(), nil
}

func (inst *innerObjectHolder) openDataWriter() (io.WriteCloser, error) {
	wtr := new(innerDataWriter)
	wtr.dataFile = inst.dataFile
	return wtr.as(), nil
}

func (inst *innerObjectHolder) readMeta(o *buckets.Object) error {

	if o == nil {
		return fmt.Errorf("object is nil")
	}

	buf := new(innerMetaBuffer)
	err := inst.loadMeta(o, buf)
	if err != nil {
		return err
	}

	o.Meta = buf.table
	return nil
}

func (inst *innerObjectHolder) loadMeta(o *buckets.Object, dst *innerMetaBuffer) error {

	if o == nil || dst == nil {
		return fmt.Errorf("params: 'obj:buckets.Object' or 'dst:innerMetaBuffer' is nil")
	}

	file := inst.metaFile
	raw, err := file.GetIO().ReadText(nil)
	if err != nil {
		return err
	}

	ptable, err := properties.Parse(raw, nil)
	if err != nil {
		return err
	}

	tmp := ptable.Export(nil)
	for key, val := range tmp {
		dst.set(MetaName(key), val)
	}
	return nil
}

func (inst *innerObjectHolder) writeMeta(o *buckets.Object, meta *innerMetaBuffer) error {

	if o == nil {
		return fmt.Errorf("object is nil")
	}

	const base = 10
	// now := lang.Now()
	// ts := strconv.FormatInt(now.Int(), base)

	ptable := properties.NewTable(nil)
	ptable.Import(meta.table.Export(nil))

	// ptable.SetProperty("meta.timestamp", ts)

	txt := properties.Format(ptable, properties.FormatWithGroups)
	file := inst.metaFile

	optmkr := new(afs.OptionsMaker)
	optmkr.SetMode(6, 4, 4)
	optmkr.WriteOnly()
	if file.Exists() {
		optmkr.Truncate()
	} else {
		optmkr.Create()
	}

	opt := optmkr.Options()
	return file.GetIO().WriteText(txt, &opt)
}

func (inst *innerObjectHolder) computeMeta(o *buckets.Object, dst *innerMetaBuffer) error {

	if o == nil {
		return fmt.Errorf("object is nil")
	}

	now := lang.Now()
	data := inst.dataFile
	sum, size, err := inst.computeSHA256sum(data)
	if err != nil {
		return err
	}

	hasOlder1 := dst.contains(MetaNameCreatedAt)
	hasOlder2 := dst.contains(MetaNameCreatedAtTime)
	setter := dst.setter()

	setter.setTimeI64(MetaNameUpdatedAt, now)
	setter.setTimeStr(MetaNameUpdatedAtTime, now)
	setter.setInt64(MetaNameLength, size)
	setter.setString(MetaNameType, o.Type.String())
	setter.setSum(MetaNameSum, &sum)
	if (!hasOlder1) || (!hasOlder2) {
		setter.setTimeI64(MetaNameCreatedAt, now)
		setter.setTimeStr(MetaNameCreatedAtTime, now)
	}

	o.Existed = true
	o.Size = size
	o.Sum = sum
	o.Meta = dst.table

	return nil
}

func (inst *innerObjectHolder) computeSHA256sum(file afs.Path) (buckets.SUM, int64, error) {

	var size int64
	sum := buckets.SUM{}
	src, err := file.GetIO().OpenReader(nil)
	if err != nil {
		return sum, size, err
	}
	defer src.Close()

	dst := sha256.New()
	cb, err := io.Copy(dst, src)
	if err != nil {
		if err != io.EOF {
			return sum, 0, err
		}
	}

	bin := dst.Sum(nil)
	sum.Algorithm = buckets.AlgorithmSHA256
	sum.Value = lang.HexFromBytes(bin)
	size = cb

	return sum, size, nil
}
