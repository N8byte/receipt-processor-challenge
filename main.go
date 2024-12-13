// Submission to the Fetch Rewards Receipt Processor Challenge, by Nathan Waskiewicz.

package main

import (
	"fetch/receipt-processor/api"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"unicode"

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
	var points int64

	// Rule 1: One point for every alphanumeric character in the retailer name
	for _, char := range receipt.Retailer {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			points += 1
		}
	}
	fmt.Println("1 points now:", points)

	// Rule 2: 50 points if the total is a round dollar amount with no cents
	if math.Mod(receipt.Total, 1) == 0 {
		points += 50
		fmt.Println("2 points now:", points)
	}

	// Rule 3: 25 points if the total is a multiple of 0.25
	if math.Mod(receipt.Total, 0.25) == 0 {
		points += 25
		fmt.Println("3 points now:", points)
	}

	// Rule 4: 5 points for every two items on the receipt
	points += (int64(len(receipt.Items) / 2)) * 5
	fmt.Println("4 points now:", points)

	// Rule 5: If the trimmed length of the item description is a multiple of 3, multiply the price
	// by 0.2 and round up to the nearest integer. The result is the number of points earned
	for _, item := range receipt.Items {
		if len(strings.TrimSpace(item.ShortDescription)) % 3 == 0 {
			points += int64(math.Ceil(item.Price * 0.2))
			fmt.Println("5 points now:", points)
		}
	}

	// Rule 6: 6 points if the day in the purchase date is odd
	if receipt.PurchaseDate.Day() % 2 == 1 {
		points += 6
		fmt.Println("6 points now:", points)
	}

	// Rule 7: 6 points if the day in the purchase date is odd
	var hour, rule6Err = strconv.ParseInt(strings.Split(receipt.PurchaseTime, ":")[0], 10, 64)
	if rule6Err != nil {
		fmt.Println("Error parsing Purchase Time: ", rule6Err.Error())
		return 0
	}
	if hour >= 14 && hour < 16 {
		points += 10
		fmt.Println("7 points now:", points)
	}

	return points
}

func main() {
	router := gin.Default()

	server := &Server{}
	
	api.RegisterHandlers(router, server)

	router.Run(":8080")
}