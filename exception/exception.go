package exception

import (
	"github.com/pkg/errors"
)

// GLOBAL ERROR EXCEPTION
func GlobalException(message error) error {
	ThrowError := errors.Errorf("Cause and exception + ", message)
	return ThrowError
}
