package routes

import (
	"finderr/books"
	homePage "finderr/home"
	torrentDownload "finderr/torrentHandler"
	"finderr/torrentHandler/websocket"
	"finderr/upload"
	"os"

	"github.com/gin-gonic/gin"
)

type SpineItem struct {
	Href  string
	Title string
}

// InitializeRoutes sets up all routes
func InitializeRoutes(router *gin.Engine, hub *websocket.Hub) {
	router.LoadHTMLGlob("templates/*")

	router.Static("/static", "./static")
	router.Static("/uploads", "./uploads")

	// WebSocket endpoint
	router.GET("/ws", gin.WrapF(hub.HandleWebSocket))

	// Download endpoint
	router.GET("/download", func(c *gin.Context) {
		torrentDownload.DownloadHandler(c.Writer, c.Request, hub)
	})

	// Routes
	router.GET("/", homePage.HomeHandler)
	router.POST("/upload", upload.UploadHandler)
	router.GET("/read/:filename", books.ReadHandler)
	router.GET("/epub-content/:filename/*path", books.EpubContentHandler)

	// Create uploads directory if it doesn't exist
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", 0o755)
	}
}
