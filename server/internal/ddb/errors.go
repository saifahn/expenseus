package ddb

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/pkg/errors"
)

var (
	ErrConflict      = errors.New("dynamodb: conflict")
	ErrAttrNotExists = errors.New("dynamodb: attribute doesn't exist")
	ErrUnexpected    = errors.New("dynamodb: something unexpected occurred")
)

func conflictOrErr(err error) error {
	dynamoErr, ok := errors.Cause(err).(awserr.Error)
	if ok && dynamoErr.Code() == "ConditionalCheckFailedException" {
		return ErrConflict
	}
	return err
}

// used with attribute_exists condition to check if the error is due to the
// condition not being met.
func attrNotExistsOrErr(err error) error {
	dynamoErr, ok := errors.Cause(err).(awserr.Error)
	if ok && dynamoErr.Code() == "ConditionalCheckFailedException" {
		return ErrAttrNotExists
	}
	return err
}
