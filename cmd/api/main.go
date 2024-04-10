package main

import (
	_ "github.com/minqyy/api/api"
	"github.com/minqyy/api/internal/config"
	"github.com/minqyy/api/internal/pkg/app"
)

// @title                        Minqyy API server
// @version                      0.1
// @description                  API Server for Minqyy app
// @host                         localhost:6969
// @BasePath                     /api
// @securityDefinitions.apikey   AccessToken
// @in                           header
// @name                         Authorization
func main() {
	conf := config.MustLoad()
	a := app.New(conf)

	a.Run()
}
