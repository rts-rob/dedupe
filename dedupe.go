package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

func main() {
	basedir := flag.String("basedir", ".", "The base directory to search.")
	flag.Parse()

	files := make(map[string][]string)

	var m sync.Mutex
	var wg sync.WaitGroup

	sem := make(chan struct{}, 12)

	wg.Add(1)
	go parseDir(basedir, files, &wg, &m, sem)
	wg.Wait()

	removeDuplicates(files)
}

func parseDir(dirname *string, files map[string][]string, wg *sync.WaitGroup, m *sync.Mutex, sem chan struct{}) {
	// If we've already reached the limit on the semaphore channel, the next line will block
	// and we won't try to open another
	sem <- struct{}{}
	// Once we've completed this goroutine remove one from the semaphore channel so life can
	// go on.
	defer func() { <-sem }()
	defer wg.Done()

	if contents, err := ioutil.ReadDir(*dirname); err != nil {
		log.Fatalf("Error parsing %s: %s\n", *dirname, err.Error())
	} else {
		for _, fileInfo := range contents {
			isDir := fileInfo.IsDir()
			name := fileInfo.Name()

			// Skip hidden files and directories
			if strings.HasPrefix(name, ".") {
				continue
			}

			fileName := *dirname + "/" + name

			if isDir {
				wg.Add(1)
				go parseDir(&fileName, files, wg, m, sem)
			} else {
				// Skip files without extensions
				if strings.Index(name, ".") == -1 {
					continue
				}

				key := name[:strings.LastIndex(name, ".")]

				// Synchronous Mutex Lock
				m.Lock()
				files[key] = append(files[key], fileName)
				m.Unlock()
				// Synchronous Mutex Unlock
			}
		}
	}
}

func removeDuplicates(hashMap map[string][]string) {
	for _, value := range hashMap {
		if len(value) > 1 && hasRequiredExtensions(value) {
			removeNonDNG(value)
		}
	}
}

func hasRequiredExtensions(list []string) bool {
	dngFound, arwFound := false, false

	for _, filename := range list {
		extension := strings.ToLower(filename[strings.LastIndex(filename, "."):])
		if extension == ".dng" {
			dngFound = true
		} else if extension == ".arw" {
			arwFound = true
		}
	}

	return dngFound && arwFound
}

func removeNonDNG(list []string) {
	for _, filename := range list {
		extension := strings.ToLower(filename[strings.LastIndex(filename, "."):])
		if extension != ".dng" {
			fmt.Printf("Removing %s\n", filename)
			os.Remove(filename)
		}
	}
}
