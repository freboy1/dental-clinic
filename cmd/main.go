package main

import (
	"fmt"

	"dental_clinic/internal/config"
)

func main() {
	cfg := config.LoadConfig()
	fmt.Println(cfg.Port)
	
}
