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
	reqPing := s.NewRequest(http.MethodGet, "/", ``)
	respPing := s.DoRequest(reqPing)
	s.EqualResponse(http.StatusOK, ``, respPing)
}

func (s *e2eTestSuite) Test_EndToEnd_CreateAndGetArticles() {
	// GET /articles -> 200
	reqGetUniqueEmpty := s.NewRequest(http.MethodGet, "/articles", "")
	respGetUniqueEmpty := s.DoRequest(reqGetUniqueEmpty)
	s.EqualResponse(http.StatusOK, `{"articles":[]}`, respGetUniqueEmpty)

	// POST /articles {"content": "..."} -> 201
	reqFirst := s.NewRequest(http.MethodPost, "/articles", `{"content":"hello world"}`)
	respFirst := s.DoRequest(reqFirst)
	s.EqualResponse(http.StatusCreated, `{"content":"hello world","duplicate_article_ids":[],"id":1}`, respFirst)

	// POST /articles {"content": "..."} -> 201
	reqDuplicate := s.NewRequest(http.MethodPost, "/articles", `{"content":"Hello a world!"}`)
	respDuplicate := s.DoRequest(reqDuplicate)
	s.EqualResponse(http.StatusCreated, `{"content":"Hello a world!","duplicate_article_ids":[1],"id":2}`, respDuplicate)

	// POST /articles {"content": "..."} -> 201
	reqUnique := s.NewRequest(http.MethodPost, "/articles", `{"content":"unique"}`)
	respUnique := s.DoRequest(reqUnique)
	s.EqualResponse(http.StatusCreated, `{"content":"unique","duplicate_article_ids":[],"id":3}`, respUnique)

	// GET /articles/1 -> 200
	reqGetID := s.NewRequest(http.MethodGet, "/articles/2", "")
	respGetID := s.DoRequest(reqGetID)
	s.EqualResponse(http.StatusOK, `{"content":"Hello a world!","duplicate_article_ids":[1],"id":2}`, respGetID)

	// GET /articles -> 200
	reqGetUniqueNonEmpty := s.NewRequest(http.MethodGet, "/articles", "")
	respGetUniqueNonEmpty := s.DoRequest(reqGetUniqueNonEmpty)
	s.EqualResponse(http.StatusOK, `{"articles":[{"content":"hello world","duplicate_article_ids":[],"id":1},{"content":"unique","duplicate_article_ids":[],"id":3}]}`, respGetUniqueNonEmpty)
}

func (s *e2eTestSuite) Test_EndToEnd_Errors() {
	// GET /articles/abc -> 400
	reqWrongID := s.NewRequest(http.MethodGet, "/articles/abc", ``)
	respWrongID := s.DoRequest(reqWrongID)
	s.EqualResponse(http.StatusBadRequest, `{"code":601,"message":"id in path must be of type int64: \"abc\""}`, respWrongID)

	// GET /articles/10000 -> 404
	reqNotFound := s.NewRequest(http.MethodGet, "/articles/10000", ``)
	respNotFound := s.DoRequest(reqNotFound)
	s.EqualResponse(http.StatusNotFound, ``, respNotFound)

	// POST /articles "" -> 400
	reqNoBody := s.NewRequest(http.MethodPost, "/articles", ``)
	respNoBody := s.DoRequest(reqNoBody)
	s.EqualResponse(http.StatusBadRequest, `{"code":602,"message":"body in body is required"}`, respNoBody)

	// POST /articles {} -> 400
	reqNoContent := s.NewRequest(http.MethodPost, "/articles", `{}`)
	respNoContent := s.DoRequest(reqNoContent)
	s.EqualResponse(http.StatusBadRequest, `{"code":602,"message":"body.content in body is required"}`, respNoContent)

	// POST /articles {"content": ""} -> 400
	reqEmptyContent := s.NewRequest(http.MethodPost, "/articles", `{"content": ""}`)
	respEmptyContent := s.DoRequest(reqEmptyContent)
	s.EqualResponse(http.StatusBadRequest, `{"message":"empty content"}`, respEmptyContent)
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
