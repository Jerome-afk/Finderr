package books

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func EpubContentHandler(c *gin.Context) {
	filename := c.Param("filename")
	path := c.Param("path")

	// Open the EPUB file as a ZIP archive
	r, err := zip.OpenReader(fmt.Sprintf("uploads/%s", filename))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to open EPUB: " + err.Error()})
		return
	}
	defer r.Close()

	// Find the requested file
	for _, f := range r.File {
		if f.Name == path[1:] { // Remove leading slash from path
			rc, err := f.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file in EPUB: " + err.Error()})
				return
			}
			defer rc.Close()

			content, err := io.ReadAll(rc)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file: " + err.Error()})
				return
			}

			// Set appropriate content type
			contentType := "text/html"
			switch {
			case strings.HasSuffix(path, ".css"):
				contentType = "text/css"
			case strings.HasSuffix(path, ".js"):
				contentType = "application/javascript"
			case strings.HasSuffix(path, ".png"):
				contentType = "image/png"
			case strings.HasSuffix(path, ".jpg"), strings.HasSuffix(path, ".jpeg"):
				contentType = "image/jpeg"
			case strings.HasSuffix(path, ".gif"):
				contentType = "image/gif"
			case strings.HasSuffix(path, ".svg"):
				contentType = "image/svg+xml"
			}

			c.Data(http.StatusOK, contentType, content)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "File not found in EPUB"})
}
