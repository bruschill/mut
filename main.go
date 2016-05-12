package main

import (
	"fmt"
	"github.com/bruschill/mut/lib/repo"
	"log"
	"os"
	"sync"
)

func updateRepos() {
	//get $MBC_WORK_ROOT env var
	rootPath := os.ExpandEnv("$MBC_WORK_ROOT")
	if rootPath == "" {
		log.Fatalln("$MBC_WORK_ROOT must be defined.")
	}

	//open dir defined by rootPath
	rootDir, err := os.Open(rootPath)
	defer rootDir.Close()
	if err != nil {
		log.Fatal(err)
	}

	//get slice of all dir names from rootDir
	repoDirs, err := rootDir.Readdirnames(0)
	if err != nil {
		log.Fatal(err)
	}

	//initialize WaitGroup
	var wg sync.WaitGroup

	//iterate through all repo dirs and update them
	for _, repoDir := range repoDirs {
		wg.Add(1)
		go func(repoDirName string) {
			r := repo.NewRepo(repoDirName, rootPath)
			updateStatus := r.Update()

			fmt.Printf("%s: %s\n", r.Name, updateStatus.Message)

			wg.Done()
		}(repoDir)
	}

	//wait until all goroutines have finished
	wg.Wait()
}

func main() {
	updateRepos()
}
