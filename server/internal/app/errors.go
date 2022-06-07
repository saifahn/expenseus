package app

import "errors"

var (
	ErrDBItemNotFound = errors.New("store: item not found")
	ErrUserNotInCtx   = errors.New("user not found in context")
	ErrTxnIDNotInCtx  = errors.New("transactionID not found in context")
)
