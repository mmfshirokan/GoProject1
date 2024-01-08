package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	jsonUploadImg   = `{"path":"/home/andreishyrakanau/temp/spudy.png","imgname":"servSpudy"}`
	jsonDownloadImg = `{"path":"/home/andreishyrakanau/temp/userSpudy.png","imgname":"test"}`
)

func TestUploadImage(t *testing.T) { // setup required: nedd spudy.png in temp; optional: chek if spudy.png is in temp if not create
	handler := NewHandler(nil, nil, nil)

	method := http.MethodPut
	target := "/users/auth/uploadImage"

	e := echo.New()
	rec := httptest.NewRecorder()

	testTable := []struct {
		name string
		body string

		destinationFilePath string
		hasError            bool
	}{
		{
			name: "standart input without error",
			body: jsonUploadImg,

			destinationFilePath: "/home/andreishyrakanau/projects/project1/GoProject1/images/servSpudy.png",
			hasError:            false,
		},
		{
			name: "input with wrong json",
			body: `{"wrong":"json"}`,

			hasError: true,
		},
	}

	for _, test := range testTable {
		req := httptest.NewRequest(method, target, strings.NewReader(test.body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)

		err := handler.UploadImage(c)

		if test.hasError {
			assert.Error(t, err, test.name)
			return
		}

		assert.Nil(t, err, test.name)
		assert.Equal(t, 200, rec.Code)

		_, err = os.Stat(test.destinationFilePath)
		if destenationFileNotExist := errors.Is(err, os.ErrNotExist); destenationFileNotExist {
			t.Error("destination file do not exist")
			return
		}

		if err = os.Remove(test.destinationFilePath); err != nil {
			t.Error("failed to remove destenation file")
		}
	}
}

func TestDownloadImage(t *testing.T) {
	handler := NewHandler(nil, nil, nil)

	method := http.MethodPut
	target := "/users/auth/downloadImage"

	e := echo.New()
	rec := httptest.NewRecorder()

	testTable := []struct {
		name string
		body string

		destinationFilePath string
		hasError            bool
	}{
		{
			name: "standart input witout error",
			body: jsonDownloadImg,

			destinationFilePath: "/home/andreishyrakanau/temp/userSpudy.png",
			hasError:            false,
		},
		{
			name: "input with wrong json",
			body: `{"wrong":"json"}`,

			hasError: true,
		},
	}

	for _, test := range testTable {
		req := httptest.NewRequest(method, target, strings.NewReader(test.body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)

		err := handler.DownloadImage(c)
		if test.hasError {
			assert.Error(t, err, test.name)
			return
		}

		assert.Nil(t, err, test.name)
		assert.Equal(t, 200, rec.Code)

		_, err = os.Stat(test.destinationFilePath)
		if destenationFileNotExist := errors.Is(err, os.ErrNotExist); destenationFileNotExist {
			t.Error("destination file do not exist")
			return
		}

		if err = os.Remove(test.destinationFilePath); err != nil {
			t.Error("failed to remove destenation file")
		}
	}
}
