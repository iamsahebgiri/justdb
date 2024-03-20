package justdb

import "errors"

var (
	ErrNoKey            = errors.New("invalid key: key is either deleted or not set yet")
	ErrChecksumMismatch = errors.New("checksum mismatch: data is corrupted")
)
