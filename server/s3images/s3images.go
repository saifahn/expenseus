package s3images

import (
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

type ImageStoreS3 struct {
	session  *session.Session
	storeKey string
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

	storeKey := os.Getenv("AWS_IMAGES_BUCKET_KEY")

	return &ImageStoreS3{
		session:  sess,
		storeKey: storeKey,
	}
}

func (i *ImageStoreS3) Upload(file multipart.File) (string, error) {
	uploader := s3manager.NewUploader(i.session)
	// generate the key
	key := uuid.New().String()

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(i.storeKey),
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
