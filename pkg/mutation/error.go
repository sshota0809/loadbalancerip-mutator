package mutation

import (
	"fmt"
)

type NoAvailableIPError struct{}

func (e *NoAvailableIPError) Error() string {
	return fmt.Sprintf("AvailableIP is not found.")
}
