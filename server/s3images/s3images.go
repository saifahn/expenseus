package s3images

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/saifahn/expenseus"
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
			Endpoint:         aws.String(os.Getenv("AWS_ENDPOINT")),
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

func (i *ImageStoreS3) Upload(file multipart.File) (string, error) {
	uploader := s3manager.NewUploader(i.session)
	// generate the key
	key := uuid.New().String()

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(i.bucket),
		Key:    aws.String(key),
		Body:   file,
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

func (i *ImageStoreS3) AddImageToExpense(expense expenseus.Expense) (expenseus.Expense, error) {
	if expense.ImageKey == "" {
		return expenseus.Expense{}, errors.New("expense is missing imageKey")
	}

	svc := s3.New(i.session)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(i.bucket),
		Key:    aws.String(expense.ImageKey),
	})

	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		return expenseus.Expense{}, errors.New(fmt.Sprintf("failed to sign image URL: %v", err))
	}

	expense.ImageURL = urlStr
	return expense, nil
}
