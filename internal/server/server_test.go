//go:build unit
// +build unit

package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	shortlyError "github.com/burkaydurdu/shortly/pkg/error"

	"github.com/burkaydurdu/shortly/internal/domain/shortly"

	"github.com/burkaydurdu/shortly/internal/db"

	"github.com/burkaydurdu/shortly/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServerSuite struct {
	suite.Suite
	config *config.Config
}

func (s *ServerSuite) SetupSuite() {
	conf, _ := config.New()

	conf.Server.Port = 5454
	conf.MemoryPath = "../../.mem"
	conf.MemoryFileName = "test_shortly"
	s.config = conf

	server := NewServer(s.config)

	go server.Start()
}

func (s *ServerSuite) TearDownSuite() {
	removeDB(s.config)
}

func (s *ServerSuite) SetupTest() {
	removeDB(s.config)
}

func (s *ServerSuite) TestHealthHandler() {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/api/v1/health", s.config.Server.Port))

	body, err := io.ReadAll(resp.Body)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), string(body), "OK")
}

func (s *ServerSuite) TestGetShortList() {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/api/v1/list", s.config.Server.Port))

	body, err := io.ReadAll(resp.Body)

	var response = make([]db.Shortly, 0)

	err = json.Unmarshal(body, &response)

	assert.NoError(s.T(), err)
	assert.Positive(s.T(), len(response))
}

func (s *ServerSuite) TestCreateShortURL() {
	body := shortly.SaveRequestDTO{
		OriginalURL: "http://burkaydurdu.github.io",
	}

	byteBody, err := json.Marshal(&body)

	resp, err := http.Post(
		fmt.Sprintf("http://localhost:%d/api/v1/create", s.config.Server.Port),
		"application/json",
		bytes.NewBuffer(byteBody),
	)

	var response = shortly.SaveResponseDTO{}

	respBody, err := io.ReadAll(resp.Body)

	err = json.Unmarshal(respBody, &response)

	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), response.ShortURL)
}

func (s *ServerSuite) TestCreateShortURL_InvalidParams() {
	var body = map[string]int{
		"original_url": 2,
	}

	byteBody, err := json.Marshal(&body)

	resp, err := http.Post(
		fmt.Sprintf("http://localhost:%d/api/v1/create", s.config.Server.Port),
		"application/json",
		bytes.NewBuffer(byteBody),
	)

	var response = shortly.ErrResponseDTO{}

	respBody, err := io.ReadAll(resp.Body)

	err = json.Unmarshal(respBody, &response)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), resp.StatusCode, 400)
	assert.Equal(s.T(), response.Message, shortlyError.ParserError)
}

func (s *ServerSuite) TestCreateShortURL_InvalidURL() {
	body := shortly.SaveRequestDTO{
		OriginalURL: "http",
	}

	byteBody, err := json.Marshal(&body)

	resp, err := http.Post(
		fmt.Sprintf("http://localhost:%d/api/v1/create", s.config.Server.Port),
		"application/json",
		bytes.NewBuffer(byteBody),
	)

	var response = shortly.ErrResponseDTO{}

	respBody, err := io.ReadAll(resp.Body)

	err = json.Unmarshal(respBody, &response)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), resp.StatusCode, 409)
	assert.Equal(s.T(), response.Message, shortlyError.InvalidURLError)
}

func (s *ServerSuite) TestRedirectURL() {
	body := shortly.SaveRequestDTO{
		OriginalURL: "Http://burkaydurdu.github.io",
	}

	byteBody, _ := json.Marshal(&body)

	resp, err := http.Post(
		fmt.Sprintf("http://localhost:%d/api/v1/create", s.config.Server.Port),
		"application/json",
		bytes.NewBuffer(byteBody),
	)

	assert.NoError(s.T(), err)

	var response = shortly.SaveResponseDTO{}

	respBody, err := io.ReadAll(resp.Body)

	err = json.Unmarshal(respBody, &response)

	assert.NoError(s.T(), err)

	resp, err = http.Get(response.ShortURL)

	respBody, err = io.ReadAll(resp.Body)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), resp.StatusCode, http.StatusOK)
}

func (s *ServerSuite) TestRedirectURL_NotFound() {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/xxxx", s.config.Server.Port))

	respBody, err := io.ReadAll(resp.Body)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), string(respBody), "404 page not found\n")
	assert.Equal(s.T(), resp.StatusCode, http.StatusNotFound)
}

func removeDB(config *config.Config) {
	_ = os.Remove(config.MemoryPath + "/" + config.MemoryFileName + ".json")
}

// Run Test Suit
func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}
