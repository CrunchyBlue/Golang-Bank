package api

import (
	"fmt"
	db "github.com/CrunchyBlue/Golang-Bank/sqlc"
	"github.com/CrunchyBlue/Golang-Bank/token"
	"github.com/CrunchyBlue/Golang-Bank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(store db.Store, config util.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.AccessTokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", validateCurrency)
		if err != nil {
			return nil, fmt.Errorf("cannot register binding validator: %w", err)
		}
	}

	server.mapRoutes()

	return server, nil
}

func (server *Server) mapRoutes() {
	router := gin.Default()

	// User
	router.POST("/user", server.createUser)
	router.POST("/user/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// Account
	authRoutes.GET("/accounts", server.getAccounts)
	authRoutes.GET("/account/:id", server.getAccount)
	authRoutes.POST("/account", server.createAccount)
	authRoutes.PUT("/account/:id", server.updateAccount)
	authRoutes.PUT("/account/:id/balance", server.updateAccountBalance)
	authRoutes.DELETE("/account/:id", server.deleteAccount)

	// Entry
	authRoutes.GET("/entries", server.getEntries)
	authRoutes.GET("/entries/:account_id", server.getEntriesForAccount)
	authRoutes.GET("/entry/:id", server.getEntry)
	authRoutes.POST("/entry", server.createEntry)
	authRoutes.PUT("/entry/:id", server.updateEntry)
	authRoutes.DELETE("/entry/:id", server.deleteEntry)

	// Transfer
	authRoutes.GET("/transfers", server.getTransfers)
	authRoutes.GET("/transfers/:account_id/outbound", server.getOutboundTransfersForAccount)
	authRoutes.GET("/transfers/:account_id/inbound", server.getInboundTransfersForAccount)
	authRoutes.GET("/transfer/:id", server.getTransfer)
	authRoutes.POST("/transfer", server.createTransfer)
	authRoutes.PUT("/transfer/:id", server.updateTransfer)
	authRoutes.DELETE("/transfer/:id", server.deleteTransfer)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
