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

func (dbe *DBEngine) GetUser(username string, password string) (*UserModel, error) {

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

func (dbe *DBEngine) getUserByUsername(username string) (*UserModel, error) {
	user := &UserModel{}

	if err := dbe.DB.
		Where("Username = ?", username).
		Take(&user).
		Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (dbe *DBEngine) GetTeams(username string) ([]TeamModel, error) {

	user, err := dbe.getUserByUsername(username)
	if err != nil {
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
	user, err := dbe.getUserByUsername(username)
	if err != nil {
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
	user, err := dbe.getUserByUsername(username)
	if err != nil {
		return nil, err
	}

	team := &TeamModel{Name: teamName, OwnerId: user.Id}
	if err := dbe.DB.
		Model(&user).
		Association("Teams").
		Append(team); err != nil {
		return nil, err
	}

	return team, nil
}

func (dbe *DBEngine) DeleteTeam(username string, teamId uint) (*TeamModel, error) {
	user, err := dbe.getUserByUsername(username)
	if err != nil {
		return nil, err
	}

	team := &TeamModel{}
	if err := dbe.DB.
		Where("Id = ? Owner_id", teamId, user.Id).
		Delete(&team).
		Error; err != nil {
		return nil, err
	}

	return team, nil
}

func (dbe *DBEngine) CreateBoard(username string, teamId uint, boardName string) (*BoardModel, error) {
	team, err := dbe.GetTeam(username, teamId)
	if err != nil {
		return nil, err
	}

	board := &BoardModel{Name: boardName}
	if err := dbe.DB.
		Model(&team).
		Association("Boards").
		Append(board); err != nil {
		return nil, err
	}

	return board, nil
}

func (dbe *DBEngine) GetBoards(username string, teamId uint) ([]BoardModel, error) {
	team, err := dbe.GetTeam(username, teamId)
	if err != nil {
		return nil, err
	}

	var boards []BoardModel
	if err = dbe.DB.
		Model(&team).
		Association("Boards").
		Find(&boards); err != nil {
		return nil, err
	}

	return boards, nil
}
