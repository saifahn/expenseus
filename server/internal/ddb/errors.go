package ddb

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/pkg/errors"
)

var ErrConflict = errors.New("dynamodb: conflict")

func conflictOrErr(err error) error {
	dynamoErr, ok := errors.Cause(err).(awserr.Error)
	if ok && dynamoErr.Code() == "ConditionalCheckFailedException" {
		return ErrConflict
	}
	return err
}