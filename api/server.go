package api

import (
	db "github.com/CrunchyBlue/Golang-Bank/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", validateCurrency)
		if err != nil {
			return nil
		}
	}

	// User
	router.POST("/user", server.createUser)

	// Account
	router.GET("/accounts", server.getAccounts)
	router.GET("/account/:id", server.getAccount)
	router.POST("/account", server.createAccount)
	router.PUT("/account/:id", server.updateAccount)
	router.PUT("/account/:id/balance", server.updateAccountBalance)
	router.DELETE("/account/:id", server.deleteAccount)

	// Entry
	router.GET("/entries", server.getEntries)
	router.GET("/entries/:account_id", server.getEntriesForAccount)
	router.GET("/entry/:id", server.getEntry)
	router.POST("/entry", server.createEntry)
	router.PUT("/entry/:id", server.updateEntry)
	router.DELETE("/entry/:id", server.deleteEntry)

	// Transfer
	router.GET("/transfers", server.getTransfers)
	router.GET("/transfers/:account_id/outbound", server.getOutboundTransfersForAccount)
	router.GET("/transfers/:account_id/inbound", server.getInboundTransfersForAccount)
	router.GET("/transfer/:id", server.getTransfer)
	router.POST("/transfer", server.createTransfer)
	router.PUT("/transfer/:id", server.updateTransfer)
	router.DELETE("/transfer/:id", server.deleteTransfer)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
