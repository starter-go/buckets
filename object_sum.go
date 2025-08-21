package buckets

import (
	"strings"

	"github.com/starter-go/base/lang"
)

type CheckSumAlgorithm string

const (
	AlgorithmMD5    CheckSumAlgorithm = "MD5"
	AlgorithmSHA1   CheckSumAlgorithm = "SHA1"
	AlgorithmSHA256 CheckSumAlgorithm = "SHA256"
)

func (a CheckSumAlgorithm) String() string {
	return string(a)
}

////////////////////////////////////////////////////////////////////////////////

type SUM struct {
	Algorithm CheckSumAlgorithm `json:"algorithm"`
	Value     lang.Hex          `json:"sum"`
}

func (s SUM) String() string {
	b := strings.Builder{}
	b.WriteString("[SUM")

	b.WriteString(" algorithm:")
	b.WriteString(s.Algorithm.String())

	b.WriteString(" hash:")
	b.WriteString(s.Value.String())

	b.WriteString("]")
	return b.String()
}
