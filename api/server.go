package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/kevenmarion/backend_master_class/db/sqlc"
)

// Server servers HTTP requests for our banking service.
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(store db.Store) *Server {
	return &Server{
		store:  store,
		router: gin.Default(),
	}
}

func (s *Server) LoadRouters() {
	s.router.POST("/accounts", s.createAccount)
	s.router.GET("/accounts/:id", s.getAccount)
	s.router.GET("/accounts", s.listAccount)

	s.router.POST("/transfers", s.createTransfer)
}

func (s *Server) LoadValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
}

// Start runs the HTTP server on a specific address.
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func erroResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
