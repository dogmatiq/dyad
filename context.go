package dyad

import "fmt"

type cloneContext struct {
	options cloneOptions
}

func (c cloneContext) Errorf(
	format string,
	args ...interface{},
) error {
	return fmt.Errorf(format, args...)
}
