package utils

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func Archive(vaultPath string, archivePath string, archiveType uint8) {
	fullPath := filepath.Join(archivePath, filepath.Base(vaultPath)) // Path to subdir in the archive dir
	time := time.Now().Format(time.RFC3339)

	err := os.MkdirAll(fullPath, 0755) // Create the subdir in the archive dir
	if err != nil {
		log.Fatalf("Failed to create archive directory: %s", err)
	}

	switch archiveType {
	case 0:
		fileName := time + ".tar" // Name of the tar file "2006-01-02T15:04:05Z07:00.tar"
		_, err := os.Create(filepath.Join(fullPath, fileName))
		if err != nil {
			log.Fatal(err)
		}
		outfile, err := os.OpenFile(filepath.Join(fullPath, fileName), os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			log.Fatal(err)
		}
		defer outfile.Close()
		err = tarArchive(vaultPath, outfile)
		if err != nil {
			log.Fatalf("Failed to create tar archive: %s", err)
		}
	case 1:
		fileName := time + ".tar.gz" // Name of the gztar file "2006-01-02T15:04:05Z07:00.tar.gz"
		_, err := os.Create(filepath.Join(fullPath, fileName))
		if err != nil {
			log.Fatal(err)
		}
		outfile, err := os.OpenFile(filepath.Join(fullPath, fileName), os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			log.Fatal(err)
		}
		defer outfile.Close()
		err = gztarArchive(vaultPath, outfile)
		if err != nil {
			log.Fatalf("Failed to create gztar archive: %s", err)
		}
	default:
		log.Print("No archive type specified, defaulting to .tar.gz")
		fileName := time + ".tar.gz" // Name of the gztar file "2006-01-02T15:04:05Z07:00.tar.gz"
		_, err := os.Create(filepath.Join(fullPath, fileName))
		if err != nil {
			log.Fatal(err)
		}
		outfile, err := os.OpenFile(filepath.Join(fullPath, fileName), os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			log.Fatal(err)
		}
		defer outfile.Close()
		err = gztarArchive(vaultPath, outfile)
		if err != nil {
			log.Fatalf("Failed to create gztar archive: %s", err)
		}
	}
	return
}

func tarArchive(vaultPath string, archive io.Writer) error {
	entries, err := os.ReadDir(vaultPath) // Read all entries in the target directory
	if err != nil {
		return err
	}

	tw := tar.NewWriter(archive)
	defer tw.Close()

	for _, file := range entries { // FIXME: Creating invalid archives / Not working with sub-directories
		if !file.IsDir() {
			err := addFile(tw, filepath.Join(vaultPath, file.Name()))
			if err != nil {
				return err
			}
		} else {
			err := tarArchive(filepath.Join(vaultPath, file.Name()), archive)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func gztarArchive(vaultPath string, archive io.Writer) error {
	gw := gzip.NewWriter(archive)
	defer gw.Close()
	err := tarArchive(vaultPath, gw) // Chaining the writers
	if err != nil {
		return err
	}
	return nil
}

func addFile(tw *tar.Writer, name string) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	// get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	tarHeader, err := tar.FileInfoHeader(fileInfo, fileInfo.Name())
	if err != nil {
		return err
	}

	tarHeader.Name = name // Preserving directory structure

	err = tw.WriteHeader(tarHeader)
	if err != nil {
		return err
	}

	_, err = io.Copy(tw, file) // Add file contents to archive
	if err != nil {
		return err
	}

	return nil
}
