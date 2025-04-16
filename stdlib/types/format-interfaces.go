package stdlib_types

import "io"

type DecodeMather interface {
	Decode(reader io.Reader, output any) error
	LooselyCompare(a any, b any) (bool, error)
	StrictlyCompare(a any, b any) (bool, error)
}
