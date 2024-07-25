package serviceerrs

import "fmt"

var (
	ErrCannotCreateMessage  = fmt.Errorf("cannot create message")
	ErrMessageNotFound      = fmt.Errorf("message not found")
	ErrCannotGetMessage     = fmt.Errorf("cannot get message")
	ErrCannotProduceMessage = fmt.Errorf("cannot produce message")
)
