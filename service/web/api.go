package web

import (
	"derek82511/jt/config"
	"derek82511/jt/model"
	"derek82511/jt/service/dataprovider"
	"derek82511/jt/service/shelljob"
	"encoding/json"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"

	uuid "github.com/iris-contrib/go.uuid"
	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

func SetupApi(app *iris.Application) {
	party := app.Party("/api")
	party.AllowMethods(iris.MethodOptions) // <<< Allow OPTIONS method for the entire party.
	party.Use(func(ctx iris.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if ctx.Method() == iris.MethodOptions {
			ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
			ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			ctx.StatusCode(iris.StatusNoContent)
			return
		}

		ctx.Next()
	})
	party.Use(iris.Compression)
	{
		party.Get("/scenario", func(ctx iris.Context) {
			scenarios := &[]*model.Scenario{}

			if err := dataprovider.GetInstance().Order("CREATE_TIME desc").Find(&scenarios).Error; err != nil {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"error": err.Error(),
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

			if err := dataprovider.GetInstance().Create(scenario).Error; err != nil {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"error": err.Error(),
				})
				return
			}

			ctx.StatusCode(iris.StatusOK)
			ctx.JSON(iris.Map{})
		})

		party.Delete("/scenario/{id}", func(ctx iris.Context) {
			id := ctx.Params().Get("id")

			scenario := &model.Scenario{}

			if err := dataprovider.GetInstance().Where("ID = ?", id).First(scenario).Error; err != nil {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"error": err.Error(),
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

			if err := dataprovider.GetInstance().Save(scenario).Error; err != nil {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"error": err.Error(),
				})
				return
			}

			ctx.JSON(iris.Map{})
		})

		party.Get("/job", func(ctx iris.Context) {
			jobs := &[]*model.Job{}

			if err := dataprovider.GetInstance().Order("CREATE_TIME desc").Find(&jobs).Error; err != nil {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"error": err.Error(),
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

			if err := dataprovider.GetInstance().Where("ID = ?", id).First(job).Error; err != nil {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"error": err.Error(),
				})
				return
			}

			if err := dataprovider.GetInstance().Where("ID = ?", job.ScenarioID).First(scenario).Error; err != nil {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"error": err.Error(),
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
			var jobRq struct {
				ScenarioID  string            `json:"scenarioId"`
				Config      map[string]string `json:"config"`
				MaxHeap     string            `json:"maxHeap"`
				MinHeap     string            `json:"minHeap"`
				ExecuteType int               `json:"executeType"`
				RemoteHost  string            `json:"remoteHost"`
			}

			if err := ctx.ReadJSON(&jobRq); err != nil {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"error": err.Error(),
				})
				return
			}

			scenario := &model.Scenario{}

			if err := dataprovider.GetInstance().Where("ID = ?", jobRq.ScenarioID).First(scenario).Error; err != nil {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"error": err.Error(),
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

			configBytes, err := json.Marshal(jobRq.Config)
			if err != nil {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"error": err.Error(),
				})
				return
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
				Config:       string(configBytes),
				Status:       model.Job_Status_Created,
				ExecuteType:  jobRq.ExecuteType,
				RemoteHost:   remoteHost,
				ReportPath:   "/jmeter/reports/" + now.Format("20060102150405") + "_" + scenario.Name + "_report/detail/main.html",
				CreateTime:   &now,
			}

			if err := dataprovider.GetInstance().Create(job).Error; err != nil {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"error": err.Error(),
				})
				return
			}

			ctx.StatusCode(iris.StatusOK)
			ctx.JSON(iris.Map{
				"id": job.ID,
			})
		})

		party.Post("/job/run/:id", func(ctx iris.Context) {
			jobId := ctx.Params().Get("id")

			var job model.Job
			if err := dataprovider.GetInstance().Where("ID = ?", jobId).First(&job).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					ctx.StatusCode(iris.StatusNotFound)
					ctx.JSON(iris.Map{
						"error": "job is not exist",
					})
					return
				}
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.JSON(iris.Map{
					"error": err.Error(),
				})
				return
			}

			go shelljob.Execute(&job)

			ctx.StatusCode(iris.StatusOK)
		})

		party.Post("/job/terminate/:id", func(ctx iris.Context) {
			jobId := ctx.Params().Get("id")
			isRecover := shelljob.Terminate(jobId)

			ctx.StatusCode(iris.StatusOK)
			ctx.JSON(iris.Map{
				"recovery": isRecover,
			})
		})

		party.Get("/jobs", func(ctx iris.Context) {
			page, err := ctx.URLParamInt("page")
			if err != nil || page < 1 {
				page = 1
			}

			pageSize, err := ctx.URLParamInt("pageSize")
			if err != nil || pageSize < 1 {
				pageSize = 10
			}

			keyword := ctx.URLParam("keyword")

			var jobs []model.Job
			var count int64

			db := dataprovider.GetInstance().Model(&model.Job{})

			if keyword != "" {
				db = db.Where("SCENARIO_NAME LIKE ?", "%"+keyword+"%")
			}

			db.Count(&count)
			db.Order("CREATE_TIME desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&jobs)

			ctx.JSON(iris.Map{
				"total": count,
				"data":  jobs,
			})
		})

		party.Post("/job/:id/report/recovery", func(ctx iris.Context) {
			// This functionality was removed during refactoring.
			// The new shelljob logic handles report generation automatically.
			ctx.StatusCode(iris.StatusNotImplemented)
			ctx.JSON(iris.Map{
				"message": "Report recovery is no longer a manual operation and has been removed.",
			})
		})

		party.Get("/job/:id/report/main", func(ctx iris.Context) {
			// This functionality was removed during refactoring.
			// Reports are accessed via static file serving.
			ctx.StatusCode(iris.StatusNotImplemented)
			ctx.JSON(iris.Map{
				"message": "This API endpoint has been removed. Reports are served statically.",
			})
		})

		party.Get("/job/:id/report/log", func(ctx iris.Context) {
			// This functionality was removed during refactoring.
			ctx.StatusCode(iris.StatusNotImplemented)
			ctx.JSON(iris.Map{
				"message": "This API endpoint has been removed.",
			})
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
