package di

import (
	"fmt"

	"github.com/mazrean/go-templates/connectrpc/internal/router"
	"go.uber.org/dig"
)

type App struct {
	*router.Router
}

func DI() (*App, error) {
	c := dig.New()

	if err := routerDI(c); err != nil {
		return nil, fmt.Errorf("failed to inject router: %w", err)
	}

	var app *App
	err := c.Invoke(func(r *router.Router) {
		app = &App{
			Router: r,
		}
	})
	if err != nil {
		return nil, fmt.Errorf("failed to inject app: %w", err)
	}

	return app, nil
}

func routerDI(c *dig.Container) error {
	if err := c.Provide(router.NewRouter); err != nil {
		return fmt.Errorf("failed to provide router: %w", err)
	}

	if err := c.Provide(router.NewExample); err != nil {
		return fmt.Errorf("failed to provide example: %w", err)
	}

	return nil
}
