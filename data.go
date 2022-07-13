package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBEngine struct {
	DB *gorm.DB
}

type DBConfig struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     string
	SSLMode  string
	Tz       string
}

func NewDBEngine(dbc DBConfig) (*DBEngine, error) {
	//"host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Novosibirsk"
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dbc.Host,
		dbc.User,
		dbc.Password,
		dbc.Name,
		dbc.Port,
		dbc.SSLMode,
		dbc.Tz)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	dbe := &DBEngine{}
	dbe.DB = db
	return dbe, nil
}

func (dbe *DBEngine) CheckSessionId(username string, session string) bool {
	var user = ServiceUserModel{}
	if err := dbe.DB.
		Where("Username = ? AND Session_ID = ?", username, session).
		Take(&user).Error; err != nil {
		return false
	}
	return true
}

func (dbe *DBEngine) GetUserInfo(username string) (*ServiceUserModel, error) {

	user := &ServiceUserModel{}

	if err := dbe.DB.
		Where("Username = ?", username).
		Take(&user).
		Error; err != nil {
		return nil, err
	}

	return user, nil
}
