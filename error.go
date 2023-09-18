package swaggerui

import "fmt"

var defaultSetupErrorMsgTmpl = "setting up swagger-ui: %s"

type SetupError struct {
	Cause error
}

func (s SetupError) Error() string {
	return fmt.Sprintf(defaultSetupErrorMsgTmpl, s.Cause)
}
