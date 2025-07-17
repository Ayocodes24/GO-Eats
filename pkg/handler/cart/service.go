package cart

import (
	"github.com/Ayocodes24/GO-Eats/pkg/handler"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CartHandler struct {
	serve      *handler.Server
	group      string
	router     gin.IRoutes
	service    *cart_order.CartService
	middleware []gin.HandlerFunc
	validate   *validator.Validate
}

func NewCartHandler(s *handler.Server, groupName string,
	service *cart_order.CartService, middleware []gin.HandlerFunc,
	validate *validator.Validate) {

	cartHandler := &CartHandler{
		s,
		groupName,
		nil,
		service,
		middleware,
		validate,
	}
	cartHandler.router = cartHandler.registerGroup(middleware...)
	cartHandler.routes()
	cartHandler.registerValidator()
}

func (s *CartHandler) registerValidator() {

}
