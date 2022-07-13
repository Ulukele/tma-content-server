package main

import (
	"gorm.io/gorm"
	"reflect"
)

type ServiceUserModel struct {
	gorm.Model
	Id       uint `gorm:"primaryKey"`
	Username string
	Teams    []TeamModel
}

type TeamModel struct {
	gorm.Model
	Id     uint `gorm:"primaryKey"`
	Name   string
	Boards []BoardModel
}

type BoardModel struct {
	gorm.Model
	Id    uint `gorm:"primaryKey"`
	Name  string
	tasks []TaskModel
}

type TaskModel struct {
	gorm.Model
	Id    uint `gorm:"primaryKey"`
	Title string
}

func (dbe *DBEngine) initTables() error {
	models := []reflect.Type{
		reflect.TypeOf((*ServiceUserModel)(nil)),
		reflect.TypeOf((*TeamModel)(nil)),
		reflect.TypeOf((*BoardModel)(nil)),
		reflect.TypeOf((*TaskModel)(nil)),
	}

	// not sure if that works
	for _, modelType := range models {
		if err := dbe.DB.AutoMigrate(modelType.Key()); err != nil {
			return err
		}
	}

	return nil
}
