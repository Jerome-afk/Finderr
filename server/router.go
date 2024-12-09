package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var SessionTokenManager = make(map[string]string)

// GenerateToken generates a simple session token
func GenerateToken(username string) string {
	return username
}

// InitializeRoutes sets up all routes
func InitializeRoutes(router *gin.Engine) {
	router.LoadHTMLGlob("*.html")
	router.Static("/static", "./static")

	// Middleware
	router.Use(func(c *gin.Context) {
		_, err := os.Stat(".env") 
		if os.IsNotExist(err) {
		// Redirect all routes to setup
			if c.FullPath() != "/setup" {
				c.Redirect(302, "/setup")
				c.Abort()
				return
			}
		} else if err == nil {
			err := godotenv.Load(".env")
			if err != nil {
				fmt.Println("Error loading .env")
			} else {
				fmt.Println("Environment variable")
			}
		}
		c.Next()
	})

	router.GET("/setup", func(c *gin.Context) {
		c.HTML(200, "setup.html", nil)
	})

	router.POST("/setup", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		email := c.PostForm("email")
		movieDBAPI := c.PostForm("movie_db_api")
		aniListAPI := c.PostForm("anilist_api")
		audioDBAPI := c.PostForm("audiodb_api")
		dbType := c.PostForm("database_type")
		dbHost := c.PostForm("database_host")
		dbPort := c.PostForm("database_port")
		dbUser := c.PostForm("database_user")
		dbPass := c.PostForm("database_password")

		// Write to env
		envData := fmt.Sprintf(`USERNAME=%s
		PASSWORD=%s
		EMAIL=%s
		MOVIEDB_API=%s
		ANILIST_API=%s
		AUDIODB_API=%s
		DATABASE_TYPE=%s
		DATABASE_HOST=%s
		DATABASE_PORT=%s
		DATABASE_USER=%s
		DATABASE_PASS=%s`, username, password, email, movieDBAPI, aniListAPI,
	                      audioDBAPI, dbType, dbHost, dbPort, dbUser, dbPass)

		err := os.WriteFile(".env", []byte(envData), 0644)
		if err != nil {
			log.Printf("Error creating .env file: %v", err)
			c.JSON(500, gin.H{"message": "Error creating .env file"})
			return
		}
		
		c.JSON(200, gin.H{"message": "Environment file created successfully"})
		c.Redirect(http.StatusFound, "/login")
	})

	// Redirect root to login
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/login")
	})

	// Login page route
	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	router.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		remember := c.PostForm("remember") == "on"

		// Authenticate the user
		if Authenticate(username, password) {
			// Generate and store the session token
			token := GenerateToken(username)
			SessionTokenManager[token] = username

			// Set a secure cookie if "remember me" is checked
			if remember {
				c.SetCookie("session_token", token, 0, "/", "", false, true)
			}

			// Redirect to the home page
			c.Redirect(http.StatusFound, "/home")
		} else {
			logEvent("server", "Invalid login attempt for user: "+username)
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{
				"error": "Invalid username or password",
			})
		}
	})

	// Middleware for user session management
	router.Use(func(c *gin.Context) {
		token, err := c.Cookie("session_token")
		if err == nil {
			if username, ok := SessionTokenManager[token]; ok {
				c.Set("username", username)
			}
		}
		c.Next()
	})

	// Home page route
	router.GET("/home", authenticatedMiddleware, func(c *gin.Context) {
		username := c.MustGet("username").(string)
		c.HTML(http.StatusOK, "home.html", gin.H{"username": username})
	})

	// Media routes
	mediaRoutes := router.Group("/media", authenticatedMiddleware)
	{
		mediaRoutes.GET("/movies", func(c *gin.Context) {
			username := c.MustGet("username").(string)
			c.HTML(http.StatusOK, "movies.html", gin.H{"username": username})
		})

		mediaRoutes.GET("/shows", func(c *gin.Context) {
			username := c.MustGet("username").(string)
			c.HTML(http.StatusOK, "series.html", gin.H{"username": username})
		})

		mediaRoutes.GET("/anime", func(c *gin.Context) {
			username := c.MustGet("username").(string)
			c.HTML(http.StatusOK, "anime.html", gin.H{"username": username})
		})
	}
}

// Middleware to enforce authentication
func authenticatedMiddleware(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists || username == "" {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}
	c.Next()
}
