package pkgerrs

import (
	"fmt"
)

type ErrInvalidArgument struct {
	Argument string
	Reason   string
}

func (e *ErrInvalidArgument) Error() string {
	return fmt.Sprintf("cueconfig: invalid argument %s (%s)", e.Argument, e.Reason)
}
