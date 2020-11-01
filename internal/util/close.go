package util

import (
	"io"
)

func Close(closer io.Closer, logger func(string, ...interface{})) {
	if err := closer.Close(); err != nil {
		logger("close failed: %v", err)
	}
}
