package buckets

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"strings"

	"github.com/starter-go/base/lang"
)

type CheckSumAlgorithm string

const (
	AlgorithmMD5    CheckSumAlgorithm = "MD5"
	AlgorithmSHA1   CheckSumAlgorithm = "SHA1"
	AlgorithmSHA256 CheckSumAlgorithm = "SHA256"
	AlgorithmSHA512 CheckSumAlgorithm = "SHA512"
)

func (a CheckSumAlgorithm) String() string {
	return string(a)
}

func (a CheckSumAlgorithm) NewHash() (hash.Hash, error) {

	switch a {

	case AlgorithmSHA1:
		return sha1.New(), nil

	case AlgorithmSHA256:
		return sha256.New(), nil

	case AlgorithmSHA512:
		return sha512.New(), nil

	case AlgorithmMD5:
		return md5.New(), nil
	}
	return nil, fmt.Errorf("buckets.CheckSumAlgorithm: Unsupported hash algorithm: %s", a)
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
