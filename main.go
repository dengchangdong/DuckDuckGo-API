package main

import (
	"os"
	"strconv"

	"github.com/dengchangdong/DuckDuckGo-API/duckduckgo"
	"github.com/dengchangdong/DuckDuckGo-API/typings"
	"github.com/acheong08/endless"
	"github.com/gin-gonic/gin"
)

func main() {
	HOST := os.Getenv("HOST")
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}
	handler := gin.Default()

	handler.GET("/search/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	handler.GET("/search", func(ctx *gin.Context) {
		chatmodParam := ctx.Query("chatmod")

		// Check the chatmod parameter
		if chatmodParam == "true" {
			handleChatmodTrue(ctx)
		} else {
			handleChatmodFalse(ctx)
		}
	})

	endless.ListenAndServe(HOST+":"+PORT, handler)
}

func handleChatmodTrue(ctx *gin.Context) {
	var search typings.Search
	search.Query = ctx.Query("query")

	if search.Query == "" {
		ctx.JSON(400, gin.H{"error": "Query is required"})
		return
	}

	// Get results
	results, err := duckduckgo.Get_results(search)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Build the response in chatmod format
	var prompt string
	prompt += search.Query + ";"

	for _, result := range results {
		prompt += result.Title + ";" + result.Snippet + ";"
	}

	ctx.JSON(200, gin.H{"prompt": prompt})
}

func handleChatmodFalse(ctx *gin.Context) {
	var search typings.Search

	// Map request to Search struct
	// Get query
	search.Query = ctx.Query("query")
	// Get region
	search.Region = ctx.Query("region")
	// Get time range
	search.TimeRange = ctx.Query("time_range")
	if search.Query == "" {
		ctx.JSON(400, gin.H{"error": "Query is required"})
		return
	}
	// Get limit and check if it's a number
	limit := ctx.Query("limit")
	if limit != "" {
		if _, err := strconv.Atoi(limit); err != nil {
			ctx.JSON(400, gin.H{"error": "Limit must be a number"})
			return
		}
		search.Limit, _ = strconv.Atoi(limit)
	}
	// Get results
	results, err := duckduckgo.Get_results(search)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// Shorten results to limit if limit is set
	if search.Limit > 0 && search.Limit < len(results) {
		results = results[:search.Limit]
	}
	// Return results
	ctx.JSON(200, results)
}
