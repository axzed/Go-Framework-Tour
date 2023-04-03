package echo

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)

	// Start server
	e.Logger.Fatal(e.Start(":8084"))
}

// Handle
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
