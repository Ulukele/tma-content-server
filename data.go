package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
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
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dbc.Host,
		dbc.User,
		dbc.Password,
		dbc.Name,
		dbc.Port,
		dbc.SSLMode,
		dbc.Tz)
	log.Printf("Use config: %s", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	dbe := &DBEngine{}
	dbe.DB = db
	return dbe, nil
}

func (dbe *DBEngine) GetUserInfo(username string, password string) (*UserModel, error) {

	user := &UserModel{}

	if err := dbe.DB.
		Where("Username = ? AND Password = ?", username, password).
		Take(&user).
		Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (dbe *DBEngine) CreateUser(username string, password string) (*UserModel, error) {
	var exists bool
	user := &UserModel{}
	if err := dbe.DB.Model(&user).
		Select("count(*) > 0").
		Where("username = ?", username).
		Find(&exists).
		Error; err != nil {
		return nil, err

	}
	if exists {
		return nil, fmt.Errorf("user with username = %s already exists", username)
	} else {
		user.Username = username
		user.Password = password
		if err := dbe.DB.Create(user).Error; err != nil {
			return nil, err
		}
		return user, nil
	}
}

func (dbe *DBEngine) GetTeams(username string) ([]TeamModel, error) {

	user := &UserModel{}

	if err := dbe.DB.
		Where("Username = ?", username).
		Take(&user).
		Error; err != nil {
		return nil, err
	}

	var teams []TeamModel

	if err := dbe.DB.
		Model(&user).
		Association("Teams").
		Find(&teams); err != nil {
		return nil, err
	}

	return teams, nil
}

func (dbe *DBEngine) GetTeam(username string, teamId uint) (*TeamModel, error) {
	user := &UserModel{}

	if err := dbe.DB.
		Where("Username = ?", username).
		Take(&user).
		Error; err != nil {
		return nil, err
	}

	var teams []TeamModel
	if err := dbe.DB.
		Model(&user).
		Where("Id = ?", teamId).
		Association("Teams").
		Find(&teams); err != nil {
		return nil, err
	}

	if len(teams) < 1 {
		return nil, fmt.Errorf("no such team")
	}

	return &teams[0], nil

}

func (dbe *DBEngine) CreateTeam(username string, teamName string) (*TeamModel, error) {
	user := &UserModel{}

	if err := dbe.DB.
		Where("Username = ?", username).
		Take(&user).
		Error; err != nil {
		return nil, err
	}

	team := &TeamModel{Name: teamName}
	if err := dbe.DB.
		Model(&user).
		Association("Teams").
		Append(team); err != nil {
		return nil, err
	}

	return team, nil

}
