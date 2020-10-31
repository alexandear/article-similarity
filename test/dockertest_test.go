package test

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/pkg/errors"
)

func defaultRetryFunc(port int) func() error {
	return func() error {
		url := fmt.Sprintf("http://0.0.0.0:%d", port)
		res, err := http.Get(url)
		if err != nil {
			return errors.Wrapf(err, "failed to execute a gen request to [%s]", url)
		}
		if err := res.Body.Close(); err != nil {
			return errors.Wrap(err, "failed to close the response body")
		}
		if res.StatusCode != http.StatusOK {
			return errors.Errorf("status code [%d] must be 200", res.StatusCode)
		}

		return nil
	}
}

func freePort() int {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}

	list, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	if err := list.Close(); err != nil {
		log.Fatal(err)
	}

	return list.Addr().(*net.TCPAddr).Port
}
