package main

import (
	"fmt"
	"github.com/minqyy/api/internal/config"
)

func main() {
	_ = config.New()
	fmt.Println("Hello, World!")
}
