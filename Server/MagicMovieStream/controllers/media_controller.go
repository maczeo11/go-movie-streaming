package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// MediaDir returns the directory that holds the video files for streaming.
func MediaDir() string {
	if dir := os.Getenv("MEDIA_DIR"); dir != "" {
		return dir
	}
	return "./media"
}

// StreamMedia serves video files with HTTP Range support so browsers can
// seek and stream instead of downloading the whole file. Path traversal is
// blocked by resolving the path and checking it stays inside MediaDir().
func StreamMedia() gin.HandlerFunc {
	return func(c *gin.Context) {
		mediaDir := MediaDir()

		name := filepath.Base(c.Param("filepath"))
		if name == "" || name == "." || name == string(os.PathSeparator) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing filename"})
			return
		}

		absMedia, err := filepath.Abs(mediaDir)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "media dir not configured"})
			return
		}
		absPath, err := filepath.Abs(filepath.Join(mediaDir, name))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filename"})
			return
		}
		if !strings.HasPrefix(absPath, absMedia+string(os.PathSeparator)) {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid path"})
			return
		}

		if _, err := os.Stat(absPath); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
			return
		}

		// gin's File() handles Content-Type sniffing and Range headers,
		// which is what lets the <video> tag seek around.
		c.File(absPath)
	}
}

// ListMedia returns the video files currently available to stream.
func ListMedia() gin.HandlerFunc {
	return func(c *gin.Context) {
		entries, err := os.ReadDir(MediaDir())
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"data": []string{}, "count": 0})
			return
		}

		var files []string
		for _, entry := range entries {
			if !entry.IsDir() && isVideoExt(filepath.Ext(entry.Name())) {
				files = append(files, entry.Name())
			}
		}
		c.JSON(http.StatusOK, gin.H{"data": files, "count": len(files)})
	}
}

func isVideoExt(ext string) bool {
	switch strings.ToLower(ext) {
	case ".mp4", ".webm", ".ogg", ".mov", ".mkv":
		return true
	}
	return false
}
