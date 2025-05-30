package handler

import (
	"strconv"

	"blog-andressa/internal/api/router"
	"blog-andressa/internal/model"
	"blog-andressa/internal/service"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (uh *UserHandler) Routes() router.Router {
	return func(route fiber.Router) {
		user := route.Group("users")
		user.Post("/", uh.CreateUserHandler)
		user.Post("/login", uh.Login)
		user.Get("/", uh.FindAllUsersHandler)
		user.Get("/:email", uh.FindUserByEmailHandler)
		user.Put("/:id", uh.UpdateUserHandler)
		user.Delete("/:id", uh.DeleteUserHandler)
	}
}

func (uh *UserHandler) CreateUserHandler(c *fiber.Ctx) error {
	userReq := new(model.UserReq)
	if err := c.BodyParser(userReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.NewErrorResponse(err, fiber.ErrBadRequest))
	}

	res := uh.service.CreateUser(userReq)

	return c.Status(res.Status).JSON(res)
}

func (uh *UserHandler) Login(c *fiber.Ctx) error {
	userReq := new(model.LoginRequest)
	if err := c.BodyParser(userReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.NewErrorResponse(err, fiber.ErrBadRequest))
	}

	res := uh.service.Login(userReq)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": res.Token})
}

func (uh *UserHandler) FindAllUsersHandler(c *fiber.Ctx) error {
	res := uh.service.FindAllUsers()

	return c.Status(res.Status).JSON(res)
}

func (uh *UserHandler) FindUserByEmailHandler(c *fiber.Ctx) error {
	userEmail := c.Params("email")

	res := uh.service.FindUserByEmail(userEmail)

	return c.Status(res.Status).JSON(res)
}

func (uh *UserHandler) UpdateUserHandler(c *fiber.Ctx) error {
	userReq := new(model.UserReq)
	if err := c.BodyParser(userReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.NewErrorResponse(err, fiber.ErrBadRequest))
	}

	userID, _ := strconv.Atoi(c.Params("id", "0"))

	res := uh.service.UpdateUser(uint(userID), userReq)

	return c.Status(res.Status).JSON(res)
}

func (uh *UserHandler) DeleteUserHandler(c *fiber.Ctx) error {
	userID, _ := strconv.Atoi(c.Params("id", "0"))

	res := uh.service.DeleteUser(uint(userID))

	return c.Status(res.Status).JSON(res)
}