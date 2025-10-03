package resp

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func OK(c echo.Context, data interface{}) error      { return c.JSON(http.StatusOK, data) }
func Created(c echo.Context, data interface{}) error { return c.JSON(http.StatusCreated, data) }
func Err(c echo.Context, code int, msg string) error {
	return c.JSON(code, map[string]string{"error": msg})
}
func Msg(c echo.Context, msg string) error {
	return c.JSON(http.StatusOK, map[string]string{"message": msg})
}
