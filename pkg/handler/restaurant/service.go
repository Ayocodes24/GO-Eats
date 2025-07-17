package restaurant

import (
	"github.com/Ayocodes24/GO-Eats/pkg/handler"
	"github.com/Ayocodes24/GO-Eats/pkg/service/restaurant"
	"github.com/gin-gonic/gin"
)

type RestaurantHandler struct {
	Serve   *handler.Server
	group   string
	router  *gin.RouterGroup
	service *restaurant.RestaurantService
}

func NewRestaurantHandler(s *handler.Server, groupName string, service *restaurant.RestaurantService) {

	restroHandler := &RestaurantHandler{
		s,
		groupName,
		&gin.RouterGroup{},
		service,
	}

	restroHandler.router = restroHandler.registerGroup()
	restroHandler.routes()
}
