// +build integration

package test

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/suite"

	"github.com/devchallenge/article-similarity/internal/restapi"
	"github.com/devchallenge/article-similarity/internal/util"
)

type e2eTestSuite struct {
	suite.Suite
	port   string
	server *restapi.ArticleServer
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, &e2eTestSuite{})
}

func (s *e2eTestSuite) SetupTest() {
	s.port = s.FreePort()
	s.Require().NoError(pflag.Set("port", s.port))

	serv, err := restapi.NewArticleServer()
	s.Require().NoError(err)
	serv.ConfigureLogger(s.T().Logf)
	s.server = serv

	go func() {
		if err := serv.Serve(); err != nil {
			s.NoError(err)
		}
	}()
}

func (s *e2eTestSuite) TearDownTest() {
	s.NoError(s.server.Close())
}

func (s *e2eTestSuite) Test_EndToEnd_AddArticle() {
	s.Run("add two duplicates and one unique", func() {
		reqFirst := s.NewRequest(http.MethodPost, "/articles", `{"content":"hello world"}`)
		respFirst := s.DoRequest(reqFirst)
		s.EqualResponse(http.StatusOK, `{"content":"hello world","duplicate_article_ids":[],"id":0}`, respFirst)

		reqDuplicate := s.NewRequest(http.MethodPost, "/articles", `{"content":"Hello a world!"}`)
		respDuplicate := s.DoRequest(reqDuplicate)
		s.EqualResponse(http.StatusOK, `{"content":"Hello a world!","duplicate_article_ids":[0],"id":1}`, respDuplicate)

		reqUnique := s.NewRequest(http.MethodPost, "/articles", `{"content":"unique"}`)
		respUnique := s.DoRequest(reqUnique)
		s.EqualResponse(http.StatusOK, `{"content":"unique","duplicate_article_ids":[],"id":2}`, respUnique)
	})
}

func (s *e2eTestSuite) Test_EndToEnd_GetArticleByID() {
	s.Run("add two duplicates and get first", func() {
		{
			req := s.NewRequest(http.MethodPost, "/articles", `{"content":"Get article by id."}`)
			resp := s.DoRequest(req)
			s.EqualResponse(http.StatusOK, `{"content":"Get article by id.","duplicate_article_ids":[],"id":0}`, resp)
		}
		{
			req := s.NewRequest(http.MethodPost, "/articles", `{"content":"Get the article by an id."}`)
			resp := s.DoRequest(req)
			s.EqualResponse(http.StatusOK, `{"content":"Get the article by an id.","duplicate_article_ids":[0],"id":1}`, resp)
		}

		reqGet := s.NewRequest(http.MethodGet, "/articles/0", "")
		respGet := s.DoRequest(reqGet)
		s.EqualResponse(http.StatusOK, `{"content":"Get article by id.","duplicate_article_ids":[1],"id":0}`, respGet)
	})
}

func (s *e2eTestSuite) NewRequest(method, path, body string) *http.Request {
	req, err := http.NewRequest(method, fmt.Sprintf("http://localhost:%s%s", s.port, path), strings.NewReader(body))
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

func (s *e2eTestSuite) FreePort() string {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	s.Require().NoError(err)

	list, err := net.ListenTCP("tcp", addr)
	s.Require().NoError(err)
	s.Require().NoError(list.Close())

	return strconv.Itoa(list.Addr().(*net.TCPAddr).Port)
}
