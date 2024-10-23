package main

import "github.com/mazrean/go-templates/connectrpc/internal/di"

func main() {
	app, err := di.DI()
	if err != nil {
		panic(err)
	}

	if err := app.Run(":8080"); err != nil {
		panic(err)
	}
}
