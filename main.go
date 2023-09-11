package main

import (
	"os"
	"strconv"
	"strings" 

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

	handler.GET("/search", func(ctx *gin.Context) {
		// Map request to Search struct
		var search typings.Search
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
		// Get text parameter
		search.Text = ctx.DefaultQuery("text", "false")

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
		// Check if we need to combine snippet fields
		if search.Text == "true" {
			var combinedSnippets []string
			for _, result := range results {
				combinedSnippets = append(combinedSnippets, result.Snippet)
			}
			combinedText := strings.Join(combinedSnippets, " ")
			ctx.String(200, combinedText)
		} else {
			ctx.JSON(200, results)
		}
	})

	endless.ListenAndServe(HOST+":"+PORT, handler)
}
