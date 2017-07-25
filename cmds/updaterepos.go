package cmds

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/bruschill/mut/lib/repo"
	"github.com/fatih/color"
)

//UpdateRepos updates all repos in the $MBC_WORK_ROOT directory
func UpdateRepos() {
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

	//create color printing functions
	var successString = color.New(color.FgGreen).SprintFunc()
	var unchangedString = color.New(color.FgWhite).SprintFunc()
	var errorString = color.New(color.FgRed, color.Bold).SprintFunc()

	//iterate through all repo dirs and update them
	for _, repoDir := range repoDirs {
		wg.Add(1)
		go func(repoDirName string) {
			r := repo.NewRepo(repoDirName, rootPath)
			updateStatus := r.Update()

			if updateStatus.Success {
				if strings.Contains(updateStatus.Message, "updated") {
					fmt.Printf("%s: %s\n", r.Name, successString(updateStatus.Message))
				} else {
					fmt.Printf("%s: %s\n", r.Name, unchangedString(updateStatus.Message))
				}
			} else {
				fmt.Printf("%s: %s\n", r.Name, errorString(updateStatus.Message))
			}

			wg.Done()
		}(repoDir)
	}

	//wait until all goroutines have finished
	wg.Wait()
}
