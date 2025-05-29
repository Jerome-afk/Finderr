package upload

import (
	"finderr/logger"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func UploadHandler(c *gin.Context) {
	file, err := c.FormFile("epub")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.LogEvent("User", err.Error())
		return
	}

	// Ensure the file has an .epub extension
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".epub") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File must be an EPUB"})
		logger.LogEvent("User", "File must be an EPUB")
		return
	}

	// Save the file
	filename := filepath.Base(file.Filename)
	dst := fmt.Sprintf("uploads/%s", filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		logger.LogEvent("Server", err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/read/%s", filename))
}
