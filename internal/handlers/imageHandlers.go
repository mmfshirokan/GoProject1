package handlers

import (
	"fmt"
	"io"
	"os"

	"github.com/caarlos0/env/v10"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/mmfshirokan/GoProject1/internal/model"
	log "github.com/sirupsen/logrus"
)

// UploadImage godoc
//
// @Summary upload image
// @Description Uploads user image to the server
// @Tags User ImageHandlers
// @Accept json
// @Param img body model.Image true "Path to image and image name"
// @Success 200
// @Router /users/auth/uploadImage [put]
func (handling *Handler) UploadImage(c echo.Context) error {
	img, err := newImage(c)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	src, err := os.Open(img.LocalImgPath)
	if err != nil {
		log.Error(fmt.Errorf("error occurred while opening file: %w", err))

		return fmt.Errorf("error occurred while opening file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(fmt.Sprint(img.ServerImgPath, img.Name, ".png"))
	if err != nil {
		log.Error(fmt.Errorf("error occurred while creating file: %w", err))

		return fmt.Errorf("error occurred while creating file: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		log.Error(fmt.Errorf("error occurred while copying file: %w", err))

		return fmt.Errorf("error occurred while copying file: %w", err)
	}

	return nil
}

// DownloadImage godoc
//
// @Summary download image
// @Description Downloads user image
// @Tags User ImageHandlers
// @Accept json
// @Param img body model.Image true "Path where to download image and image name"
// @Success 1
// @Router /users/auth/downloadImage [put]
func (handling *Handler) DownloadImage(c echo.Context) error {
	img, err := newImage(c)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	src, err := os.Open(fmt.Sprint(img.ServerImgPath, img.Name, ".png"))
	if err != nil {
		log.Error(fmt.Errorf("error occurred while opening file: %w", err))

		return fmt.Errorf("error occurred while opening file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(img.LocalImgPath)
	if err != nil {
		log.Error(fmt.Errorf("error occurred while creating file: %w", err))

		return fmt.Errorf("error occurred while creating file: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		log.Error(fmt.Errorf("error occurred while copying file: %w", err))

		return fmt.Errorf("error occurred while copying file: %w", err)
	}

	return nil
}

func newImage(c echo.Context) (model.Image, error) {
	logInit()

	var img model.Image
	if err := c.Bind(&img); err != nil {
		log.Error(fmt.Errorf("bind error in handlers.UploadImage: %w", err))

		return model.Image{}, fmt.Errorf("method: handlers.UploadImage; bind error: %w", err)
	}

	if err := env.Parse(&img); err != nil {
		log.Error(fmt.Errorf("env.Parse: %w", err))

		return model.Image{}, fmt.Errorf("env.Parse: %w", err)
	}

	val := validator.New(validator.WithRequiredStructEnabled())
	if err := val.Struct(&img); err != nil {
		log.Error(fmt.Errorf("invalid path at model.Image: %w", err))

		return model.Image{}, fmt.Errorf("model.Image err: %w", err)
	}

	return img, nil
}
