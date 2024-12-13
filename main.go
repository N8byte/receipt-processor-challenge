// Submission to the Fetch Rewards Receipt Processor Challenge, by Nathan Waskiewicz.

package main

import (
	"fetch/receipt-processor/api"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Map of receipts keyed on uuid
var receiptsMap = make(map[uuid.UUID]api.Receipt)

// Server implements ServerInterface in /api/api_server.go
type Server struct{}

// Endpoint for /receipts/process
func (s *Server) PostReceiptsProcess(c *gin.Context) {
	var newReceipt api.Receipt
	var newUuid = uuid.New()

	// Attempt to parse the receipt
	if err := c.BindJSON(&newReceipt); err != nil {
		log.Println(err.Error())
		return
	}

	// Store the new receipt
	receiptsMap[newUuid] = newReceipt

	c.IndentedJSON(http.StatusCreated, gin.H{"id": newUuid})
}

// Endpoint for /receipts/{id}/points
func (s *Server) GetReceiptsIdPoints(c *gin.Context, id string) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var receipt = receiptsMap[parsedId]
	var points = CalculatePoints(receipt)
	
	c.IndentedJSON(http.StatusOK, gin.H{"points": points})
}

// Calculate the points for a given receipt
func CalculatePoints(receipt api.Receipt) int64 {
	return 1234
}

func main() {
	router := gin.Default()

	server := &Server{}
	
	api.RegisterHandlers(router, server)

	router.Run(":8080")
}