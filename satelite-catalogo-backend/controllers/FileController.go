package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// ServeArchivo godoc
// @Summary      Serve uploaded file
// @Description  Serves files from the uploads directory in a safe manner
// @Tags         files
// @Produce      octet-stream
// @Param        filepath  path  string  true  "File path"
// @Success      200  {file}  string
// @Failure      400  {object} map[string]interface{}
// @Router       /archivos/{filepath} [get]
func ServeArchivo(ctx *gin.Context) {
	fileParam := ctx.Param("filepath")
	cleanPath := filepath.Clean(fileParam)
	uploadRoot := os.Getenv("UPLOAD_DIR")
	if uploadRoot == "" {
		uploadRoot = "./uploads"
	}
	fullPath := filepath.Join(uploadRoot, cleanPath)

	if !strings.HasPrefix(fullPath, filepath.Clean(uploadRoot)) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid path"})
		return
	}
	ctx.File(fullPath)
}
