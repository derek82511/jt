package task

import "os/exec"

var wsJobCmdMap *map[string]*exec.Cmd

func init() {
	wsJobCmdMap = &map[string]*exec.Cmd{}
}

func SetCmd(jobID string, cmd *exec.Cmd) {
	(*wsJobCmdMap)[jobID] = cmd
}

func GetCmd(jobID string) *exec.Cmd {
	return (*wsJobCmdMap)[jobID]
}

func DeleteCmd(jobID string) {
	cmd := GetCmd(jobID)
	if cmd != nil {
		delete((*wsJobCmdMap), jobID)
	}
}
