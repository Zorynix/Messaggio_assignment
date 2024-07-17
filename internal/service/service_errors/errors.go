package serviceerrs

import "fmt"

var (
	ErrMessageAlreadyExists = fmt.Errorf("message already exists")
	ErrCannotCreateMessage  = fmt.Errorf("cannot create message")
	ErrMessageNotFound      = fmt.Errorf("message not found")
	ErrCannotGetMessage     = fmt.Errorf("cannot get message")
)
