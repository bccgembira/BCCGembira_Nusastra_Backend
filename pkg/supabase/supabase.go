package supabase

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/gofiber/fiber/v2"
)

type ISupabase interface {
	Upload(file *multipart.FileHeader) (string, error)
	Delete(link string) error
	ConvertFile(link string) string
}

type supabase struct {
	url       string
	publicUrl string
	token     string
	client    http.Client
}

func NewSupabase() ISupabase {
	url := fmt.Sprintf("%s/storage/v1/object/%s/", os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_BUCKET"))
	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/", os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_BUCKET"))

	return &supabase{
		url:       url,
		publicUrl: publicURL,
		token:     os.Getenv("SUPABASE_TOKEN"),
		client:    http.Client{},
	}
}

func (us *supabase) ConvertFile(link string) string {
	return strings.ReplaceAll(link, us.publicUrl, "")
}

func (us *supabase) Upload(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return "", err
	}

	url := us.url + file.Filename

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(fileBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+ us.token)
	req.Header.Set("Content-Type", file.Header.Get("Content-Type"))

	resp, err := us.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", err
	}

	publicURL := us.publicUrl + file.Filename
	return publicURL, nil
}

func (us *supabase) Delete(link string) error {
	fileName := us.ConvertFile(link)
	url := us.url + "/" + fileName
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "Bearer "+us.token)
	response, err := us.client.Do(request)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("error read file : %v", err)
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		log.Printf("%v %v \n", color.RedString("Received non-200 response:"), response.StatusCode)
		return fiber.ErrBadRequest
	}
	return nil
}
