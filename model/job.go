package model

import "time"

type Job struct {
	ID           string     `gorm:"column:ID;PRIMARY_KEY;AUTO_INCREMENT"`
	ScenarioID   string     `gorm:"column:SCENARIO_ID"`
	ScenarioName string     `gorm:"column:SCENARIO_Name"`
	MinHeap      string     `gorm:"column:MIN_HEAP"`
	MaxHeap      string     `gorm:"column:MAX_HEAP"`
	Config       string     `gorm:"column:CONFIG"`
	Status       int        `gorm:"column:STATUS"`
	ExecuteType  int        `gorm:"column:EXECUTE_TYPE"`
	RemoteHost   string     `gorm:"column:REMOTE_HOST"`
	ReportPath   string     `gorm:"column:REPORT_PATH"`
	CreateTime   *time.Time `gorm:"column:CREATE_TIME"`
}

const (
	Job_Status_Created    = 0
	Job_Status_InProgress = 1
	Job_Status_Completed  = 2

	Job_ExecuteType_Local  = 0
	Job_ExecuteType_Remote = 1
)
