/*
Utilities required for go-archive-it

https://github.com/korbexmachina/go-archive-it
*/
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

/*
Archive takes 3 arguments, and calls the appropriate archive function, passing on it's other arguments.

args:
	vaultPath string: The path to the directory that will be archived
	archivePath string: The name of the directory where all of the archives are to be stored
	archiveType uint8: The type of archive that is to be created
		- 0 = .tar
		- 1 = .tar.gz

This function does not return a value and logs it's own errors.

Archive creates any directories neccesary for it to function.
*/
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
	tw := tar.NewWriter(archive)
	defer tw.Close()

	link, err := isSymlink(vaultPath)
	if err != nil {
		return err
	}

	if link {
		linkPath, err := os.Open(vaultPath)
		if err != nil {
			return err
		}

		// TODO: Handle symlinks

	}

	// Traverse the directory and all of its subdirectories and add each file found to the archive
	err = filepath.Walk(vaultPath,
		func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					err = addFile(tw, path, vaultPath)
					if err != nil {
						return err
					}
				}
			return nil
	})
	if err != nil {
		return err
	}
	return nil
}

/*
gztarArchive takes 2 arguments

```
args:
- vaultPath string: The path to the directory being archived
- archive io.writer: An io.Writer
```

gztarArchive initializes a ne gzip writer, and chains it onto a tar writer by calling tarArchive
*/
func gztarArchive(vaultPath string, archive io.Writer) error {
	gw := gzip.NewWriter(archive)
	defer gw.Close()
	err := tarArchive(vaultPath, gw) // Chaining the writers
	if err != nil {
		return err
	}
	return nil
}

func addFile(tw *tar.Writer, name string, vaultPath string) error {
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

	tarHeader.Name, err = filepath.Rel(vaultPath, name) // Preserving directory structure relative to the directory being archived
	if err != nil {
		return err
	}

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

/*
Cleanup takes 3 arguments and returns an error

args:
archivePath string: The path to the archive that is being cleaned up
retention uint8: The number of archives that should be retained
verbose bool: whether or not the verbose flag was specified

Cleanup returns an error if something goes wrong
*/
func Cleanup(archivePath string, retention uint8, verbose bool) error {
	files, err := os.ReadDir(archivePath)
	if len(files) < int(retention) {
		return nil
	} else {
		if verbose == true {
			log.Printf("Retention cap [[ %d ]] exceeded - Cleaning up %s...", retention, archivePath)
		}
	}
	if err != nil {
		return err
	}

	oldestFile, err := files[0].Info()
	if err != nil {
		return err
	}
	oldest := time.Now()
	for _, file := range files[1:] {
		info, err := file.Info()
		if err != nil {
			return err
		}
		if info.ModTime().Before(oldest) {
			oldestFile = info
			oldest = info.ModTime()
		}
	}

	os.Remove(filepath.Join(archivePath, oldestFile.Name()))
	if verbose == true {
		log.Printf("%s succesfully cleaned up!", archivePath)
	}
	return nil
}

func isSymlink (path string) (bool, error) {
	link := os.ModeSymlink.Perm()

	if link == os.ModeSymlink {
		return true, nil // File is a symlink
	} else {
		return false, nil // File is not a symlink
	}
}
