package k8s

import (
	"fmt"
)

type NoHomeDirError struct{}

func (e *NoHomeDirError) Error() string {
	return fmt.Sprintf("Home directory is not found.")
}
