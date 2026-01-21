package main

import (
	"github.com/lomifile/api/config"
	"github.com/lomifile/api/internal/app"
)

func main() {
	c := config.New()
	err := c.ParseConfig()
	if err != nil {
		panic(err)
	}

	app.Start(c)
}
