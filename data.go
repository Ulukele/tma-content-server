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

func (dbe *DBEngine) InternalGetUser(username string) (*UserModel, error) {
	user := &UserModel{}
	if err := dbe.DB.
		Where("Username = ?", username).
		Take(&user).
		Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (dbe *DBEngine) GetUser(userId uint) (*UserModel, error) {
	user := &UserModel{}
	if err := dbe.DB.
		Where("Id = ?", userId).
		Take(&user).
		Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (dbe *DBEngine) GetTeamById(teamId uint) (*TeamModel, error) {
	team := &TeamModel{}
	if err := dbe.DB.
		Where("Id = ?", teamId).
		Take(&team).
		Error; err != nil {
		return nil, err
	}

	return team, nil
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

func (dbe *DBEngine) GetTeams(userId uint) ([]TeamModel, error) {

	user, err := dbe.GetUser(userId)
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

func (dbe *DBEngine) GetTeamUsers(teamId uint) ([]*UserModel, error) {

	team, err := dbe.GetTeamById(teamId)
	if err != nil {
		return nil, err
	}

	var users []*UserModel

	if err := dbe.DB.
		Model(&team).
		Association("Users").
		Find(&users); err != nil {
		return nil, err
	}

	return users, nil
}

func (dbe *DBEngine) GetTeam(userId uint, teamId uint) (*TeamModel, error) {
	user, err := dbe.GetUser(userId)
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

func (dbe *DBEngine) CreateTeam(userId uint, teamName string, teamPassword string) (*TeamModel, error) {
	user, err := dbe.GetUser(userId)
	if err != nil {
		return nil, err
	}

	team := &TeamModel{Name: teamName, OwnerId: user.Id, Password: teamPassword}
	if err := dbe.DB.
		Model(&user).
		Association("Teams").
		Append(team); err != nil {
		return nil, err
	}

	return team, nil
}

func (dbe *DBEngine) DeleteTeam(userId uint, teamId uint) (*TeamModel, error) {
	user, err := dbe.GetUser(userId)
	if err != nil {
		return nil, err
	}

	team := &TeamModel{}
	if err := dbe.DB.
		Where("Id = ? AND Owner_id = ?", teamId, user.Id).
		Delete(&team).
		Error; err != nil {
		return nil, err
	}

	return team, nil
}

func (dbe *DBEngine) JoinTeam(userId uint, teamId uint, teamPassword string) (*TeamModel, error) {
	user, err := dbe.GetUser(userId)
	if err != nil {
		return nil, err
	}

	team := &TeamModel{}
	if err := dbe.DB.
		Where("Id = ? AND Password = ?", teamId, teamPassword).
		Take(&team).
		Error; err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	if err := dbe.DB.
		Model(&team).
		Association("Users").
		Append(user); err != nil {
		return nil, err
	}

	return team, nil
}

func (dbe *DBEngine) LeaveTeam(userId uint, teamId uint) (*TeamModel, error) {
	user, err := dbe.GetUser(userId)
	if err != nil {
		return nil, err
	}

	team, err := dbe.GetTeam(userId, teamId)
	if err != nil {
		return nil, err
	}

	if err := dbe.DB.
		Model(&team).
		Association("Users").
		Delete(user); err != nil {
		return nil, err
	}

	return team, nil
}

func (dbe *DBEngine) CreateBoard(userId uint, teamId uint, boardName string) (*BoardModel, error) {
	team, err := dbe.GetTeam(userId, teamId)
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

func (dbe *DBEngine) GetBoards(userId uint, teamId uint) ([]BoardModel, error) {
	team, err := dbe.GetTeam(userId, teamId)
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

func (dbe *DBEngine) GetBoard(userId uint, teamId uint, boardId uint) (*BoardModel, error) {
	team, err := dbe.GetTeam(userId, teamId)
	if err != nil {
		return nil, err
	}

	var boards []BoardModel
	if err := dbe.DB.
		Model(&team).
		Where("Id = ?", boardId).
		Association("Boards").
		Find(&boards); err != nil {
		return nil, err
	}

	if len(boards) < 1 {
		return nil, fmt.Errorf("no such board")
	}

	return &boards[0], nil
}

func (dbe *DBEngine) GetBoardById(boardId uint) (*BoardModel, error) {
	board := &BoardModel{}
	if err := dbe.DB.
		Where("Id = ?", boardId).
		Take(&board).
		Error; err != nil {
		return nil, err
	}

	return board, nil
}

func (dbe *DBEngine) GetBoardTasks(boardId uint) ([]*TaskModel, error) {

	board, err := dbe.GetBoardById(boardId)
	if err != nil {
		return nil, err
	}

	var tasks []*TaskModel

	if err := dbe.DB.
		Model(&board).
		Association("Tasks").
		Find(&tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (dbe *DBEngine) GetTasks(userId uint, teamId uint, boardId uint) ([]TaskModel, error) {
	board, err := dbe.GetBoard(userId, teamId, boardId)
	if err != nil {
		return nil, err
	}

	var tasks []TaskModel
	if err = dbe.DB.
		Model(&board).
		Association("Tasks").
		Find(&tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (dbe *DBEngine) DeleteBoard(userId uint, teamId uint, boardId uint) (*BoardModel, error) {
	team, err := dbe.GetTeam(userId, teamId)
	if err != nil {
		return nil, err
	}

	board := &BoardModel{}
	if err := dbe.DB.
		Where("Id = ? AND Team_id = ? ", boardId, team.Id).
		Delete(&board).
		Error; err != nil {
		return nil, err
	}

	return board, nil
}

func (dbe *DBEngine) CreateTask(userId uint, teamId uint, boardId uint, title string) (*TaskModel, error) {
	board, err := dbe.GetBoard(userId, teamId, boardId)
	if err != nil {
		return nil, err
	}

	task := &TaskModel{Title: title}
	if err := dbe.DB.
		Model(&board).
		Association("Tasks").
		Append(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (dbe *DBEngine) GetTask(userId uint, teamId uint, boardId uint, taskId uint) (*TaskModel, error) {
	board, err := dbe.GetBoard(userId, teamId, boardId)
	if err != nil {
		return nil, err
	}

	var tasks []TaskModel
	if err := dbe.DB.
		Model(&board).
		Where("Id = ?", taskId).
		Association("Tasks").
		Find(&tasks); err != nil {
		return nil, err
	}

	if len(tasks) < 1 {
		return nil, fmt.Errorf("no such task")
	}

	return &tasks[0], nil
}

func (dbe *DBEngine) DeleteTask(userId uint, teamId uint, boardId uint, taskId uint) (*TaskModel, error) {
	board, err := dbe.GetBoard(userId, teamId, boardId)
	if err != nil {
		return nil, err
	}

	task := &TaskModel{}
	if err := dbe.DB.
		Where("Id = ? AND Board_id = ? ", taskId, board.Id).
		Delete(&task).
		Error; err != nil {
		return nil, err
	}

	return task, nil
}
