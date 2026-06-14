package notification

import (
	"net/http"

	"github.com/Ayocodes24/GO-Eats/pkg/handler"
	"github.com/Ayocodes24/GO-Eats/pkg/service/notification"
	"github.com/Ayocodes24/GO-Eats/pkg/wsclients"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
)

type NotifyHandler struct {
	serve             *handler.Server
	group             string
	middlewareGuarded gin.IRoutes
	router            gin.IRoutes
	service           *notification.NotificationService
	middleware        []gin.HandlerFunc
	validate          *validator.Validate
	ws                *websocket.Upgrader
	clients           *wsclients.Registry
}

func NewNotifyHandler(s *handler.Server, group string,
	service *notification.NotificationService, middleware []gin.HandlerFunc,
	validate *validator.Validate, clients *wsclients.Registry) {

	var ws = &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	notifyHandler := &NotifyHandler{
		serve:             s,
		group:             group,
		service:           service,
		middleware:        middleware,
		validate:          validate,
		ws:                ws,
		clients:           clients,
	}
	notifyHandler.middlewareGuarded = notifyHandler.registerMiddlewareGroup(middleware...)
	notifyHandler.router = notifyHandler.registerGroup()
	notifyHandler.regularRoutes()
	notifyHandler.middlewareRoutes()
	notifyHandler.registerValidator()
}

func (s *NotifyHandler) registerValidator() {}
