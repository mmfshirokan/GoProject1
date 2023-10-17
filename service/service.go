package service

import "github.com/mmfshirokan/GoProject1/repository"

type Service struct {
	rep *repository.Repository
	err error
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		rep: repo,
	}
}

func (serv *Service) GetUserTroughID(id string) (string, string, error) {
	return serv.rep.GetUserTroughID(id)
}

func (serv *Service) SaveUser(id string, name string, male string) error {
	return serv.rep.SaveUser(id, name, male)
}

func (serv *Service) UpdateUser(id string, name string, male string) error {
	return serv.rep.UpdateUser(id, name, male)
}

func (serv *Service) DeleteUser(id string) error {
	return serv.rep.DeleteUser(id)
}
