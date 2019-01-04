package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"
)

func main() {
	basedir := flag.String("basedir", ".", "The base directory to search.")
	flag.Parse()
	files := make(map[string][]string)

	var wg sync.WaitGroup
	wg.Add(1)
	go parseDir(basedir, files, &wg)
	wg.Wait()
	fmt.Println("Finished Waiting")
	printMap(files)
}

func parseDir(dirname *string, files map[string][]string, wg *sync.WaitGroup) {
	defer wg.Done()

	if contents, err := ioutil.ReadDir(*dirname); err != nil {
		log.Fatalf("Error parsing %s: %s\n", *dirname, err.Error())
	} else {
		for _, fileInfo := range contents {
			// Skip hidden files and directories
			if fileInfo.Name()[0] == '.' {
				continue
			}

			fileName := *dirname + "/" + fileInfo.Name()

			if fileInfo.IsDir() {
				wg.Add(1)
				go parseDir(&fileName, files, wg)
			} else {
				key := fileInfo.Name()[:strings.LastIndex(fileInfo.Name(), ".")]
				if value, ok := files[key]; !ok {
					var newList []string
					newList = append(newList, fileName)
					files[key] = newList
				} else {
					files[key] = append(value, fileName)
				}
			}
		}
	}
}

func printMap(hashMap map[string][]string) {
	for key, value := range hashMap {
		fmt.Printf("Key: %s\n", key)
		for _, filename := range value {
			fmt.Printf("\tValue: %s\n", filename)
		}
	}
}
