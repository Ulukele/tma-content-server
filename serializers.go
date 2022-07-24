package main

func (s *Server) SerializeUser(model *UserModel) User {
	return User{
		Id: model.Id,
	}
}

func (s *Server) SerializeTeam(model *TeamModel) Team {
	return Team{
		Id:       model.Id,
		Name:     model.Name,
		OwnerId:  model.OwnerId,
		Password: model.Password,
	}
}

func (s *Server) SerializeTeamExtended(model *TeamModel) (TeamExtended, error) {
	usersModels, err := s.contentDBEngine.GetTeamUsersIDs(model.Id)
	if err != nil {
		return TeamExtended{}, err
	}
	users := make([]User, 0)
	for _, user := range usersModels {
		users = append(users, s.SerializeUser(user))
	}
	teamExtended := TeamExtended{
		Users: users,
	}
	teamExtended.Id = model.Id
	teamExtended.Name = model.Name
	teamExtended.OwnerId = model.OwnerId
	teamExtended.Password = model.Password

	return teamExtended, nil
}

func (s *Server) SerializeTask(model *TaskModel) Task {
	return Task{
		Id:         model.Id,
		Title:      model.Title,
		Solved:     model.Solved,
		BoardId:    model.BoardId,
		Importance: model.Importance,
		WorkerId:   model.WorkerId,
	}
}

func (s *Server) SerializeBoard(model *BoardModel) Board {
	return Board{
		Id:     model.Id,
		Name:   model.Name,
		TeamId: model.TeamId,
	}
}

func (s *Server) SerializeBoardExtended(model *BoardModel) (BoardExtended, error) {
	tasksModels, err := s.contentDBEngine.GetBoardTasks(model.Id)
	if err != nil {
		return BoardExtended{}, err
	}
	tasks := make([]Task, 0)
	for _, task := range tasksModels {
		tasks = append(tasks, s.SerializeTask(task))
	}
	teamExtended := BoardExtended{
		Tasks: tasks,
	}
	teamExtended.Id = model.Id
	teamExtended.Name = model.Name
	teamExtended.TeamId = model.TeamId

	return teamExtended, nil
}
