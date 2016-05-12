package gcmd

import "os/exec"

//GitCommand is the base structure for an executable git command
type GitCommand struct {
	//name of bin
	binName string

	//additional args for call, appended to reqArgs
	addlArgs []string

	//path to exec command on
	execPath string
}

//NewGitCommand returns a new instance of GitCommand
func NewGitCommand(execPath string, addlArgs []string) *GitCommand {
	return &GitCommand{
		binName:  "git",
		addlArgs: addlArgs,
		execPath: execPath,
	}
}

func (gc *GitCommand) reqArgs() []string {
	return []string{"-C", gc.execPath}
}

//Run executes GitCommand based on its members
func (gc *GitCommand) Run() error {
	args := append(gc.reqArgs(), gc.addlArgs...)
	return exec.Command(gc.binName, args...).Run()
}

//Output executes GitCommand based on its members
func (gc *GitCommand) Output() (string, error) {
	args := append(gc.reqArgs(), gc.addlArgs...)
	out, err := exec.Command(gc.binName, args...).Output()

	if err != nil {
		return "", err
	}

	return string(out), nil
}
