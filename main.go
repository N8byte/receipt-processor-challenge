// Submission to the Fetch Rewards Receipt Processor Challenge, by Nathan Waskiewicz.

package main

import (
	"fetch/receipt-processor/api"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Map of receipts keyed on uuid.
var receiptsMap = make(map[uuid.UUID]api.Receipt)

type Server struct{}

func (s *Server) PostReceiptsProcess(c *gin.Context) {
	var newReceipt api.Receipt
	var newUuid = uuid.New()

	if err := c.BindJSON(&newReceipt); err != nil {
		log.Println(err.Error())
		return
	}

	receiptsMap[newUuid] = newReceipt

	c.IndentedJSON(http.StatusCreated, gin.H{"id": newUuid})
}

func (s *Server) GetReceiptsIdPoints(c *gin.Context, id string) {

}

func main() {
	router := gin.Default()

	server := &Server{}
	
	api.RegisterHandlers(router, server)

	router.Run(":8080")
}