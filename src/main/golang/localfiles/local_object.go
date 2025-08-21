package localfiles

import (
	"crypto/sha256"

	"fmt"
	"io"

	"github.com/starter-go/afs"
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

	return nil
}

func (inst *innerObjectHolder) writeMeta(o *buckets.Object) error {
	if o == nil {
		return fmt.Errorf("object is nil")
	}

	return nil
}

func (inst *innerObjectHolder) computeMeta(o *buckets.Object) error {
	if o == nil {
		return fmt.Errorf("object is nil")
	}

	data := inst.dataFile
	sum, size, err := inst.computeSHA256sum(data)
	if err != nil {
		return err
	}

	// info := data.GetInfo()

	o.Existed = true
	o.Size = size
	o.Sum = sum
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
	if err == io.EOF {
		bin := dst.Sum(nil)
		sum.Algorithm = buckets.AlgorithmSHA256
		sum.Value = lang.HexFromBytes(bin)
		size = cb
	} else if err != nil {
		return sum, size, err
	}

	return sum, size, nil
}
