package main

// UserModel shadow model -- don't exist in db
type UserModel struct {
	Id uint
}

type TeamModel struct {
	Id       uint `gorm:"primaryKey"`
	Name     string
	OwnerId  uint
	Password string
	Boards   []BoardModel `gorm:"foreignKey:TeamId"`
}

type TeamUserRelation struct {
	Id     uint `gorm:"primaryKey"`
	TeamId uint
	UserId uint
}

type BoardModel struct {
	Id     uint `gorm:"primaryKey"`
	Name   string
	TeamId uint
	Tasks  []TaskModel `gorm:"foreignKey:BoardId"`
}

type TaskModel struct {
	Id      uint `gorm:"primaryKey"`
	Title   string
	Solved  bool
	BoardId uint
}

func (dbe *DBEngine) initTables() error {

	if err := dbe.DB.AutoMigrate(&TeamModel{}); err != nil {
		return err
	} else if err := dbe.DB.AutoMigrate(&TeamUserRelation{}); err != nil {
		return err
	} else if err := dbe.DB.AutoMigrate(&BoardModel{}); err != nil {
		return err
	} else if err := dbe.DB.AutoMigrate(&TaskModel{}); err != nil {
		return err
	}

	return nil
}
