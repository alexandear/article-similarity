package util

import (
	"io"
	"log"
)

func Close(closer io.Closer) {
	if err := closer.Close(); err != nil {
		log.Println("close failed: ", err)
	}
}
