package cachedmysql

import "errors"

var (
	ErrKeyNotFound = errors.New("not found")
	ErrGetFailed   = errors.New("get failed")
)
