package service

import (
	"blog-andressa/internal/model"
	"blog-andressa/internal/repository"
	"blog-andressa/utils"
	"blog-andressa/utils/auth"
)

type UserService interface {
	CreateUser(userReq *model.UserReq) *model.Response
	Login(req *model.LoginRequest) *model.LoginResponse
	FindAllUsers() *model.Response
	FindUserByEmail(email string) *model.Response
	UpdateUser(id uint, userReq *model.UserReq) *model.Response
	DeleteUser(id uint) *model.Response
}

type UserServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &UserServiceImpl{
		repo: repo,
	}
}

func (us *UserServiceImpl) CreateUser(userReq *model.UserReq) *model.Response {
	user := userReq.ToUser()

	hash , err := utils.HashPassword(user.Password)
	if err != nil {
		return model.NewErrorResponse(err,"Error ao gerar o hash da senha.")
	}

	user.Password = hash

	createUser, err := us.repo.Create(user)
	if err != nil {
		return model.NewErrorResponse(err, 500)
	}

	return model.NewCreatedResponse(createUser.ToUserRes())
}

func (us *UserServiceImpl) Login(req *model.LoginRequest) *model.LoginResponse {
	user, _ := us.repo.FindByEmail(req.Email)

	if !utils.CompareHash(req.Password, user.Password) {
		return nil
	}

	token, err := auth.GenerateToken(user.Name)
	if err != nil {
		return nil
	}

	return &model.LoginResponse{Token: token} 
}

func (us *UserServiceImpl) FindAllUsers() *model.Response {
	users, err := us.repo.FindAll()
	if err != nil {
		return model.NewErrorResponse(err, 404)
	}

	usersResponse := []*model.UserRes{}
	for _, u := range users {
		usersResponse = append(usersResponse, u.ToUserRes())
	}

	return model.NewSuccessResponse(usersResponse)
}

func (us *UserServiceImpl) FindUserByEmail(email string) *model.Response {
	user, err := us.repo.FindByEmail(email)
	if err != nil {
		return model.NewErrorResponse(err, 404)
	}

	return model.NewSuccessResponse(user.ToUserRes())
}

func (us *UserServiceImpl) UpdateUser(id uint, userReq *model.UserReq) *model.Response {
	user, err := us.repo.FindByID(id)
	if err != nil {
		return model.NewErrorResponse(err, 404)
	}

	user.Name = userReq.Name
	user.Email = userReq.Email

	updateUser, err := us.repo.Update(user)
	if err != nil {
		return model.NewErrorResponse(err, 500)
	}

	return model.NewSuccessResponse(updateUser.ToUserRes())
}

func (us *UserServiceImpl) DeleteUser(id uint) *model.Response {
	userID, err := us.repo.FindByID(id)
	if err != nil {
		return model.NewErrorResponse(err, 404)
	}
	_, err = us.repo.Delete(userID.ID)
	if err != nil {
		return model.NewErrorResponse(err, 500)
	}

	return model.NewSuccessResponse(nil)
}