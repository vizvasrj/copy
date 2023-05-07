package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func copyFile(src, dst string, wg *sync.WaitGroup) {
	defer wg.Done()

	source, err := os.Open(src)
	if err != nil {
		fmt.Printf("Failed to open source file: %v\n", err)
		return
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		fmt.Printf("Failed to create destination file: %v\n", err)
		return
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		fmt.Printf("Failed to copy file: %v\n", err)
		return
	}

	// fmt.Printf("Copied file: %s\n", dst)
}

func copyDirectory(src, dst string, wg *sync.WaitGroup) {
	defer wg.Done()

	files, err := ioutil.ReadDir(src)
	if err != nil {
		fmt.Printf("Failed to read source directory: %v\n", err)
		return
	}

	wg.Add(len(files)) // Increment the wait group counter for all files and directories

	for _, file := range files {
		source := filepath.Join(src, file.Name())
		destination := filepath.Join(dst, file.Name())

		if file.IsDir() {
			err := os.MkdirAll(destination, file.Mode())
			if err != nil {
				fmt.Printf("Failed to create directory: %v\n", err)
				continue
			}
			copyDirectory(source, destination, wg)
		} else {
			go copyFile(source, destination, wg)
		}
	}
}

func main() {
	now := time.Now()
	// Configuration
	sourcePath := "demo_files4/"
	destinationPath := "/home/tmp"

	// Create destination directory if it doesn't exist
	err := os.MkdirAll(destinationPath, 0755)
	if err != nil {
		fmt.Printf("Failed to create destination directory: %v\n", err)
		return
	}

	// Get list of files in source path
	files, err := ioutil.ReadDir(sourcePath)
	if err != nil {
		fmt.Printf("Failed to read source directory: %v\n", err)
		return
	}

	// Create a WaitGroup to wait for all copy operations to finish
	var wg sync.WaitGroup

	wg.Add(len(files)) // Increment the wait group counter for all files and directories

	// Perform the copy operation using multiple goroutines
	for _, file := range files {
		source := filepath.Join(sourcePath, file.Name())
		destination := filepath.Join(destinationPath, file.Name())

		if file.IsDir() {
			err := os.MkdirAll(destination, file.Mode())
			if err != nil {
				fmt.Printf("Failed to create directory: %v\n", err)
				continue
			}
			go copyDirectory(source, destination, &wg)
		} else {
			go copyFile(source, destination, &wg)
		}
	}
	wg.Wait()

	fmt.Println("All files copied. in seconds", time.Since(now))
}
