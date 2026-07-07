package userserver

import (
	"context"

	"github.com/Ayocodes24/GO-Eats/pkg/auth"
	userpb "github.com/Ayocodes24/GO-Eats/pkg/proto/user"
)

type Server struct {
	userpb.UnimplementedUserServiceServer
}

func New() *Server {
	return &Server{}
}

// ValidateToken validates a JWT and returns the user_id if valid.
// Called by Order Service on every authenticated request.
func (s *Server) ValidateToken(_ context.Context, req *userpb.ValidateTokenRequest) (*userpb.ValidateTokenResponse, error) {
	valid, userID := auth.ValidateToken(req.Token)
	return &userpb.ValidateTokenResponse{Valid: valid, UserId: userID}, nil
}
