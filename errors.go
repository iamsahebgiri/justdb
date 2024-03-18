package justdb

import "errors"

var (
	ErrNoKey    = errors.New("invalid key: key is either deleted or not set yet")
	ErrChecksum = errors.New("checksum mismatch")
)
