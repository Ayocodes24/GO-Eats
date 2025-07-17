package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// its an instance of the main GIN engine that handles all the incoming 	HTTP requests
// If s.group is "/users", that means any route you add to this group will automatically start with /users.
// Gin sets up its internal routing so that when a request arrives at /users/..., it gets dispatched to this subgroup.
func (s *UserHandler) registerGroup() *gin.RouterGroup {
	return s.Serve.Gin.Group(s.group)
}

func (s *UserHandler) routes() http.Handler {
	s.router.POST("/", s.addUser)
	s.router.DELETE("/:id", s.deleteUser)
	s.router.POST("/login", s.loginUser)
	return s.Serve.Gin
}
