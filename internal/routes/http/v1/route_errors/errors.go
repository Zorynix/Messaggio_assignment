package routeerrs

import (
	"errors"

	"github.com/labstack/echo/v4"
)

func NewErrorResponse(c echo.Context, errStatus int, message string) {
	err := errors.New(message)
	_, ok := err.(*echo.HTTPError)
	if !ok {
		report := echo.NewHTTPError(errStatus, err.Error())
		_ = c.JSON(errStatus, report)
	}
	c.Error(errors.New("internal server error"))
}
