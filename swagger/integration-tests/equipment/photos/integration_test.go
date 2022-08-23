package photos

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/client/photos"
	utils "git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/integration-tests"
)

var (
	testLogin    string
	testPassword string
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	beClient := utils.SetupClient()

	var err error
	testLogin, testPassword, err = utils.GenerateLoginAndPassword()
	if err != nil {
		log.Fatalf("GenerateLoginAndPassword: %v", err)
	}
	_, err = utils.CreateUser(ctx, beClient, testLogin, testPassword)
	if err != nil {
		log.Fatalf("CreateUser: %v", err)
	}
	os.Exit(m.Run())
}

func TestIntegration_PhotosUpload(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	beClient := utils.SetupClient()

	loginUser, err := utils.LoginUser(ctx, beClient, testLogin, testPassword)
	require.NoError(t, err)

	t.Run("", func(t *testing.T) {
		token := loginUser.GetPayload().AccessToken
		fileName := "testfile.txt"

		//id := "testimagename"
		//url := "http://localhost:8080/api/equipments/photos/testimagename"
		//fileName := "testimagename.jpg"

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

		params := photos.NewCreateNewPhotoParamsWithContext(ctx)
		params.File = f
		res, err := beClient.Photos.CreateNewPhoto(params, utils.AuthInfoFunc(token))
		require.NoError(t, err)

		assert.Equal(t, fileName, res.Payload.Data.FileName)
	})
}

func generateImageBytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, image.Rect(0, 0, 100, 100), nil)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
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
