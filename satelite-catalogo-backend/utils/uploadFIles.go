package utils

import (
	"fmt"
	"io"
	"mime/multipart"

	"os"
	"path/filepath"
	"strings"
)

// funcion que se encarga de guardar los archivos subidos en una carpeta según su prefijo
// de que linea pertenecen, a si misma todos  los archivos que tengan el prefijo "linea_X" se guardan en una subcarpeta
// segun el id de la solicitud
func GuardarArchivos(archivos []*multipart.FileHeader, carpetaID string) ([]string, error) {
	uploadRoot := os.Getenv("UPLOAD_DIR")
	destinoRaiz := filepath.Join(uploadRoot, carpetaID)

	if err := os.MkdirAll(destinoRaiz, os.ModePerm); err != nil {
		return nil, err
	}

	var rutasRelativas []string

	for _, fileHeader := range archivos {
		src, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer src.Close()

		nombreOriginal := filepath.Base(fileHeader.Filename)

		// Extraer prefijo tipo "linea_1" y el resto del nombre
		var subCarpeta string
		var nombreFinal string

		// Detectar si el nombre sigue el patrón "linea_X_..."
		if strings.HasPrefix(nombreOriginal, "linea_") {
			partes := strings.SplitN(nombreOriginal, "_", 3)
			if len(partes) >= 3 {
				subCarpeta = fmt.Sprintf("linea_%s", partes[1])
				nombreFinal = strings.TrimPrefix(partes[2], "_")
			} else {
				subCarpeta = "sin_linea"
				nombreFinal = nombreOriginal
			}
		} else {
			subCarpeta = "sin_linea"
			nombreFinal = nombreOriginal
		}

		subCarpeta = filepath.Base(subCarpeta)
		nombreFinal = filepath.Base(nombreFinal)
		if strings.Contains(nombreFinal, "..") {
			return nil, fmt.Errorf("invalid file name")
		}

		// Crear subcarpeta si no existe
		destinoSub := filepath.Join(destinoRaiz, subCarpeta)
		if err := os.MkdirAll(destinoSub, os.ModePerm); err != nil {
			return nil, err
		}

		pathDestino := filepath.Join(destinoSub, nombreFinal)
		if !strings.HasPrefix(filepath.Clean(pathDestino), filepath.Clean(destinoSub)) {
			return nil, fmt.Errorf("invalid file path")
		}

		dst, err := os.Create(pathDestino)
		if err != nil {
			return nil, err
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return nil, err
		}

		// Ruta relativa para base de datos
		rutaRelativa := filepath.ToSlash(filepath.Join("/archivos", carpetaID, subCarpeta, nombreFinal))
		rutasRelativas = append(rutasRelativas, rutaRelativa)
	}

	return rutasRelativas, nil
}
