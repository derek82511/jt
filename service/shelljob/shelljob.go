package shelljob

import (
	"bufio"
	"derek82511/jt/config"
	"derek82511/jt/model"
	"derek82511/jt/service/dataprovider"
	"derek82511/jt/service/log"
	"derek82511/jt/service/web/socket"
	"encoding/json"
	"fmt"
	"os/exec"
	"path"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var (
	runningCmds = make(map[string]*exec.Cmd)
	cmdMu       sync.Mutex
)

func setCmd(jobId string, cmd *exec.Cmd) {
	cmdMu.Lock()
	defer cmdMu.Unlock()
	runningCmds[jobId] = cmd
}

func getCmd(jobId string) *exec.Cmd {
	cmdMu.Lock()
	defer cmdMu.Unlock()
	return runningCmds[jobId]
}

func deleteCmd(jobId string) {
	cmdMu.Lock()
	defer cmdMu.Unlock()
	delete(runningCmds, jobId)
}

func writeMsg(jobId string, event string, msg string) {
	conn := socket.GetConnection(jobId)
	if conn != nil {
		log.WebLogger.Info("write to ws: ", jobId, " ", event, " ", msg)
		resp := struct {
			Event   string `json:"event"`
			Payload string `json:"payload"`
		}{
			Event:   event,
			Payload: msg,
		}
		respBytes, _ := json.Marshal(resp)
		err := conn.WriteMessage(websocket.TextMessage, respBytes)
		if err != nil {
			log.WebLogger.Error("write to ws error: ", err)
			socket.RemoveConnection(jobId)
		}
	}
}

func Execute(job *model.Job) {
	conn := socket.GetConnection(job.ID)
	if conn == nil {
		log.WebLogger.Info("ws conn is nil, stop to execute job ", job.ID)
		return
	}

	scenario := &model.Scenario{}
	if err := dataprovider.GetInstance().Where("ID = ?", job.ScenarioID).First(scenario).Error; err != nil {
		writeMsg(job.ID, "console", "find scenario failed: "+err.Error())
		dataprovider.GetInstance().Model(&model.Job{}).Where("ID = ?", job.ID).Update("status", 2)
		writeMsg(job.ID, "finish", "ok")
		return
	}

	writeMsg(job.ID, "console", "start running job "+job.ScenarioName)
	dataprovider.GetInstance().Model(&model.Job{}).Where("ID = ?", job.ID).Update("status", 1)

	// Prepare arguments for run.sh
	jvmArgs := fmt.Sprintf("JVM_ARGS=\"-Xms%s -Xmx%s\"", job.MinHeap, job.MaxHeap)
	scriptName := scenario.ScriptName
	jmeterRootFolder := config.JMETER_ROOT_FOLDER
	reportFolderName := job.CreateTime.Format("20060102150405") + "_" + scenario.Name + "_report"

	var configMap map[string]string
	configArgs := ""
	if err := json.Unmarshal([]byte(job.Config), &configMap); err == nil {
		propSyntax := "-J"
		if job.ExecuteType == model.Job_ExecuteType_Remote {
			propSyntax = "-G"
		}
		for k, v := range configMap {
			configArgs += fmt.Sprintf(" %s%s=%s", propSyntax, k, v)
		}
	}

	remoteHostConfig := ""
	if job.ExecuteType == model.Job_ExecuteType_Remote {
		remoteHostConfig = "-R " + job.RemoteHost
	}

	cmd := exec.Command("/bin/sh", path.Join(config.JMETER_SHELL_FOLDER, "run.sh"),
		jvmArgs,
		scriptName,
		jmeterRootFolder,
		reportFolderName,
		configArgs,
		remoteHostConfig,
	)

	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	setCmd(job.ID, cmd)

	stdout, _ := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		writeMsg(job.ID, "console", "Failed to start run.sh: "+err.Error())
		dataprovider.GetInstance().Model(&model.Job{}).Where("ID = ?", job.ID).Update("status", 2)
		deleteCmd(job.ID)
		writeMsg(job.ID, "finish", "ok")
		return
	}
	dataprovider.GetInstance().Model(&model.Job{}).Where("ID = ?", job.ID).Update("pid", cmd.Process.Pid)

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		writeMsg(job.ID, "console", m)
	}

	cmd.Wait()
	deleteCmd(job.ID)

	writeMsg(job.ID, "console", "finish job "+job.ScenarioName)
	writeMsg(job.ID, "finish", "ok")
	dataprovider.GetInstance().Model(&model.Job{}).Where("ID = ?", job.ID).Update("status", 2)
}

func Terminate(jobId string) (isRecover bool) {
	var job model.Job
	if err := dataprovider.GetInstance().Where("ID = ?", jobId).First(&job).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			writeMsg(jobId, "console", "job is not exist")
		}
		return false
	}

	if job.Status != 1 {
		writeMsg(jobId, "console", "job is not running")
		return false
	}

	log.WebLogger.Info("terminate job ", job.ID, " pid: ", job.Pid)
	cmd := getCmd(job.ID)
	if cmd != nil && cmd.Process != nil {
		// Kill the entire process group
		if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
			writeMsg(jobId, "console", "Failed to kill process group: "+err.Error())
		} else {
			writeMsg(jobId, "console", "Job process terminated.")
		}
	} else {
		writeMsg(jobId, "console", "Job process not found, could not terminate.")
	}

	dataprovider.GetInstance().Model(&model.Job{}).Where("ID = ?", job.ID).Update("status", 2)
	now := time.Now()
	dataprovider.GetInstance().Model(&model.Job{}).Where("ID = ?", job.ID).Update("finish_time", &now)

	// Since we are now using the robust process group kill,
	// manual recovery is less of an issue. Returning true for compatibility.
	return true
}
