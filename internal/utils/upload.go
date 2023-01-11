package utils

import (
	"errors"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/go-template-api/internal/domain"
)

func UploadResizeMultipleFile(c *gin.Context, info *domain.ImageInput) ([]string, error) {
	filePaths := []string{}
	// c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, int64(30<<20))
	form, err := c.MultipartForm()
	if err != nil {
		return filePaths, err
	}

	pathDir := fmt.Sprintf("public/%s/%s", info.UserID, info.Service)
	if _, err := os.Stat(pathDir); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(pathDir, os.ModePerm)
		if err != nil {
			// log.Println(err)
			return filePaths, err
		}
	}

	files := form.File["images[]"]
	for _, file := range files {
		fileExt := filepath.Ext(file.Filename)
		fmt.Println("fileExt", file.Filename)
		originalFileName := strings.TrimSuffix(filepath.Base(file.Filename), filepath.Ext(file.Filename))
		now := time.Now()
		filenameOriginal := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix()) + fileExt
		filePaths = append(filePaths, filenameOriginal)
		filenameLg := "lg-" + strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix()) + fileExt
		// filePaths["lg"] = filenameLg
		filenameXs := "xs-" + strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix()) + fileExt
		// filePaths["xs"] = filenameXs

		readerFile, _ := file.Open()
		imageFile, _, err := image.Decode(readerFile)
		if err != nil {
			// log.Fatal(err)
			return filePaths, err
		}

		err = imaging.Save(imageFile, fmt.Sprintf("%s/%v", pathDir, filenameOriginal), imaging.JPEGQuality(80))
		if err != nil {
			// log.Fatalf("failed to save image: %v", err)
			return filePaths, err
		}
		src := imaging.Resize(imageFile, 1000, 0, imaging.Lanczos)
		err = imaging.Save(src, fmt.Sprintf("%s/%v", pathDir, filenameLg), imaging.JPEGQuality(80))
		if err != nil {
			// log.Fatalf("failed to save image: %v", err)
			return filePaths, err
		}
		src_xs := imaging.Resize(imageFile, 30, 0, imaging.Lanczos)
		err = imaging.Save(src_xs, fmt.Sprintf("%s/%v", pathDir, filenameXs), imaging.JPEGQuality(30))
		if err != nil {
			// log.Fatalf("failed to save image: %v", err)
			return filePaths, err
		}
	}

	// ctx.JSON(http.StatusOK, gin.H{"filepaths": filePaths})
	return filePaths, nil
}
