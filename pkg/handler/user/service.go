package user

import (
	userModel "github.com/Ayocodes24/GO-Eats/pkg/database/models/user"
	"github.com/Ayocodes24/GO-Eats/pkg/handler"
	userService "github.com/Ayocodes24/GO-Eats/pkg/service/user"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	Serve    *handler.Server
	group    string
	router   *gin.RouterGroup
	service  *userService.UsrService
	validate *validator.Validate
}

func NewUserHandler(
	s *handler.Server,
	groupName string,
	service *userService.UsrService,
	validate *validator.Validate,
) *UserHandler {
	h := &UserHandler{
		Serve:    s,
		group:    groupName,
		router:   &gin.RouterGroup{},
		service:  service,
		validate: validate,
	}
	h.router = h.registerGroup()
	h.routes()
	h.registerValidator()
	return h
}

func (s *UserHandler) registerValidator() {
	_ = s.validate.RegisterValidation("name", userModel.NameValidator)
	_ = s.validate.RegisterValidation("email", userModel.EmailValidator)
}
