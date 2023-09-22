package router

import (
	"github.com/labstack/echo/v4"
)

var logger echo.Logger

func setupLogger(e *echo.Echo){
	logger = e.Logger
}

func SetRouter(e *echo.Echo) {
	setupLogger(e)

	// GET "/"
	root_get(e)
	// POST "/"
	root_post(e)
}
