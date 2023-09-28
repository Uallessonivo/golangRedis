package main

import (
	"context"
	"fmt"
	"golangRedis/application"
)

func main() {
	app := application.New()

	err := app.Start(context.TODO())
	if err != nil {
		fmt.Println("failed to start application: ", err)
	}
}
