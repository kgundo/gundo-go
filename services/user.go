package services

import (
	"errors"

	"github.com/kgundo/gundo-go/repositories"
)

type UserService struct {
	CommonService
	userRepo *repositories.UserRepository
}

func InitUserService() *UserService {
	commonService := InitCommonService()
	userRepo := repositories.InitUserRepository()
	return &UserService{CommonService: commonService, userRepo: userRepo}
}

func (s *UserService) CreateUser(user interface{}) (interface{}, interface{}) {
	newUserId, err := s.userRepo.CreateUser(s.db, user)
	if err != nil {
		return nil, err
	}
	return s.userRepo.GetUserByID(s.db, newUserId.(int))
}

func (s *UserService) GetAllUsers(limit int, page int) (interface{}, error) {
	if limit <= 0 {
		limit = 10
	}
	if page < 1 {
		page = 1
	}
	offset := limit * (page - 1)
	return s.userRepo.GetAllUsers(s.db, limit, offset)
}

func (s *UserService) GetUserByID(id int) (interface{}, error) {
	result, err := s.userRepo.GetUserByID(s.db, id)
	if int(result.ID) != id {
		return nil, errors.New("Not Found")
	}
	return result, err
}

func (s *UserService) DeleteUser(id int) error {
	existsUser, err := s.userRepo.GetUserByID(s.db, id)
	if int(existsUser.ID) != id || err != nil {
		return errors.New("Not Found")
	}
	err = s.userRepo.DeleteUser(s.db, id)
	return err
}

func (s *UserService) UpdateUser(id int, user map[string]interface{}) (interface{}, interface{}) {
	existsUser, err := s.userRepo.GetUserByID(s.db, id)
	if int(existsUser.ID) != id || err != nil {
		return nil, errors.New("Not Found")
	}
	result, err := s.userRepo.UpdateUser(s.db, user, id)
	return result, err
}
