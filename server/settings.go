package server

import (
	"encoding/json"
	//"net/http"
	"fmt"
	"os"

	//"github.com/gin-gonic/gin"
)

// Check user authorization
func isAuthorized(username string) bool {
	filePath := fmt.Sprintf("users/%s.json", username)
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	var userData struct {
		UserLevel int `json:"user_level"`
	}

	if err := json.NewDecoder(file).Decode(&userData); err != nil {
		return false
	}

	return userData.UserLevel == 1
}

// Setings route
// func initSettingsRoutes(router *gin.Context) {
// 	router.GET("/settings", func(c *gin.Context) {
// 		username, exists := c.Get("username")
// 		if !exists || !isAuthorized(username.(string)) {
// 			c.HTML(http.StatusUnauthorized, "home.html", nil)
// 			return
// 		}
// 		c.HTML(http.StatusOK, "settings.html", nil)
// 	})

// 	router.POST("/settings/downloads", func(c *gin.Context))
// }