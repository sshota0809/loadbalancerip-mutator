package ip

import "fmt"

type CidrFormatError struct {
	error error
}

func (e *CidrFormatError) Error() string {
	return fmt.Sprintf("CIDR format is invalid. %s", e.error.Error())
}
