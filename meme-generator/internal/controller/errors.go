package controller

import "errors"

var (
	errDatabase            = errors.New("database error")
	errSessionSave         = errors.New("failed to save session")
	errInvalidIntegerInput = errors.New("invalid integer input")
)
