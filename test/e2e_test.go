package test

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"

	"github.com/devchallenge/article-similarity/internal/util"
)

type e2eTestSuite struct {
	suite.Suite
	serverPort string
	pool       *dockertest.Pool
	resources  []*dockertest.Resource
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, &e2eTestSuite{})
}

func (s *e2eTestSuite) SetupSuite() {
	rand.Seed(time.Now().Unix())

	pool, err := dockertest.NewPool("")
	s.Require().NoError(err)
	s.pool = pool

	s.Require().NoError(s.setupServer(s.pool))
}

func (s *e2eTestSuite) setupServer(pool *dockertest.Pool) error {
	ex, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "failed to get current directory")
	}
	dockerfilePath := path.Join(filepath.Dir(ex), "Dockerfile")

	port := freePort()
	portStr := strconv.Itoa(port)

	options := &dockertest.RunOptions{
		Name:         "article-similarity",
		Repository:   "article-similarity",
		Tag:          "latest",
		ExposedPorts: []string{portStr},
		Env:          []string{"PORT=" + portStr},
		PortBindings: map[docker.Port][]docker.PortBinding{docker.Port(portStr): {{HostPort: portStr}}},
	}

	resource, err := pool.BuildAndRunWithOptions(dockerfilePath, options)
	if err != nil {
		return errors.Wrap(err, "failed to build and run")
	}

	if err := pool.Retry(defaultRetryFunc(port)); err != nil {
		return errors.Wrap(err, "failed to retry")
	}

	s.resources = append(s.resources, resource)
	s.serverPort = portStr
	return nil
}

func (s *e2eTestSuite) TearDownSuite() {
	for _, r := range s.resources {
		s.NoError(s.pool.Purge(r))
	}
}

func (s *e2eTestSuite) Test_EndToEnd() {
	// POST /articles
	reqFirst := s.NewRequest(http.MethodPost, "/articles", `{"content":"hello world"}`)
	respFirst := s.DoRequest(reqFirst)
	s.EqualResponse(http.StatusCreated, `{"content":"hello world","duplicate_article_ids":[],"id":0}`, respFirst)

	// POST /articles
	reqDuplicate := s.NewRequest(http.MethodPost, "/articles", `{"content":"Hello a world!"}`)
	respDuplicate := s.DoRequest(reqDuplicate)
	s.EqualResponse(http.StatusCreated, `{"content":"Hello a world!","duplicate_article_ids":[0],"id":1}`, respDuplicate)

	// POST /articles
	reqUnique := s.NewRequest(http.MethodPost, "/articles", `{"content":"unique"}`)
	respUnique := s.DoRequest(reqUnique)
	s.EqualResponse(http.StatusCreated, `{"content":"unique","duplicate_article_ids":[],"id":2}`, respUnique)

	// GET /articles/1
	reqGet := s.NewRequest(http.MethodGet, "/articles/1", "")
	respGet := s.DoRequest(reqGet)
	s.EqualResponse(http.StatusOK, `{"content":"Hello a world!","duplicate_article_ids":[0],"id":1}`, respGet)
}

func (s *e2eTestSuite) NewRequest(method, path, body string) *http.Request {
	req, err := http.NewRequest(method, fmt.Sprintf("http://localhost:%s%s", s.serverPort, path), strings.NewReader(body))
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")

	return req
}

func (s *e2eTestSuite) DoRequest(req *http.Request) *http.Response {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Do(req)
	s.Require().NoError(err)

	return resp
}

func (s *e2eTestSuite) EqualResponse(expectedStatusCode int, expectedBody string, actual *http.Response) {
	s.Require().NotNil(actual)
	s.Require().NotNil(actual.Body)
	defer util.Close(actual.Body)

	s.Equal(expectedStatusCode, actual.StatusCode)

	byteBody, err := ioutil.ReadAll(actual.Body)
	s.Require().NoError(err)
	s.Equal(expectedBody, strings.Trim(string(byteBody), "\n"))
}
