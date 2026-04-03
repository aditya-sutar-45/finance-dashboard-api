package main

import (
	"context"
	"fmt"

	"github.com/aditya-sutar-45/finance-dashboard-api/app"
)

func main() {
	app := app.New()
	err := app.Start(context.TODO())
	if err != nil {
		fmt.Println(err)
	}
}
