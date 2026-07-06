package main

import (
	"fmt"

	"github.com/SoFarCalm/blobgator/internal/config"
)

func main() {
	cfg := config.ReadConfig()
	if err := cfg.SetUser("Lonnie"); err != nil {
		fmt.Println("failed to set user:", err)
	}
	cfg = config.ReadConfig()
	fmt.Printf("DbURL: %s, CurrentUser: %s\n", cfg.DbURL, cfg.CurrentUsername)
}
