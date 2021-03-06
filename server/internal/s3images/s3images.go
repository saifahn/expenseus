package s3images

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/saifahn/expenseus/internal/app"
)

type ImageStoreS3 struct {
	session *session.Session
	bucket  string
}

func New(useConfig bool) *ImageStoreS3 {
	var sess *session.Session
	if useConfig {
		sess = session.Must(session.NewSession(&aws.Config{
			Credentials:      credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
			Endpoint:         aws.String(os.Getenv("S3_ENDPOINT_LOCAL")),
			Region:           aws.String("ap-northeast-1"),
			DisableSSL:       aws.Bool(true),
			S3ForcePathStyle: aws.Bool(true),
		}))
	} else {
		sess = session.Must(session.NewSession())
	}

	bucket := os.Getenv("AWS_IMAGES_BUCKET_KEY")

	return &ImageStoreS3{
		session: sess,
		bucket:  bucket,
	}
}

func (i *ImageStoreS3) Upload(file multipart.File, header multipart.FileHeader) (string, error) {
	uploader := s3manager.NewUploader(i.session)
	// uuid generated to keep key unique but filename is preserved
	key := fmt.Sprintf("%v/%v", uuid.New().String(), header.Filename)

	// file must be read to byte buffer before detecting content type
	buffer := make([]byte, header.Size)
	file.Read(buffer)
	fileType := http.DetectContentType(buffer)
	// return the file position to the start before uploading
	file.Seek(0, 0)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(i.bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(fileType),
	})
	if err != nil {
		return "", err
	}

	return key, nil
}

func (i *ImageStoreS3) Validate(file multipart.File) (bool, error) {
	// TODO: actually check the image file is OK
	return true, nil
}

func (i *ImageStoreS3) AddImageToTransaction(transaction app.Transaction) (app.Transaction, error) {
	if transaction.ImageKey == "" {
		return app.Transaction{}, errors.New("transaction is missing imageKey")
	}

	svc := s3.New(i.session)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(i.bucket),
		Key:    aws.String(transaction.ImageKey),
	})

	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		return app.Transaction{}, errors.New(fmt.Sprintf("failed to sign image URL: %v", err))
	}

	transaction.ImageURL = urlStr
	return transaction, nil
}
