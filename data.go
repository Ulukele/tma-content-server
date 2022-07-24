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

func (dbe *DBEngine) GetTeams(userId uint) ([]*TeamModel, error) {

	var teams []*TeamModel
	var relations []*TeamUserRelation

	if err := dbe.DB.
		Where("user_id = ?", userId).
		Find(&relations).
		Error; err != nil {
		return nil, err
	}

	for _, relation := range relations {
		team, err := dbe.GetTeamById(relation.TeamId)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}

	return teams, nil
}

func (dbe *DBEngine) GetTeamUsersIDs(teamId uint) ([]*UserModel, error) {

	var users []*UserModel
	var relations []*TeamUserRelation

	if err := dbe.DB.
		Where("team_id = ?", teamId).
		Find(&relations).
		Error; err != nil {
		return nil, err
	}

	for _, relation := range relations {
		users = append(users, &UserModel{Id: relation.UserId})
	}

	return users, nil
}

func (dbe *DBEngine) GetTeam(userId uint, teamId uint) (*TeamModel, error) {
	var relation *TeamUserRelation

	if err := dbe.DB.
		Where("user_id = ? AND team_id = ?", userId, teamId).
		Take(&relation).
		Error; err != nil {
		return nil, err
	}

	team, err := dbe.GetTeamById(relation.TeamId)
	if err != nil {
		return nil, err
	}

	return team, nil
}

func (dbe *DBEngine) CreateTeam(userId uint, teamName string, teamPassword string) (*TeamModel, error) {
	team := &TeamModel{Name: teamName, Password: teamPassword, OwnerId: userId}
	if err := dbe.DB.Create(team).Error; err != nil {
		return nil, err
	}
	relation := &TeamUserRelation{TeamId: team.Id, UserId: userId}
	if err := dbe.DB.Create(relation).Error; err != nil {
		return nil, err
	}
	return team, nil
}

func (dbe *DBEngine) DeleteTeam(userId uint, teamId uint) (*TeamModel, error) {
	team := &TeamModel{}
	var relations []*TeamUserRelation

	// delete team and all relations
	err := dbe.DB.Transaction(func(tx *gorm.DB) error {
		var exists bool
		if err := dbe.DB.
			Model(&team).
			Select("count(*) > 0").
			Where("Id = ? AND Owner_id = ?", teamId, userId).
			Find(&exists).
			Error; err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("can't find command with teamId=%d and ownerId=%d", teamId, userId)
		}

		if err := dbe.DB.
			Where("Id = ? AND Owner_id = ?", teamId, userId).
			Delete(&team).
			Error; err != nil {
			return err
		}

		if err := dbe.DB.
			Where("Team_id = ?", teamId).
			Delete(&relations).
			Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return team, nil
}

func (dbe *DBEngine) JoinTeam(userId uint, teamId uint, teamPassword string) (*TeamModel, error) {

	team := &TeamModel{}
	if err := dbe.DB.
		Where("Id = ? AND Password = ?", teamId, teamPassword).
		Take(&team).
		Error; err != nil {
		return nil, err
	}

	relation := &TeamUserRelation{TeamId: team.Id, UserId: userId}
	if err := dbe.DB.Create(relation).Error; err != nil {
		return nil, err
	}

	return team, nil
}

func (dbe *DBEngine) LeaveTeam(userId uint, teamId uint) (*TeamModel, error) {

	team, err := dbe.GetTeamById(teamId)
	if err != nil {
		return nil, err
	}

	relation := &TeamUserRelation{}
	if err := dbe.DB.
		Where("team_id = ? AND user_id = ?", teamId, userId).
		Delete(&relation).
		Error; err != nil {
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

func (dbe *DBEngine) CreateTask(userId uint, teamId uint, boardId uint, task Task) (*TaskModel, error) {
	board, err := dbe.GetBoard(userId, teamId, boardId)
	if err != nil {
		return nil, err
	}

	if task.WorkerId != 0 {
		var relation *TeamUserRelation
		var exists bool
		if err := dbe.DB.
			Model(&relation).
			Select("count(*) > 0").
			Where("user_id = ? AND team_id = ?", task.WorkerId, teamId).
			Find(&exists).
			Error; err != nil {
			return nil, err
		}
		if !exists {
			return nil, fmt.Errorf("can't assign on such user")
		}
	}

	taskModel := &TaskModel{Title: task.Title, Importance: task.Importance, WorkerId: task.WorkerId}
	if err := dbe.DB.
		Model(&board).
		Association("Tasks").
		Append(taskModel); err != nil {
		return nil, err
	}

	return taskModel, nil
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
