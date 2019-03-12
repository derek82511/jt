package model

import "time"

type Scenario struct {
	ID         string     `gorm:"column:ID;PRIMARY_KEY;AUTO_INCREMENT"`
	Name       string     `gorm:"column:NAME;unique"`
	ScriptName string     `gorm:"column:SCRIPT_NAME;unique"`
	Status     int        `gorm:"column:STATUS"`
	CreateTime *time.Time `gorm:"column:CREATE_TIME"`
}

const (
	Scenario_Status_Created = 0
	Scenario_Status_Deleted = 1
)
