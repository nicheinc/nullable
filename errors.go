package nullable

import "fmt"

// ScanTypeError indicates a type error while trying to Scan Src into Dest.
type ScanTypeError struct {
	Src  interface{}
	Dest interface{}
}

func (e ScanTypeError) Error() string {
	return fmt.Sprintf("cannot scan %v (type %T) into %T", e.Src, e.Src, e.Dest)
}
