package utils

import (
	"bytes"
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io/ioutil"
	"log"
	"os"
)

func Archive(vaultPath string, archivePath string, archiveType uint8) {
	switch archiveType {
	case 0:
		tarArchive(vaultPath, archivePath)
	case 1:
		gztarArchive(vaultPath, archivePath)
	case 2:
		zipArchive(vaultPath, archivePath)
	default:
		log.Print("No archive type specified, defaulting to .tar.gz")
		gztarArchive(vaultPath, archivePath)
	}
}

func tarArchive(vaultPath string, archivePath string) {
	// TODO: Implement
}

func gztarArchive(vaultPath string, archivePath string) {
	// TODO: Implement
}

func zipArchive(vaultPath string, archivePath string) {
	// TODO: Implement
}
