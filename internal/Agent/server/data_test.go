package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestAgentServer_PostCredentials(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	body, err := json.Marshal(RespData{})
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/credentials" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err = server.PostCredentials(context.TODO(), &ReqData{})
	assert.NoError(t, err)

}
func TestAgentServer_PostCredentials404(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	body, err := json.Marshal(RespData{})
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/credentialss" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err = server.PostCredentials(context.TODO(), &ReqData{})
	assert.Error(t, err)

}
func TestAgentServer_PostCredentials500BadJson(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/credentials" {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				//Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header: make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err := server.PostCredentials(context.TODO(), &ReqData{})
	assert.Error(t, err)

}
func TestAgentServer_PostCredentials500(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	body, err := json.Marshal(RespError{})
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/credentials" {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err = server.PostCredentials(context.TODO(), &ReqData{})
	assert.Error(t, err)

}
func TestAgentServer_PostCredentialsBadBody(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/credentials" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer([]byte("test"))),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err := server.PostCredentials(context.TODO(), &ReqData{})
	assert.Error(t, err)

}

func TestAgentServer_PostCrateFileStartChunks(t *testing.T) {
	ctx := context.Background()
	data := []byte("test data")
	fileName := "test.txt"
	uuidChunk := "test-uuid"
	nStart := 0
	nEnd := 5
	maxSize := 10
	reqData := []byte(`{"info": "test"}`)

	client := resty.New()
	server := &AgentServer{
		host:     "http://localhost:8080",
		JWTToken: "test-token",
		client:   client,
	}

	body, err := json.Marshal(RespData{})
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/fileChunks" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, respData, err := server.PostCrateFileStartChunks(ctx, data, fileName, uuidChunk, nStart, nEnd, maxSize, reqData)

	assert.NoError(t, err)
	assert.NotNil(t, respData)

}
func TestAgentServer_PostCrateFileEndChunks(t *testing.T) {
	ctx := context.Background()
	data := []byte("test data")
	fileName := "test.txt"
	uuidChunk := "test-uuid"
	nStart := 5
	nEnd := 5
	maxSize := 5
	reqData := []byte(`{"info": "test"}`)

	client := resty.New()
	server := &AgentServer{
		host:     "http://localhost:8080",
		JWTToken: "test-token",
		client:   client,
	}

	body, err := json.Marshal(RespData{})
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/fileChunks" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, respData, err := server.PostCrateFileStartChunks(ctx, data, fileName, uuidChunk, nStart, nEnd, maxSize, reqData)

	assert.NoError(t, err)
	assert.NotNil(t, respData)
}
func TestAgentServer_PostCrateFileBad400(t *testing.T) {
	ctx := context.Background()
	data := []byte("test data")
	fileName := "test.txt"
	uuidChunk := "test-uuid"
	nStart := 5
	nEnd := 5
	maxSize := 5
	reqData := []byte(`{"info": "test"}`)

	client := resty.New()
	server := &AgentServer{
		host:     "http://localhost:8080",
		JWTToken: "test-token",
		client:   client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, _, err := server.PostCrateFileStartChunks(ctx, data, fileName, uuidChunk, nStart, nEnd, maxSize, reqData)
	assert.Error(t, err)
}
func TestAgentServer_PostCrateFileEndChunksBadReq(t *testing.T) {
	ctx := context.Background()
	data := []byte("test data")
	fileName := "test.txt"
	uuidChunk := "test-uuid"
	nStart := 5
	nEnd := 5
	maxSize := 5
	reqData := []byte(`{"info": "test"}`)

	client := resty.New()
	server := &AgentServer{
		host:     "http://localhost:8080",
		JWTToken: "test-token",
		client:   client,
	}

	body, err := json.Marshal(RespError{})
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/fileChunks" {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, _, err = server.PostCrateFileStartChunks(ctx, data, fileName, uuidChunk, nStart, nEnd, maxSize, reqData)

	assert.Error(t, err)

}
func TestAgentServer_PostCrateFileEndChunksBadReqJson(t *testing.T) {
	ctx := context.Background()
	data := []byte("test data")
	fileName := "test.txt"
	uuidChunk := "test-uuid"
	nStart := 5
	nEnd := 5
	maxSize := 5
	reqData := []byte(`{"info": "test"}`)

	client := resty.New()
	server := &AgentServer{
		host:     "http://localhost:8080",
		JWTToken: "test-token",
		client:   client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/fileChunks" {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				//Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header: make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, _, err := server.PostCrateFileStartChunks(ctx, data, fileName, uuidChunk, nStart, nEnd, maxSize, reqData)

	assert.Error(t, err)

}
func TestAgentServer_PostCrateFileStartChunks2(t *testing.T) {
	ctx := context.Background()
	data := []byte("test data")
	fileName := "test.txt"
	uuidChunk := "test-uuid"
	nStart := 0
	nEnd := 5
	maxSize := 10
	reqData := []byte(`{"info": "test"}`)

	client := resty.New()
	server := &AgentServer{
		host:     "http://localhost:8080",
		JWTToken: "test-token",
		client:   client,
	}

	//body, err := json.Marshal(RespData{})
	//assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/fileChunks" {
			return &http.Response{
				StatusCode: http.StatusOK,
				//Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header: make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, _, err := server.PostCrateFileStartChunks(ctx, data, fileName, uuidChunk, nStart, nEnd, maxSize, reqData)

	assert.NoError(t, err)

}
func TestAgentServer_PostCrateFileStartChunksBadJson(t *testing.T) {
	ctx := context.Background()
	data := []byte("test data")
	fileName := "test.txt"
	uuidChunk := "test-uuid"
	nStart := 0
	nEnd := 5
	maxSize := 10
	reqData := []byte(`{"info": "test"}`)

	client := resty.New()
	server := &AgentServer{
		host:     "http://localhost:8080",
		JWTToken: "test-token",
		client:   client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/fileChunks" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer([]byte("saddasds"))),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, _, err := server.PostCrateFileStartChunks(ctx, data, fileName, uuidChunk, nStart, nEnd, maxSize, reqData)

	assert.Error(t, err)

}

func TestAgentServer_PostCrateFile(t *testing.T) {

}
