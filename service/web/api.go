package web

import (
	"derek82511/jt/config"
	"derek82511/jt/model"
	"derek82511/jt/service/dataprovider"
	"derek82511/jt/service/log"
	"derek82511/jt/service/shelljob"
	"derek82511/jt/service/shelljob/task"
	"derek82511/jt/service/web/socket"
	"encoding/json"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	uuid "github.com/iris-contrib/go.uuid"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
)

func SetupApi(app *iris.Application) {
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{iris.MethodGet, iris.MethodPost, iris.MethodDelete, iris.MethodOptions},
		AllowCredentials: true,
	})

	party := app.Party("/api", crs).AllowMethods(iris.MethodOptions)
	{
		party.Get("/scenario", func(ctx iris.Context) {
			scenarios := &[]*model.Scenario{}

			if errs := dataprovider.GetInstance().Order("CREATE_TIME desc").Find(&scenarios).GetErrors(); len(errs) > 0 {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"error": errs[0].Error(),
				})
				return
			}

			results := &[]map[string]string{}

			for _, scenarioEntity := range *scenarios {
				if scenarioEntity.Status == model.Scenario_Status_Deleted {
					continue
				}

				*results = append(*results, map[string]string{
					"id":         scenarioEntity.ID,
					"name":       scenarioEntity.Name,
					"scriptName": scenarioEntity.ScriptName,
					"createTime": scenarioEntity.CreateTime.Format("2006-01-02 15:04:05"),
				})
			}

			ctx.StatusCode(iris.StatusOK)
			ctx.JSON(results)
		})

		party.Post("/scenario", func(ctx iris.Context) {
			name := ctx.FormValue("name")

			var filename string

			now := time.Now()
			datetime := now.Format("20060102150405")

			maxSize := ctx.Application().ConfigurationReadOnly().GetPostMaxMemory()

			err := ctx.Request().ParseMultipartForm(maxSize)
			if err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.JSON(iris.Map{
					"error": err.Error(),
				})
				return
			}

			form := ctx.Request().MultipartForm

			files := form.File["uploadedFile"]
			failures := 0
			for _, file := range files {
				file.Filename = datetime + "_" + file.Filename
				_, err = saveUploadedFile(file, config.JMETER_SCRIPTS_FOLDER)
				if err != nil {
					failures++
				} else {
					filename = file.Filename
				}
			}

			if err != nil {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"error": err.Error(),
				})
				return
			}

			uuid, err := uuid.NewV4()
			if err != nil {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"error": err.Error(),
				})
				return
			}

			scenario := &model.Scenario{
				ID:         uuid.String(),
				Name:       name,
				ScriptName: filename,
				Status:     model.Scenario_Status_Created,
				CreateTime: &now,
			}

			if errs := dataprovider.GetInstance().Create(scenario).GetErrors(); len(errs) > 0 {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"error": errs[0].Error(),
				})
				return
			}

			ctx.StatusCode(iris.StatusOK)
			ctx.JSON(iris.Map{})
		})

		party.Delete("/scenario/{id}", func(ctx iris.Context) {
			id := ctx.Params().Get("id")

			scenario := &model.Scenario{}

			if errs := dataprovider.GetInstance().Where("ID = ?", id).First(scenario).GetErrors(); len(errs) > 0 {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"error": errs[0].Error(),
				})
				return
			}

			// soft delete

			// if err := os.Remove(filepath.Join(JMETER_SCRIPTS_FOLDER, scenario.ScriptName)); err != nil {
			// 	ctx.StatusCode(iris.StatusBadRequest)
			// 	ctx.JSON(iris.Map{
			// 		"error": err.Error(),
			// 	})
			// 	return
			// }

			scenario.Status = model.Scenario_Status_Deleted

			if errs := dataprovider.GetInstance().Save(scenario).GetErrors(); len(errs) > 0 {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"error": errs[0].Error(),
				})
				return
			}

			ctx.JSON(iris.Map{})
		})

		party.Get("/job", func(ctx iris.Context) {
			jobs := &[]*model.Job{}

			if errs := dataprovider.GetInstance().Order("CREATE_TIME desc").Find(&jobs).GetErrors(); len(errs) > 0 {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"error": errs[0].Error(),
				})
				return
			}

			results := &[]map[string]string{}

			for _, jobEntity := range *jobs {
				*results = append(*results, map[string]string{
					"id":           jobEntity.ID,
					"scenarioId":   jobEntity.ScenarioID,
					"scenarioName": jobEntity.ScenarioName,
					"minHeap":      jobEntity.MinHeap,
					"maxHeap":      jobEntity.MaxHeap,
					"config":       jobEntity.Config,
					"status":       strconv.Itoa(jobEntity.Status),
					"executeType":  strconv.Itoa(jobEntity.ExecuteType),
					"remoteHost":   jobEntity.RemoteHost,
					"reportPath":   jobEntity.ReportPath,
					"createTime":   jobEntity.CreateTime.Format("2006-01-02 15:04:05"),
				})
			}

			ctx.StatusCode(iris.StatusOK)
			ctx.JSON(results)
		})

		party.Get("/job/{id}", func(ctx iris.Context) {
			id := ctx.Params().Get("id")

			job := &model.Job{}
			scenario := &model.Scenario{}

			if errs := dataprovider.GetInstance().Where("ID = ?", id).First(job).GetErrors(); len(errs) > 0 {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"error": errs[0].Error(),
				})
				return
			}

			if errs := dataprovider.GetInstance().Where("ID = ?", job.ScenarioID).First(scenario).GetErrors(); len(errs) > 0 {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"error": errs[0].Error(),
				})
				return
			}

			ctx.StatusCode(iris.StatusOK)
			ctx.JSON(iris.Map{
				"id":                 job.ID,
				"scenarioId":         job.ScenarioID,
				"scenarioName":       scenario.Name,
				"scenarioScriptName": scenario.ScriptName,
				"minHeap":            job.MinHeap,
				"maxHeap":            job.MaxHeap,
				"config":             job.Config,
				"status":             strconv.Itoa(job.Status),
				"executeType":        strconv.Itoa(job.ExecuteType),
				"remoteHost":         job.RemoteHost,
				"reportPath":         job.ReportPath,
				"createTime":         job.CreateTime.Format("2006-01-02 15:04:05"),
			})
		})

		party.Post("/job", func(ctx iris.Context) {
			jobRq := &model.JobRq{}

			if err := ctx.ReadJSON(jobRq); err != nil {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"error": err.Error(),
				})
				return
			}

			scenario := &model.Scenario{}

			if errs := dataprovider.GetInstance().Where("ID = ?", jobRq.ScenarioID).First(scenario).GetErrors(); len(errs) > 0 {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"error": errs[0].Error(),
				})
				return
			}

			now := time.Now()

			uuid, err := uuid.NewV4()
			if err != nil {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"error": err.Error(),
				})
				return
			}

			config := jobRq.Config
			if config == "" {
				config = "{}"
			}

			remoteHost := ""
			if jobRq.ExecuteType == model.Job_ExecuteType_Remote {
				remoteHost = jobRq.RemoteHost
			}

			job := &model.Job{
				ID:           uuid.String(),
				ScenarioID:   jobRq.ScenarioID,
				ScenarioName: scenario.Name,
				MinHeap:      jobRq.MinHeap,
				MaxHeap:      jobRq.MaxHeap,
				Config:       config,
				Status:       model.Job_Status_Created,
				ExecuteType:  jobRq.ExecuteType,
				RemoteHost:   remoteHost,
				ReportPath:   "/jmeter/reports/" + now.Format("20060102150405") + "_" + scenario.Name + "_report/detail/index.html",
				CreateTime:   &now,
			}

			if errs := dataprovider.GetInstance().Create(job).GetErrors(); len(errs) > 0 {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"error": errs[0].Error(),
				})
				return
			}

			ctx.StatusCode(iris.StatusOK)
			ctx.JSON(iris.Map{
				"id": job.ID,
			})
		})

		party.Post("/job/run/{id}", func(ctx iris.Context) {
			id := ctx.Params().Get("id")

			job := &model.Job{}
			scenario := &model.Scenario{}

			if errs := dataprovider.GetInstance().Where("ID = ?", id).First(job).GetErrors(); len(errs) > 0 {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"error": errs[0].Error(),
				})
				return
			}

			if errs := dataprovider.GetInstance().Where("ID = ?", job.ScenarioID).First(scenario).GetErrors(); len(errs) > 0 {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"error": errs[0].Error(),
				})
				return
			}

			if job.Status != model.Job_Status_Created {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"error": "invalid operation",
				})
				return
			}

			job.Status = model.Job_Status_InProgress

			if errs := dataprovider.GetInstance().Save(job).GetErrors(); len(errs) > 0 {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"error": errs[0].Error(),
				})
				return
			}

			reportFolderName := job.CreateTime.Format("20060102150405") + "_" + scenario.Name + "_report"

			jvmArgs := ""
			if job.MinHeap != "-1" && job.MaxHeap != "-1" {
				jvmArgs = "JVM_ARGS=\"-Xms" + job.MinHeap + " -Xmx" + job.MaxHeap + "\""
			}

			var configMap map[string]string

			if err := json.Unmarshal([]byte(job.Config), &configMap); err != nil {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"error": err.Error(),
				})
				return
			}

			configArgs := ""

			propertySyntax := " -J"
			remoteHostConfig := ""
			if job.ExecuteType == model.Job_ExecuteType_Remote {
				propertySyntax = " -G"
				remoteHostConfig = "-R " + job.RemoteHost
			}

			for k, v := range configMap {
				configArgs += propertySyntax + k + "=" + v
			}

			go shelljob.DoRunJob(job.ID, jvmArgs, scenario.ScriptName, config.JMETER_ROOT_FOLDER, reportFolderName, configArgs, remoteHostConfig)

			ctx.StatusCode(iris.StatusOK)
			ctx.JSON(iris.Map{})
		})

		party.Post("/job/terminate/{id}", func(ctx iris.Context) {
			id := ctx.Params().Get("id")

			job := &model.Job{}

			if errs := dataprovider.GetInstance().Where("ID = ?", id).First(job).GetErrors(); len(errs) > 0 {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"error": errs[0].Error(),
				})
				return
			}

			if job.Status != model.Job_Status_InProgress {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"error": "invalid operation",
				})
				return
			}

			cmd := task.GetCmd(job.ID)

			if cmd == nil {
				shelljob.RecoveryReport(job.ID)

				job.Status = model.Job_Status_Completed

				if errs := dataprovider.GetInstance().Save(job).GetErrors(); len(errs) > 0 {
					ctx.StatusCode(iris.StatusBadRequest)
					ctx.JSON(iris.Map{
						"error": errs[0].Error(),
					})
					return
				}

				ctx.StatusCode(iris.StatusOK)
				ctx.JSON(iris.Map{
					"recovery": "1",
				})
				return
			}

			conn := socket.GetConnection(job.ID)

			if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
				shelljob.RecoveryReport(job.ID)

				job.Status = model.Job_Status_Completed

				if errs := dataprovider.GetInstance().Save(job).GetErrors(); len(errs) > 0 {
					ctx.StatusCode(iris.StatusBadRequest)
					ctx.JSON(iris.Map{
						"error": errs[0].Error(),
					})
					return
				}

				msg := "Job " + job.ID + " terminate failed: " + err.Error()
				log.Logger.Error(msg)
				if conn != nil {
					(*conn).Emit("console", msg)
				}

				ctx.StatusCode(iris.StatusOK)
				ctx.JSON(iris.Map{
					"recovery": "1",
				})
			} else {
				msg := "Job " + job.ID + " send kill request success."
				log.Logger.Info(msg)
				if conn != nil {
					(*conn).Emit("console", msg)
				}

				ctx.StatusCode(iris.StatusOK)
				ctx.JSON(iris.Map{})
			}
		})

	}
}

func saveUploadedFile(fh *multipart.FileHeader, destDirectory string) (int64, error) {
	src, err := fh.Open()
	if err != nil {
		return 0, err
	}
	defer src.Close()

	out, err := os.OpenFile(filepath.Join(destDirectory, fh.Filename),
		os.O_WRONLY|os.O_CREATE, os.FileMode(0666))

	if err != nil {
		return 0, err
	}
	defer out.Close()

	return io.Copy(out, src)
}
