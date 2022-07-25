package main

import (
	"context"
	app "main/src/core/app"
)

func main() {
	container := app.BuildContainer()
	err := container.Invoke(func(application *app.Application) {
		ctx := context.Background()
		err := application.Init(ctx)
		if err != nil {
			panic(err)
		}

		err = application.Run(ctx)
		if err != nil {
			panic(err)
		}
	})

	if err != nil {
		panic(err)
	}

}
