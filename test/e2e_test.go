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

func (s *e2eTestSuite) Test_EndToEnd_Ping() {
	// GET / -> OK
	s.AssertRequestResponse(http.MethodGet, "/", ``, http.StatusOK, ``)
}

func (s *e2eTestSuite) Test_EndToEnd_Success() {
	// GET /articles -> 200
	s.AssertRequestResponse(http.MethodGet, "/articles", "",
		http.StatusOK, `{"articles":[]}`)

	// POST /articles {"content": "..."} -> 201
	s.AssertRequestResponse(http.MethodPost, "/articles", `{"content":"first"}`,
		http.StatusCreated, `{"content":"first","duplicate_article_ids":[],"id":1}`)

	// POST /articles {"content": "..."} -> 201
	s.AssertRequestResponse(http.MethodPost, "/articles", `{"content":"First!"}`,
		http.StatusCreated, `{"content":"First!","duplicate_article_ids":[1],"id":2}`)

	// GET /articles/2 -> 200
	s.AssertRequestResponse(http.MethodGet, "/articles/2", "",
		http.StatusOK, `{"content":"First!","duplicate_article_ids":[1],"id":2}`)

	// POST /articles {"content": "..."} -> 201
	s.AssertRequestResponse(http.MethodPost, "/articles", `{"content":"second"}`,
		http.StatusCreated, `{"content":"second","duplicate_article_ids":[],"id":3}`)

	// POST /articles {"content": "..."} -> 201
	s.AssertRequestResponse(http.MethodPost, "/articles", `{"content":"the first"}`,
		http.StatusCreated, `{"content":"the first","duplicate_article_ids":[1,2],"id":4}`)

	// GET /articles/2 -> 200
	s.AssertRequestResponse(http.MethodGet, "/articles/2", "",
		http.StatusOK, `{"content":"First!","duplicate_article_ids":[1,4],"id":2}`)

	// GET /articles -> 200
	s.AssertRequestResponse(http.MethodGet, "/articles", "",
		http.StatusOK, `{"articles":[{"content":"first","duplicate_article_ids":[],"id":1},{"content":"second","duplicate_article_ids":[],"id":3}]}`)

	// POST /articles {"content": "..."} -> 201
	s.AssertRequestResponse(http.MethodPost, "/articles", `{"content":"go go go"}`,
		http.StatusCreated, `{"content":"go go go","duplicate_article_ids":[],"id":5}`)

	// POST /articles {"content": "..."} -> 201
	s.AssertRequestResponse(http.MethodPost, "/articles", `{"content":"go went gone"}`,
		http.StatusCreated, `{"content":"go went gone","duplicate_article_ids":[5],"id":6}`)

	// GET /duplicate_groups -> 200
	s.AssertRequestResponse(http.MethodGet, "/duplicate_groups", ``,
		http.StatusOK, `{"duplicate_groups":[[1,2,4],[5,6]]}`)
}

func (s *e2eTestSuite) Test_EndToEnd_Errors() {
	// GET /articles/abc -> 400
	s.AssertRequestResponse(http.MethodGet, "/articles/abc", ``,
		http.StatusBadRequest, `{"code":601,"message":"id in path must be of type int64: \"abc\""}`)

	// GET /articles/10000 -> 404
	s.AssertRequestResponse(http.MethodGet, "/articles/10000", ``,
		http.StatusNotFound, ``)

	// POST /articles "" -> 400
	s.AssertRequestResponse(http.MethodPost, "/articles", ``,
		http.StatusBadRequest, `{"code":602,"message":"body in body is required"}`)

	// POST /articles {} -> 400
	s.AssertRequestResponse(http.MethodPost, "/articles", `{}`,
		http.StatusBadRequest, `{"code":602,"message":"body.content in body is required"}`)

	// POST /articles {"content": ""} -> 400
	s.AssertRequestResponse(http.MethodPost, "/articles", `{"content": ""}`,
		http.StatusBadRequest, `{"message":"empty content"}`)
}

func (s *e2eTestSuite) AssertRequestResponse(reqMethod, reqPath, reqBody string, expectedStatus int, expectedBody string) {
	req := s.NewRequest(reqMethod, reqPath, reqBody)
	resp := s.DoRequest(req)
	s.EqualResponse(expectedStatus, expectedBody, resp)
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
