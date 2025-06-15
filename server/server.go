package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/tylerdevx/geetest-solver/server/middlewares"
	"github.com/tylerdevx/geetest-solver/server/routes"
)

func Start(serverHost string, serverPort string) error {
	app := echo.New()

	middlewares.SetupLoggerMiddleware(app)
	middlewares.SetupTimeoutMiddleware(app)

	routes.SetupRoutes(app)

	serverAddress := fmt.Sprintf("%s:%s", serverHost, serverPort)

	if err := app.Start(serverAddress); err != nil {
		return err
	}

	return nil
}
