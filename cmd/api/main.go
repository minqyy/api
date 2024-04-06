package main

import (
	"github.com/minqyy/api/internal/config"
	"github.com/minqyy/api/internal/pkg/app"
)

func main() {
	conf := config.MustLoad()
	a := app.New(conf)

	a.Run()
}
