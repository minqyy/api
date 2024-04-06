package main

import (
	"github.com/minqyy/api/internal/config"
)

func main() {
	c := config.MustLoad()
	_ = c
}
