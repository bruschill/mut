package repo

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/bruschill/mut/lib/gcmd"
)

//Repo represents a git repository
type Repo struct {
	//name of the repository
	Name string

	//full path to repository
	fullPath string

	//name of branch initially selected
	origBranch string

	//tracks whether or not unstaged changes were stashed
	hasStashedChanges bool
}

//UpdateStatus represents the status of calling Repo.Update()
type UpdateStatus struct {
	Success bool
	Message string
}

//NewRepo returns a new instance of Repo
func NewRepo(dirName string, rootPath string) *Repo {
	repo := &Repo{
		Name:     dirName,
		fullPath: filepath.Join(rootPath, dirName),
	}

	branch, err := repo.currentBranch()
	if err != nil {
		log.Fatal(err)
	}

	repo.origBranch = branch

	return repo
}

//hasUnstagedChanges returns true if the repo has unstaged changes, false otherwise
func (r *Repo) hasUnstagedChanges() bool {
	addlArgs := []string{"diff", "--exit-code"}
	gc := gcmd.NewGitCommand(r.fullPath, addlArgs)

	err := gc.Run()
	if err != nil {
		return true
	}

	return false
}

//stashChanges stashes changes for current branch
func (r *Repo) stashChanges() {
	if r.hasUnstagedChanges() {
		r.hasStashedChanges = true
		addlArgs := []string{"stash"}
		gc := gcmd.NewGitCommand(r.fullPath, addlArgs)

		gc.Run()
	}
}

//retrieveChanges retrieves changes that were previously stashed
func (r *Repo) retrieveChanges() {
	if r.hasStashedChanges {
		addlArgs := []string{"stash", "pop"}
		gc := gcmd.NewGitCommand(r.fullPath, addlArgs)

		gc.Run()
	}
}

//currentBranch returns name of current branch as string
func (r *Repo) currentBranch() (string, error) {
	addlArgs := []string{"rev-parse", "--abbrev-ref", "HEAD"}
	gc := gcmd.NewGitCommand(r.fullPath, addlArgs)

	branchName, err := gc.Output()
	if err != nil {
		return "", err
	}

	return branchName, nil
}

//checkoutBranch checks out branch specified by branchName
func (r *Repo) checkoutBranch(branchName string) {
	curBranch, _ := r.currentBranch()

	if branchName != curBranch {
		addlArgs := []string{"checkout", branchName}
		gc := gcmd.NewGitCommand(r.fullPath, addlArgs)

		gc.Run()
	}
}

//Update updates repo's master branch and returns repo back to original state before the update
func (r *Repo) Update() *UpdateStatus {
	r.stashChanges()
	r.checkoutBranch("master")

	addlArgs := []string{"pull"}
	gc := gcmd.NewGitCommand(r.fullPath, addlArgs)

	out, err := gc.Output()

	if err != nil {
		statusMsg := "An error occurred when attempting to update " + r.Name + "\n"
		return &UpdateStatus{Success: false, Message: statusMsg}
	}

	statusMsg := ""
	if strings.Contains(out, "is up to date") {
		statusMsg = "no changes."
	} else {
		statusMsg = "updated."
	}

	r.checkoutBranch(r.origBranch)
	r.retrieveChanges()

	return &UpdateStatus{Success: true, Message: statusMsg}
}
