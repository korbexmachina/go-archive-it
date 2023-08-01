package main

import (
	"log"
	"path/filepath"
	"time"

	"github.com/korbexmachina/go-archive-it/utils"
)



func main() {
	start := time.Now()
	// config_path, err := filepath.Abs("~/.config/go-archive-it/config.yaml")
	configPath, _ := filepath.Abs("./test-conf/config.yaml") // test path

	utils.ConfigExists(configPath)
	config := utils.LoadConfig(configPath)

	for _, path := range config.VaultPath {
		utils.Archive(path, config.ArchivePath, config.ArchiveType)
	}

	elapsed := time.Since(start)
	log.Printf("Archive(s) created in [[ %s ]]", elapsed)
}
