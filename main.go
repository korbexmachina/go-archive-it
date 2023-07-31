package main

import (
	"fmt"
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

	vaultPath := config.VaultPath
	archivePath := config.ArchivePath
	archiveType := config.ArchiveType
	
	// fmt.Printf("Archive Path = %s\nArchive Type = %d\n", archivePath, archiveType)



	elapsed := time.Since(start)
	log.Printf("Archive(s) created in [[ %s ]]", elapsed)
}
