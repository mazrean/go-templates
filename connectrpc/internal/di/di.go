package di

import (
	"fmt"

	"github.com/mazrean/go-templates/connectrpc/internal/config"
	"github.com/mazrean/go-templates/connectrpc/internal/pkg/log"
	"github.com/mazrean/go-templates/connectrpc/internal/router"
	"go.uber.org/dig"
)

type App struct {
	config *config.Config
	router *router.Router
}

func (a *App) Run() error {
	return a.router.Run(a.config.Addr)
}

func DI() (*App, error) {
	c := dig.New()

	if err := c.Provide(config.NewConfig); err != nil {
		return nil, fmt.Errorf("failed to provide config: %w", err)
	}

	if err := routerDI(c); err != nil {
		return nil, fmt.Errorf("failed to inject router: %w", err)
	}

	var app *App
	err := c.Invoke(func(c *config.Config, r *router.Router) {
		log.Setup(c.Debug)

		app = &App{
			config: c,
			router: r,
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
