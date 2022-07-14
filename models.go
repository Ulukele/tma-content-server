package main

import (
	"gorm.io/gorm"
)

// UserModel -- model from sessions db
type UserModel struct {
	gorm.Model
	id       uint `gorm:"primaryKey"`
	Username string
}

type ServiceUserModel struct {
	gorm.Model
	Id       uint `gorm:"primaryKey"`
	Username string
	Teams    []*TeamModel `gorm:"many2many:user_teams;"`
}

type TeamModel struct {
	gorm.Model
	Id     uint `gorm:"primaryKey"`
	Name   string
	UserId uint
	Users  []*ServiceUserModel `gorm:"many2many:user_teams;"`
}

type BoardModel struct {
	gorm.Model
	Id   uint `gorm:"primaryKey"`
	Name string
}

type TaskModel struct {
	gorm.Model
	Id    uint `gorm:"primaryKey"`
	Title string
}

func (dbe *DBEngine) initTables() error {

	if err := dbe.DB.AutoMigrate(&ServiceUserModel{}); err != nil {
		return err
	} else if err := dbe.DB.AutoMigrate(&TeamModel{}); err != nil {
		return err
	} else if err := dbe.DB.AutoMigrate(&BoardModel{}); err != nil {
		return err
	} else if err := dbe.DB.AutoMigrate(&TaskModel{}); err != nil {
		return err
	}

	return nil
}
