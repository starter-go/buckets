package localfiles

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/starter-go/base/lang"
	"github.com/starter-go/buckets"
)

type MetaName string

const (
	MetaNamePrefix MetaName = "object."

	MetaNameLength MetaName = MetaNamePrefix + "length"
	MetaNameType   MetaName = MetaNamePrefix + "type"
	MetaNameSum    MetaName = MetaNamePrefix + "sum"

	MetaNameCreatedAt     MetaName = MetaNamePrefix + "created-at"
	MetaNameCreatedAtTime MetaName = MetaNamePrefix + "created-at-time"
	MetaNameUpdatedAt     MetaName = MetaNamePrefix + "updated-at"
	MetaNameUpdatedAtTime MetaName = MetaNamePrefix + "updated-at-time"
)

////////////////////////////////////////////////////////////////////////////////

type innerMetaBuffer struct {
	table map[string]string
}

func (inst *innerMetaBuffer) set(name MetaName, value string) {
	if name == "" || value == "" {
		return
	}
	t := inst.table
	if t == nil {
		t = make(map[string]string)
		inst.table = t
	}
	key := string(name)
	t[key] = value
}

func (inst *innerMetaBuffer) get(name MetaName) string {
	t := inst.table
	if t == nil {
		return ""
	}
	key := string(name)
	return t[key]
}

func (inst *innerMetaBuffer) getRequiredValue(name MetaName) (string, error) {
	value := inst.get(name)
	if value == "" {
		return "", fmt.Errorf("no meta value with name '%s'", name)
	}
	return value, nil
}

func (inst *innerMetaBuffer) contains(name MetaName) bool {
	value := inst.get(name)
	return (value != "")
}

func (inst *innerMetaBuffer) getter() *innerMetaGetter {
	return &innerMetaGetter{
		buffer: inst,
	}
}

func (inst *innerMetaBuffer) setter() *innerMetaSetter {
	return &innerMetaSetter{
		buffer: inst,
	}
}

////////////////////////////////////////////////////////////////////////////////

type innerMetaGetter struct {
	buffer *innerMetaBuffer
}

////////////////////////////////////////////////////////////////////////////////

type innerMetaSetter struct {
	buffer *innerMetaBuffer
}

func (inst *innerMetaSetter) setTimeI64(name MetaName, value lang.Time) {
	n := value.Int()
	inst.setInt64(name, n)
}

func (inst *innerMetaSetter) setTimeStr(name MetaName, value lang.Time) {
	txt := value.String()
	inst.buffer.set(name, txt)
}

func (inst *innerMetaSetter) setInt64(name MetaName, value int64) {
	const base = 10
	txt := strconv.FormatInt(value, base)
	inst.buffer.set(name, txt)
}

func (inst *innerMetaSetter) setString(name MetaName, value string) {
	inst.buffer.set(name, value)
}

func (inst *innerMetaSetter) setSum(name MetaName, value *buckets.SUM) {

	if value == nil {
		return
	}

	alg := value.Algorithm
	hex := value.Value
	b := strings.Builder{}

	b.WriteString(alg.String())
	b.WriteRune('(')
	b.WriteString(hex.String())
	b.WriteRune(')')

	str := b.String()
	str = strings.ToLower(str)

	inst.buffer.set(name, str)
}

////////////////////////////////////////////////////////////////////////////////
// EOF
