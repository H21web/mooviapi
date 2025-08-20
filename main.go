package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/Jisin0/filmigo/imdb"
	"github.com/Jisin0/filmigo/omdb"
	"github.com/Jisin0/filmigo/justwatch"
)

// MovieAPI struct holds all the movie service clients
type MovieAPI struct {
	imdbClient      *imdb.Client
	omdbClient      *omdb.Client
	justwatchClient *justwatch.Client
}

// NewMovieAPI creates a new instance of MovieAPI
func NewMovieAPI() *MovieAPI {
	omdbKey := os.Getenv("OMDB_API_KEY")
	if omdbKey == "" {
		log.Println("Warning: OMDB_API_KEY not set. OMDB features will be disabled.")
	}

	return &MovieAPI{
		imdbClient:      imdb.NewClient(),
		omdbClient:      omdb.NewClient(omdbKey),
		justwatchClient: justwatch.NewClient(),
	}
}

// Health check endpoint
func (api *MovieAPI) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "OK",
		"service":   "filmigo-api",
		"version":   "1.0.0",
		"message":   "Filmigo API is running successfully!",
	})
}

// Get movie by IMDB ID
func (api *MovieAPI) getMovieByID(c *gin.Context) {
	id := c.Param("id")
	
	// Validate IMDB ID format
	if !strings.HasPrefix(id, "tt") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid IMDB ID format. Must start with 'tt'",
			"example": "tt0111161",
		})
		return
	}

	movie, err := api.imdbClient.GetMovie(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Movie not found",
			"id":    id,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    movie,
	})
}

// Search movies using OMDB
func (api *MovieAPI) searchMovies(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Query parameter 'q' is required",
			"example": "/api/v1/movies/search?q=inception",
		})
		return
	}

	if api.omdbClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "OMDB service unavailable. API key not configured.",
		})
		return
	}

	results, err := api.omdbClient.Search(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Search failed",
			"query": query,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"query":   query,
		"results": results,
	})
}

// Get movie from OMDB by IMDB ID
func (api *MovieAPI) getOMDBMovie(c *gin.Context) {
	id := c.Param("id")
	
	if api.omdbClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "OMDB service unavailable. API key not configured.",
		})
		return
	}

	movie, err := api.omdbClient.GetMovie(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Movie not found in OMDB",
			"id":    id,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"source":  "omdb",
		"data":    movie,
	})
}

// CORS middleware
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}

// Request logging middleware
func requestLoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %s %s %d %s\n",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
		)
	})
}

func main() {
	// Initialize API
	api := NewMovieAPI()

	// Set Gin mode based on environment
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	r := gin.Default()

	// Middleware
	r.Use(corsMiddleware())
	r.Use(requestLoggingMiddleware())
	r.Use(gin.Recovery())

	// Root endpoint with API documentation
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "🎬 Welcome to Filmigo API",
			"version": "1.0.0",
			"description": "A powerful REST API for movie data",
			"endpoints": gin.H{
				"health":      "GET /api/v1/health",
				"movie_by_id": "GET /api/v1/movies/{imdb_id}",
				"search":      "GET /api/v1/movies/search?q={query}",
				"omdb_movie":  "GET /api/v1/omdb/{imdb_id}",
			},
			"examples": gin.H{
				"get_movie":     "/api/v1/movies/tt0111161",
				"search_movies": "/api/v1/movies/search?q=inception",
				"omdb_data":     "/api/v1/omdb/tt0111161",
				"health_check":  "/api/v1/health",
			},
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", api.healthCheck)
		
		// Movie endpoints
		v1.GET("/movies/:id", api.getMovieByID)
		v1.GET("/movies/search", api.searchMovies)
		
		// OMDB specific endpoint
		v1.GET("/omdb/:id", api.getOMDBMovie)
	}

	// Get port from environment (Render uses PORT env var)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default for local development
	}

	log.Printf("🎬 Filmigo API starting on port %s", port)
	log.Printf("🌐 API will be available at your Render URL")
	log.Printf("📋 Health check endpoint: /api/v1/health")
	
	// Start server - IMPORTANT: Bind to 0.0.0.0 for Render
	if err := r.Run("0.0.0.0:" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
