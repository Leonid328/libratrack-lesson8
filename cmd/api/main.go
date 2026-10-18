package main

import (
	"log"
	"net/http"
	"os"

	"libratrack/internal/handler"
	"libratrack/internal/middleware"
	"libratrack/internal/repository"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return envError("DATABASE_URL is required")
	}
	migrPath := os.Getenv("MIGRATIONS_PATH")
	if migrPath == "" {
		migrPath = "migrations"
	}
	if err := repository.RunMigrations(dbURL, migrPath); err != nil {
		log.Printf("migrations: %v", err)
	}
	db, err := repository.Open(dbURL)
	if err != nil {
		return err
	}

	bh := handler.NewBookHandler(repository.NewBookRepository(db))
	ah := handler.NewAuthorHandler(repository.NewAuthorRepository(db))
	ch := handler.NewCategoryHandler(repository.NewCategoryRepository(db))
	cph := handler.NewCopyHandler(repository.NewCopyRepository(db))

	r := gin.New()
	r.Use(gin.Recovery(), middleware.Logging(), cors())
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/swagger/index.html", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, swaggerHTML)
	})

	v1 := r.Group("/api/v1")
	{
		v1.GET("/books", bh.List)
		v1.GET("/books/:id", bh.Get)
		v1.POST("/books", bh.Create)
		v1.PUT("/books/:id", bh.Update)
		v1.DELETE("/books/:id", bh.Delete)
		v1.POST("/books/:id/copies", cph.Create)
		v1.GET("/books/:id/copies", cph.List)
		v1.GET("/books/:id/availability", cph.Availability)
		v1.POST("/authors", ah.Create)
		v1.GET("/authors", ah.List)
		v1.GET("/authors/:id", ah.Get)
		v1.POST("/categories", ch.Create)
		v1.GET("/categories", ch.List)
		v1.GET("/categories/:id", ch.Get)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return r.Run(":" + port)
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

type envError string

func (e envError) Error() string { return string(e) }

const swaggerHTML = `<!DOCTYPE html><html><head><title>LibraTrack API</title></head>
<body><h1>LibraTrack OpenAPI</h1><p>CRUD books, authors, categories, copies.</p></body></html>
`
