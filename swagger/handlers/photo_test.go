package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/generated/restapi/operations/photos"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-openapi/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/ent"
	repomock "git.epam.com/epm-lstr/epm-lstr-lc/be/internal/mocks/repositories"
	servicesmock "git.epam.com/epm-lstr/epm-lstr-lc/be/internal/mocks/services"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/generated/models"
)

type PhotoTestSuite struct {
	suite.Suite
	logger       *zap.Logger
	repository   *repomock.PhotoRepository
	fileManager  *servicesmock.FileManager
	handler      *Photo
	serverURL    string
	photoURLPath string
}

func TestPhotoSuite(t *testing.T) {
	suite.Run(t, new(PhotoTestSuite))
}

func (s *PhotoTestSuite) SetupTest() {
	s.logger = zap.NewNop()
	s.serverURL = "http://localhost:8080/"
	s.photoURLPath = "api/equipment/photos/"
	s.repository = &repomock.PhotoRepository{}
	s.fileManager = &servicesmock.FileManager{}
	s.handler = NewPhoto(s.serverURL, s.logger)
}

func (s *PhotoTestSuite) TestPhoto_CreatePhoto_EmptyFile() {
	t := s.T()
	request := http.Request{}

	fileName := "testimagename.jpg"
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(fileName)
	defer f.Close()

	data := photos.CreateNewPhotoParams{
		HTTPRequest: &request,
		File:        f,
	}

	handlerFunc := s.handler.CreateNewPhotoFunc(s.repository, s.fileManager)
	resp := handlerFunc.Handle(data)

	responseRecorder := httptest.NewRecorder()
	producer := runtime.JSONProducer()
	resp.WriteResponse(responseRecorder, producer)
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)

	response := models.Error{}
	err = json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, response)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, "File is empty", response.Data.Message)

	s.repository.AssertExpectations(t)
}

func (s *PhotoTestSuite) TestPhoto_CreatePhoto_WrongMimeType() {
	t := s.T()
	request := http.Request{}
	fileName := "testfile.txt"

	err := createNonEmptyFile(fileName, []byte("This is txt"))
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(fileName)

	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	data := photos.CreateNewPhotoParams{
		HTTPRequest: &request,
		File:        f,
	}

	handlerFunc := s.handler.CreateNewPhotoFunc(s.repository, s.fileManager)
	resp := handlerFunc.Handle(data)

	responseRecorder := httptest.NewRecorder()
	producer := runtime.JSONProducer()
	resp.WriteResponse(responseRecorder, producer)
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)

	response := models.Error{}
	err = json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, response)
	assert.NotEmpty(t, response.Data)
	assert.Containsf(t, response.Data.Message, "Wrong file format", "returned wrong error")

	s.repository.AssertExpectations(t)
}

func (s *PhotoTestSuite) TestPhoto_CreatePhoto_OK() {
	t := s.T()
	request := http.Request{}
	ctx := request.Context()

	id := "testimagename"
	url := "http://localhost:8080/api/equipments/photos/testimagename"
	fileName := "testimagename.jpg"

	img, err := generateImageBytes()
	if err != nil {
		log.Fatal(err)
	}
	err = createNonEmptyFile(fileName, img)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(fileName)

	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	data := photos.CreateNewPhotoParams{
		HTTPRequest: &request,
		File:        f,
	}
	s.repository.On("CreatePhoto", ctx, models.Photo{
		ID:       id,
		URL:      &url,
		FileName: fileName,
	}).Return(&ent.Photo{
		ID:       id,
		URL:      url,
		FileName: fileName,
	}, nil)
	s.fileManager.On("GenerateFileName").Return(id, nil)
	s.fileManager.On("SaveDataToFile", img, fileName).Return(nil)
	s.fileManager.On("BuildFileURL", s.serverURL, s.photoURLPath, id).Return(url, nil)

	handlerFunc := s.handler.CreateNewPhotoFunc(s.repository, s.fileManager)
	resp := handlerFunc.Handle(data)

	responseRecorder := httptest.NewRecorder()
	producer := runtime.JSONProducer()
	resp.WriteResponse(responseRecorder, producer)
	assert.Equal(t, http.StatusCreated, responseRecorder.Code)

	returnedPhoto := models.CreateNewPhotoResponse{}
	err = json.Unmarshal(responseRecorder.Body.Bytes(), &returnedPhoto)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, id, returnedPhoto.Data.ID)
	assert.Equal(t, url, *returnedPhoto.Data.URL)
	assert.Equal(t, fileName, returnedPhoto.Data.FileName)

	s.repository.AssertExpectations(t)
}

func (s *PhotoTestSuite) TestPhoto_GetPhoto_OK() {
	t := s.T()
	request := http.Request{}
	ctx := request.Context()

	id := "testimagename"
	url := "http://localhost:8080/api/equipments/photos/testimagename"
	fileName := "testimagename.jpg"

	data := photos.GetPhotoParams{
		HTTPRequest: &request,
		PhotoID:     id,
	}

	s.repository.On("PhotoByID", ctx, data.PhotoID).Return(&ent.Photo{
		ID:       id,
		URL:      url,
		FileName: fileName,
	}, nil)
	s.fileManager.On("ReadFile", fileName).Return([]byte{1, 1, 1}, nil)

	handlerFunc := s.handler.GetPhotoFunc(s.repository, s.fileManager)
	resp := handlerFunc.Handle(data)

	responseRecorder := httptest.NewRecorder()
	producer := runtime.JSONProducer()
	resp.WriteResponse(responseRecorder, producer)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Equal(t, "image/jpg", responseRecorder.Header().Get("Content-Type"))

	assert.Equal(t, []byte{1, 1, 1}, responseRecorder.Body.Bytes())

	s.repository.AssertExpectations(t)
}

func (s *PhotoTestSuite) TestPhoto_GetPhoto_RepoErr() {
	t := s.T()
	request := http.Request{}
	ctx := request.Context()

	id := "testimagename"
	data := photos.GetPhotoParams{
		HTTPRequest: &request,
		PhotoID:     id,
	}

	errorToReturn := errors.New("repo err")
	s.repository.On("PhotoByID", ctx, data.PhotoID).Return(nil, errorToReturn)

	handlerFunc := s.handler.GetPhotoFunc(s.repository, s.fileManager)
	resp := handlerFunc(data)

	responseRecorder := httptest.NewRecorder()
	producer := runtime.JSONProducer()
	resp.WriteResponse(responseRecorder, producer)
	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
	assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"))

	response := models.Error{}
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, response)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, errorToReturn.Error(), response.Data.Message)

	s.repository.AssertExpectations(t)
}

func (s *PhotoTestSuite) TestPhoto_DownloadPhoto_OK() {
	t := s.T()
	request := http.Request{}
	ctx := request.Context()

	id := "testimagename"
	url := "http://localhost:8080/api/equipments/photos/testimagename"
	fileName := "testimagename.jpg"

	data := photos.DownloadPhotoParams{
		HTTPRequest: &request,
		PhotoID:     id,
	}

	bytesToReturn := []byte{1, 1, 1, 1}
	s.repository.On("PhotoByID", ctx, data.PhotoID).Return(&ent.Photo{
		ID:       id,
		URL:      url,
		FileName: fileName,
	}, nil)
	s.fileManager.On("ReadFile", fileName).Return(bytesToReturn, nil)

	handlerFunc := s.handler.DownloadPhotoFunc(s.repository, s.fileManager)
	resp := handlerFunc.Handle(data)

	responseRecorder := httptest.NewRecorder()
	producer := runtime.JSONProducer()
	resp.WriteResponse(responseRecorder, producer)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Equal(t, "application/octet-stream", responseRecorder.Header().Get("Content-Type"))

	assert.Equal(t, bytesToReturn, responseRecorder.Body.Bytes())

	s.repository.AssertExpectations(t)
}

func (s *PhotoTestSuite) TestPhoto_DeletePhoto_NotExists() {
	t := s.T()
	request := http.Request{}
	ctx := request.Context()

	id := "testimagename"
	data := photos.DeletePhotoParams{
		HTTPRequest: &request,
		PhotoID:     id,
	}

	errorToReturn := errors.New("not found")
	s.repository.On("PhotoByID", ctx, data.PhotoID).Return(nil, errorToReturn)

	handlerFunc := s.handler.DeletePhotoFunc(s.repository, s.fileManager)
	resp := handlerFunc.Handle(data)

	responseRecorder := httptest.NewRecorder()
	producer := runtime.JSONProducer()
	resp.WriteResponse(responseRecorder, producer)
	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)

	response := models.Error{}
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, response)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, errorToReturn.Error(), response.Data.Message)

	s.repository.AssertExpectations(t)
}

func (s *PhotoTestSuite) TestPhoto_DeletePhoto_OK() {
	t := s.T()
	request := http.Request{}
	ctx := request.Context()

	id := "testimagename"
	url := "http://localhost:8080/api/equipments/photos/testimagename"
	fileName := "testimagename.jpg"

	data := photos.DeletePhotoParams{
		HTTPRequest: &request,
		PhotoID:     id,
	}
	s.repository.On("PhotoByID", ctx, data.PhotoID).Return(&ent.Photo{
		ID:       id,
		URL:      url,
		FileName: fileName,
	}, nil)
	s.repository.On("DeletePhotoByID", ctx, data.PhotoID).Return(nil)
	s.fileManager.On("DeleteFile", fileName).Return(nil)

	handlerFunc := s.handler.DeletePhotoFunc(s.repository, s.fileManager)
	resp := handlerFunc.Handle(data)

	responseRecorder := httptest.NewRecorder()
	producer := runtime.JSONProducer()
	resp.WriteResponse(responseRecorder, producer)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	s.repository.AssertExpectations(t)
}

func createNonEmptyFile(name string, content []byte) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	_, err = f.Write(content)
	if err != nil {
		if err := os.Remove(name); err != nil {
			return err
		}
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

func generateImageBytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, image.Rect(0, 0, 100, 100), nil)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
