package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type e2eTestSuite struct {
	suite.Suite
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, &e2eTestSuite{})
}

func (s *e2eTestSuite) SetupSuite() {
	cmd := exec.Command("docker-compose", "up", "-d")
	s.Require().NoError(cmd.Run())

	time.Sleep(3 * time.Second)
}

func (s *e2eTestSuite) TearDownSuite() {
	s.Require().NoError(exec.Command("docker-compose", "down").Run())
}

func (s *e2eTestSuite) Test_EndToEnd() {
	// GET /
	reqPing := s.NewRequest(http.MethodGet, "/", "")
	respPing := s.DoRequest(reqPing)
	s.EqualResponse(http.StatusOK, ``, respPing)

	// POST /articles
	reqFirst := s.NewRequest(http.MethodPost, "/articles", `{"content":"hello world"}`)
	respFirst := s.DoRequest(reqFirst)
	s.EqualResponse(http.StatusCreated, `{"content":"hello world","duplicate_article_ids":[],"id":1}`, respFirst)

	// POST /articles
	reqDuplicate := s.NewRequest(http.MethodPost, "/articles", `{"content":"Hello a world!"}`)
	respDuplicate := s.DoRequest(reqDuplicate)
	s.EqualResponse(http.StatusCreated, `{"content":"Hello a world!","duplicate_article_ids":[1],"id":2}`, respDuplicate)

	// POST /articles
	reqUnique := s.NewRequest(http.MethodPost, "/articles", `{"content":"unique"}`)
	respUnique := s.DoRequest(reqUnique)
	s.EqualResponse(http.StatusCreated, `{"content":"unique","duplicate_article_ids":[],"id":3}`, respUnique)

	// GET /articles/1
	reqGet := s.NewRequest(http.MethodGet, "/articles/2", "")
	respGet := s.DoRequest(reqGet)
	s.EqualResponse(http.StatusOK, `{"content":"Hello a world!","duplicate_article_ids":[1],"id":2}`, respGet)
}

func (s *e2eTestSuite) NewRequest(method, path, body string) *http.Request {
	req, err := http.NewRequest(method, fmt.Sprintf("http://localhost:80%s", path), strings.NewReader(body))
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

	s.Equal(expectedStatusCode, actual.StatusCode)

	byteBody, err := ioutil.ReadAll(actual.Body)
	s.Require().NoError(err)
	s.Equal(expectedBody, strings.Trim(string(byteBody), "\n"))

	s.Require().NoError(actual.Body.Close())
}
