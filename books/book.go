package books

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type SpineItem struct {
	Href  string
	Title string
}

func ReadHandler(c *gin.Context) {
	filename := c.Param("filename")
	filePath := fmt.Sprintf("uploads/%s", filename)

	// Open the EPUB file as a ZIP archive
	r, err := zip.OpenReader(filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to open EPUB: " + err.Error()})
		return
	}
	defer r.Close()

	// Try to find the container.xml to locate the OPF file
	var opfPath string
	for _, f := range r.File {
		if f.Name == "META-INF/container.xml" {
			rc, err := f.Open()
			if err != nil {
				continue
			}
			defer rc.Close()

			buf := new(bytes.Buffer)
			_, err = io.Copy(buf, rc)
			if err != nil {
				continue
			}

			// Simple XML parsing to find the OPF path
			content := buf.String()
			if strings.Contains(content, "rootfile") {
				start := strings.Index(content, "full-path=\"") + len("full-path=\"")
				end := strings.Index(content[start:], "\"")
				if start > len("full-path=\"") && end > 0 {
					opfPath = content[start : start+end]
				}
			}
			break
		}
	}

	if opfPath == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find OPF file in EPUB"})
		return
	}

	// Parse the OPF file to get the spine (reading order)
	var readingOrder []SpineItem
	for _, f := range r.File {
		if f.Name == opfPath {
			rc, err := f.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open OPF file: " + err.Error()})
				return
			}
			defer rc.Close()

			buf := new(bytes.Buffer)
			_, err = io.Copy(buf, rc)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read OPF file: " + err.Error()})
				return
			}

			// Simple XML parsing to find the spine items
			content := buf.String()
			readingOrder = parseSpineFromOPF(content)
			break
		}
	}

	// Get the title (simplified)
	title := "Untitled"
	for _, f := range r.File {
		if f.Name == opfPath {
			rc, err := f.Open()
			if err != nil {
				continue
			}
			defer rc.Close()

			buf := new(bytes.Buffer)
			_, err = io.Copy(buf, rc)
			if err != nil {
				continue
			}

			// Simple title extraction
			content := buf.String()
			if strings.Contains(content, "<dc:title>") {
				start := strings.Index(content, "<dc:title>") + len("<dc:title>")
				end := strings.Index(content[start:], "</dc:title>")
				if start > len("<dc:title>") && end > 0 {
					title = content[start : start+end]
				}
			}
			break
		}
	}

	c.HTML(http.StatusOK, "reader.html", gin.H{
		"Title":         title,
		"Filename":      filename,
		"Reading Order": readingOrder,
	})
}

func parseSpineFromOPF(content string) []SpineItem {
	var spineItems []SpineItem

	// This is a simplified parser - a real implementation would need proper XML parsing
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, "<item ") {
			// Parse manifest items
			if strings.Contains(line, "id=\"") && strings.Contains(line, "href=\"") {
				idStart := strings.Index(line, "id=\"") + len("id=\"")
				idEnd := strings.Index(line[idStart:], "\"")
				id := line[idStart : idStart+idEnd]

				hrefStart := strings.Index(line, "href=\"") + len("href=\"")
				hrefEnd := strings.Index(line[hrefStart:], "\"")
				href := line[hrefStart : hrefStart+hrefEnd]

				spineItems = append(spineItems, SpineItem{
					Href:  href,
					Title: id,
				})
			}
		}
	}

	return spineItems
}