package shelljob

import (
	"bufio"
	"derek82511/jt/config"
	"derek82511/jt/model"
	"derek82511/jt/service/dataprovider"
	"derek82511/jt/service/file"
	"derek82511/jt/service/log"
	"derek82511/jt/service/shelljob/task"
	"derek82511/jt/service/web/socket"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"syscall"
	"unicode/utf8"

	"github.com/Clever/csvlint"
)

func runJobOnError(jobID string, msg string) {
	conn := socket.GetConnection(jobID)

	log.Logger.Error(msg)
	(*conn).Emit("console", msg)

	job := &model.Job{}

	if errs := dataprovider.GetInstance().Where("ID = ?", jobID).First(job).GetErrors(); len(errs) > 0 {
		log.Logger.Error("error: " + errs[0].Error())

		(*conn).Emit("finish", "ok")
		return
	}

	job.Status = model.Job_Status_Completed

	if errs := dataprovider.GetInstance().Save(job).GetErrors(); len(errs) > 0 {
		log.Logger.Error("error: " + errs[0].Error())

		(*conn).Emit("finish", "ok")
		return
	}

	(*conn).Emit("finish", "ok")
}

func DoRunJob(jobID string, args ...string) {
	conn := socket.GetConnection(jobID)

	cmd := exec.Command("/bin/sh", config.JMETER_SHELL_FOLDER+"/run.sh", args[0], args[1], args[2], args[3], args[4], args[5])
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	task.SetCmd(jobID, cmd)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		runJobOnError(jobID, "StdoutPipe error: "+err.Error())
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		runJobOnError(jobID, "StderrPipe error: "+err.Error())
		return
	}

	if err := cmd.Start(); err != nil {
		runJobOnError(jobID, "Start error: "+err.Error())
		return
	}

	scannerOut := bufio.NewScanner(stdout)
	scannerOut.Split(bufio.ScanLines)
	for scannerOut.Scan() {
		text := scannerOut.Text()

		msg := "[INFO] " + text
		log.Logger.Info(msg)
		(*conn).Emit("console", msg)
	}

	scannerErr := bufio.NewScanner(stderr)
	scannerErr.Split(bufio.ScanLines)
	for scannerErr.Scan() {
		text := scannerErr.Text()

		msg := "[ERROR] " + text
		log.Logger.Error(msg)
		(*conn).Emit("console", msg)
	}

	if err := cmd.Wait(); err != nil {
		log.Logger.Error(err.Error())
		(*conn).Emit("console", err.Error())

		RecoveryReport(jobID)
	}

	task.DeleteCmd(jobID)

	job := &model.Job{}

	if errs := dataprovider.GetInstance().Where("ID = ?", jobID).First(job).GetErrors(); len(errs) > 0 {
		msg := "[DATA ERROR] " + errs[0].Error()
		log.Logger.Error(msg)
		(*conn).Emit("console", msg)
		return
	}

	job.Status = model.Job_Status_Completed

	if errs := dataprovider.GetInstance().Save(job).GetErrors(); len(errs) > 0 {
		msg := "[DATA ERROR] " + errs[0].Error()
		log.Logger.Error(msg)
		(*conn).Emit("console", msg)
		return
	}

	endMsg := "[INFO] Job " + job.ID + " finished."
	log.Logger.Info(endMsg)
	(*conn).Emit("console", endMsg)

	(*conn).Emit("finish", "ok")

	socket.ReleaseConnection(job.ID)
}

func RecoveryReport(jobID string) {
	job := &model.Job{}
	scenario := &model.Scenario{}

	if errs := dataprovider.GetInstance().Where("ID = ?", jobID).First(job).GetErrors(); len(errs) > 0 {
		msg := "[DATA ERROR] " + errs[0].Error()
		log.Logger.Error(msg)
		return
	}

	if errs := dataprovider.GetInstance().Where("ID = ?", job.ScenarioID).First(scenario).GetErrors(); len(errs) > 0 {
		msg := "[DATA ERROR] " + errs[0].Error()
		log.Logger.Error(msg)
		return
	}

	conn := socket.GetConnection(jobID)

	reportFolderName := job.CreateTime.Format("20060102150405") + "_" + scenario.Name + "_report"
	jobReportFolder := config.JMETER_REPORTS_FOLDER + "/" + reportFolderName
	logFilePath := jobReportFolder + "/log.csv"
	detailPath := jobReportFolder + "/detail"

	fDetail, err := os.Open(detailPath + "/main.html")
	if err == nil {
		defer fDetail.Close()
		return
	}

	fLog, err := os.Open(logFilePath)
	if err != nil {
		msg := "Log file error: " + err.Error()
		log.Logger.Error(msg)
		if conn != nil {
			(*conn).Emit("console", msg)
		}
	}
	defer fLog.Close()

	comma, _ := utf8.DecodeRuneInString(",")

	errs, _, _ := csvlint.Validate(fLog, comma, false)

	if len(errs) != 0 {
		regex := regexp.MustCompile("#(.*?) ")

		for i := len(errs) - 1; i >= 0; i-- {
			match := regex.FindStringSubmatch(errs[i].Error())
			lineNum, _ := strconv.Atoi(match[1])
			if err := file.RemoveLines(logFilePath, lineNum+1, 1); err != nil {
				msg := "RemoveLines error: " + err.Error()
				log.Logger.Error(msg)
				if conn != nil {
					(*conn).Emit("console", msg)
				}
				return
			}
		}
	}

	cmd := exec.Command("/bin/sh", config.JMETER_SHELL_FOLDER+"/recovery.sh", config.JMETER_ROOT_FOLDER, reportFolderName)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		msg := "StdoutPipe error: " + err.Error()
		log.Logger.Error(msg)
		if conn != nil {
			(*conn).Emit("console", msg)
		}
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		msg := "StderrPipe error: " + err.Error()
		log.Logger.Error(msg)
		if conn != nil {
			(*conn).Emit("console", msg)
		}
		return
	}

	if err := cmd.Start(); err != nil {
		msg := "Start error: " + err.Error()
		log.Logger.Error(msg)
		if conn != nil {
			(*conn).Emit("console", msg)
		}
		return
	}

	scannerOut := bufio.NewScanner(stdout)
	scannerOut.Split(bufio.ScanLines)
	for scannerOut.Scan() {
		text := scannerOut.Text()

		msg := "[INFO] " + text
		log.Logger.Info(msg)
		if conn != nil {
			(*conn).Emit("console", msg)
		}
	}

	scannerErr := bufio.NewScanner(stderr)
	scannerErr.Split(bufio.ScanLines)
	for scannerErr.Scan() {
		text := scannerErr.Text()

		msg := "[ERROR] " + text
		log.Logger.Error(msg)
		if conn != nil {
			(*conn).Emit("console", msg)
		}
	}

	if err := cmd.Wait(); err != nil {
		log.Logger.Error(err.Error())
		if conn != nil {
			(*conn).Emit("console", err.Error())
		}
	}

	recoveryMsg := "Finish recovery report."
	log.Logger.Info(recoveryMsg)
	if conn != nil {
		(*conn).Emit("console", recoveryMsg)
	}
}
