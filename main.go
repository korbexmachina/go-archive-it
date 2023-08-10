package main

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/korbexmachina/go-archive-it/utils"
)

func main() {
	start := time.Now()
	configDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to resolve user config directory: %s", err)
	}
	configDir = filepath.Join(configDir, ".config")
	// log.Fatalf(configDir)

	configPath := filepath.Join(configDir, "go-archive-it/config.yaml") // Production path
	// configPath, _ := filepath.Abs("./test-conf/config.yaml") // test path

	utils.ConfigExists(configPath)
	config := utils.LoadConfig(configPath)

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	dir := usr.HomeDir

	archivePath := config.ArchivePath
	if archivePath == "~" {
		archivePath = dir
	} else if strings.HasPrefix(archivePath, "~/") {
		archivePath = filepath.Join(dir, archivePath[2:])
	}

	var wg sync.WaitGroup
	// The loop that actually runs everything
	for _, path := range config.VaultPath {
		// Tilda expansion
		if path == "~" {
			path = dir
		} else if strings.HasPrefix(path, "~/") {
			path = filepath.Join(dir, path[2:])
		}

		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			utils.Archive(path, archivePath, config.ArchiveType)
		}(path)

		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			err := utils.Cleanup(filepath.Join(archivePath, filepath.Base(path)), config.Retention)
			if err != nil {
				log.Fatalf("Failed to cleanup: %s", err)
			}
		}(path)
	}

	wg.Wait()
	elapsed := time.Since(start)
	log.Printf("Archive(s) created in [[ %f ]] seconds", elapsed.Seconds())
}
