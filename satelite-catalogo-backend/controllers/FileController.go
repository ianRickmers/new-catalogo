package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// ServeArchivo sirve archivos del directorio de uploads de forma segura
func ServeArchivo(ctx *gin.Context) {
	fileParam := ctx.Param("filepath")
	cleanPath := filepath.Clean(fileParam)
	uploadRoot := os.Getenv("UPLOAD_DIR")
	fullPath := filepath.Join(uploadRoot, cleanPath)

	if !strings.HasPrefix(fullPath, filepath.Clean(uploadRoot)) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid path"})
		return
	}
	ctx.File(fullPath)
}
