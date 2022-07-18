package main

func (s *Server) SerializeUser(model *UserModel) User {
	return User{
		Id:       model.Id,
		Username: model.Username,
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
	usersModels, err := s.contentDBEngine.GetTeamUsers(model.Id)
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
