package main

import (
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	Id       uint `gorm:"primaryKey"`
	Username string
	Password string
	Teams    []*TeamModel `gorm:"many2many:users_teams;"`
}

type TeamModel struct {
	gorm.Model
	Id       uint `gorm:"primaryKey"`
	Name     string
	OwnerId  uint
	Password string
	Users    []*UserModel `gorm:"many2many:users_teams;"`
	Boards   []BoardModel `gorm:"foreignKey:TeamId"`
}

type BoardModel struct {
	gorm.Model
	Id     uint `gorm:"primaryKey"`
	Name   string
	TeamId uint
}

type TaskModel struct {
	gorm.Model
	Id    uint `gorm:"primaryKey"`
	Title string
}

func (dbe *DBEngine) initTables() error {

	if err := dbe.DB.AutoMigrate(&UserModel{}); err != nil {
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
