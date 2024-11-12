package claudinary

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"strconv"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Service interface {
	UploadImage(file *multipart.FileHeader, imageType string, publicID string) (string, error)
	DeleteImage(publicID string) (bool, error)
}

type service struct {
	cld *cloudinary.Cloudinary
}

func NewService() (*service, error) {
	var CLOUDINARY_URL = os.Getenv("CLOUDINARY_URL")
	cld, err := cloudinary.NewFromURL(CLOUDINARY_URL)
	if err != nil {
		return nil, fmt.Errorf("cloudinary URL not set in environment variables")
	}	

	return &service{
		cld: cld,
	}, nil
}

func (s *service) UploadImage(file *multipart.FileHeader, imageType string, publicID string) (string, error) {

	ctx := context.Background()

	folderImg := ""

	switch imageType {
	case "campaign":
			folderImg = "bwastartup/campaign-images"
	case "user":
			folderImg = "bwastartup/user-images"
	}


	resp, err := s.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: folderImg,
		PublicID: strconv.FormatInt(time.Now().Unix(), 10) + "-" + publicID,
	})

	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}

func (s *service) DeleteImage(publicID string) (bool, error) {
	ctx := context.Background()

	_, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})

	if err != nil {
		return false, err
	}


	return true, nil
}