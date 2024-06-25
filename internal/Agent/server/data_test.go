package server

import (
	"GophKeeper/pkg/store"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
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
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	tt := time.Now()
	body, err := json.Marshal(RespData{
		UserDataId: 1,
		Hash:       "22",
		CreatedAt:  &tt,
		UpdateAt:   &tt,
	})
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/file" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err = server.PostCrateFile(context.TODO(), &ReqData{})
	assert.NoError(t, err)

}

func TestAgentServer_PostCrateFileBadReq(t *testing.T) {
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

	_, err = server.PostCrateFile(context.TODO(), &ReqData{})
	assert.Error(t, err)

}
func TestAgentServer_PostCrateFile400(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/file" {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				//Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header: make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err := server.PostCrateFile(context.TODO(), &ReqData{})
	assert.Error(t, err)

}
func TestAgentServer_PostCrateFile500(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}
	body, err := json.Marshal(RespError{})
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/file" {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err = server.PostCrateFile(context.TODO(), &ReqData{})
	assert.Error(t, err)

}
func TestAgentServer_PostCrateFileBadBody(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/file" {
			return &http.Response{
				StatusCode: http.StatusOK,
				//Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header: make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err := server.PostCrateFile(context.TODO(), &ReqData{})
	assert.Error(t, err)

}

func TestAgentServer_PostCreditCard(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	tt := time.Now()
	body, err := json.Marshal(RespData{
		UserDataId: 1,
		Hash:       "22",
		CreatedAt:  &tt,
		UpdateAt:   &tt,
	})
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/creditCard" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err = server.PostCreditCard(context.TODO(), &ReqData{})
	assert.NoError(t, err)
}
func TestAgentServer_PostCreditCardBadBody(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/creditCard" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err := server.PostCreditCard(context.TODO(), &ReqData{})
	assert.Error(t, err)
}
func TestAgentServer_PostCreditCardBadReq(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	tt := time.Now()
	body, err := json.Marshal(RespData{
		UserDataId: 1,
		Hash:       "22",
		CreatedAt:  &tt,
		UpdateAt:   &tt,
	})
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/creditCsard" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err = server.PostCreditCard(context.TODO(), &ReqData{})
	assert.Error(t, err)
}
func TestAgentServer_PostCreditCard400(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/creditCard" {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err := server.PostCreditCard(context.TODO(), &ReqData{})
	assert.Error(t, err)
}
func TestAgentServer_PostCreditCard500(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}
	body, err := json.Marshal(RespError{})
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/creditCard" {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err = server.PostCreditCard(context.TODO(), &ReqData{})
	assert.Error(t, err)
}

func TestAgentServer_GetCheckChanges(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	tt := time.Now()
	body, err := json.Marshal([]store.UsersData{
		{
			UserDataId: 1,
			Hash:       "22",
			CreatedAt:  &tt,
			UpdateAt:   &tt,
		},
	})
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/changes" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err = server.GetCheckChanges(context.TODO(), &ReqData{}, time.Now())
	assert.NoError(t, err)
}
func TestAgentServer_GetCheckChangesBadJson(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/changes" {
			return &http.Response{
				StatusCode: http.StatusOK,
				//Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header: make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err := server.GetCheckChanges(context.TODO(), &ReqData{}, time.Now())
	assert.Error(t, err)
}
func TestAgentServer_GetCheckChangesBadReq(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	tt := time.Now()
	body, err := json.Marshal([]store.UsersData{
		{
			UserDataId: 1,
			Hash:       "22",
			CreatedAt:  &tt,
			UpdateAt:   &tt,
		},
	})
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/changes1" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err = server.GetCheckChanges(context.TODO(), &ReqData{}, time.Now())
	assert.Error(t, err)
}
func TestAgentServer_GetCheckChanges400(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/changes" {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err := server.GetCheckChanges(context.TODO(), &ReqData{}, time.Now())
	assert.Error(t, err)
}
func TestAgentServer_GetCheckChanges500(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	body, err := json.Marshal(RespError{})
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/changes" {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err = server.GetCheckChanges(context.TODO(), &ReqData{}, time.Now())
	assert.Error(t, err)
}
func TestAgentServer_GetGetData(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	body, err := json.Marshal(RespUsersData{
		InfoUsersData: &store.UsersData{},
		EncryptData:   &store.DataFile{},
	})
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/1" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err = server.GetData(context.TODO(), 1)
	assert.NoError(t, err)
}
func TestAgentServer_GetGetDataTypeBinErr(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	body, err := json.Marshal(RespUsersData{
		InfoUsersData: &store.UsersData{
			DataType: 3,
		},
		EncryptData: &store.DataFile{},
	})
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/1" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err = server.GetData(context.TODO(), 1)
	assert.Error(t, err)
}
func TestAgentServer_GetGetDataBadReq(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	body, err := json.Marshal(RespUsersData{
		InfoUsersData: &store.UsersData{},
		EncryptData:   &store.DataFile{},
	})
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/dat2a/1" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err = server.GetData(context.TODO(), 1)
	assert.Error(t, err)
}
func TestAgentServer_GetGetDataBadJson(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/1" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err := server.GetData(context.TODO(), 1)
	assert.Error(t, err)
}
func TestAgentServer_GetGetData400(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/1" {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err := server.GetData(context.TODO(), 1)
	assert.Error(t, err)
}
func TestAgentServer_GetGetData500(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	body, err := json.Marshal(RespError{})
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/1" {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err = server.GetData(context.TODO(), 1)
	assert.Error(t, err)
}

func TestAgentServer_getFileSize(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	reqFileSize := struct {
		FileSize int64 `json:"fileSize"`
	}{}

	body, err := json.Marshal(reqFileSize)
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/fileSize/1" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err = server.getFileSize(context.TODO(), 1)
	assert.NoError(t, err)
}
func TestAgentServer_getFileSizeBadReq(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	reqFileSize := struct {
		FileSize int64 `json:"fileSize"`
	}{}

	body, err := json.Marshal(reqFileSize)
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/filesSize/1" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err = server.getFileSize(context.TODO(), 1)
	assert.Error(t, err)
}
func TestAgentServer_getFileSize400(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/fileSize/1" {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err := server.getFileSize(context.TODO(), 1)
	assert.Error(t, err)
}
func TestAgentServer_getFileSize500(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	body, err := json.Marshal(RespError{})
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/fileSize/1" {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err = server.getFileSize(context.TODO(), 1)
	assert.Error(t, err)
}
func TestAgentServer_getFileSizeBadBody(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/fileSize/1" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err := server.getFileSize(context.TODO(), 1)
	assert.Error(t, err)
}

func TestAgentServer_getFileData(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	reqFileSize := struct {
		FileSize int64 `json:"fileSize"`
	}{}

	body, err := json.Marshal(reqFileSize)
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/fileChunks/1" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err = server.getFileData(context.TODO(), 1, 10)
	assert.NoError(t, err)
}

func TestAgentServer_getFileChunks(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}
	var (
		userDataId int64 = 1
		fileSize   int64 = 10
		startChunk int   = 1
		endChunk   int   = 1
	)

	reqFileSize := struct {
		FileSize int64 `json:"fileSize"`
	}{}

	body, err := json.Marshal(reqFileSize)
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/fileChunks/1" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err = server.getFileChunks(context.TODO(), userDataId, fileSize, startChunk, endChunk)
	assert.NoError(t, err)
}
func TestAgentServer_getFileChunksBadReq(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}
	var (
		userDataId int64 = 1
		fileSize   int64 = 10
		startChunk int   = 1
		endChunk   int   = 1
	)

	reqFileSize := struct {
		FileSize int64 `json:"fileSize"`
	}{}

	body, err := json.Marshal(reqFileSize)
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/fileChu2nks/1" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err = server.getFileChunks(context.TODO(), userDataId, fileSize, startChunk, endChunk)
	assert.Error(t, err)
}
func TestAgentServer_getFileChunks400(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}
	var (
		userDataId int64 = 1
		fileSize   int64 = 10
		startChunk int   = 1
		endChunk   int   = 1
	)

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/fileChunks/1" {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err := server.getFileChunks(context.TODO(), userDataId, fileSize, startChunk, endChunk)
	assert.Error(t, err)
}
func TestAgentServer_getFileChunks500(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}
	var (
		userDataId int64 = 1
		fileSize   int64 = 10
		startChunk int   = 1
		endChunk   int   = 1
	)

	body, err := json.Marshal(RespError{})
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/fileChunks/1" {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err = server.getFileChunks(context.TODO(), userDataId, fileSize, startChunk, endChunk)
	assert.Error(t, err)
}

func TestAgentServer_GetListData(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	reqFileSize := struct {
		FileSize int64 `json:"fileSize"`
	}{}

	body, err := json.Marshal(reqFileSize)
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err = server.GetListData(context.TODO())
	assert.NoError(t, err)
}
func TestAgentServer_GetListDataBadReq(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	reqFileSize := struct {
		FileSize int64 `json:"fileSize"`
	}{}

	body, err := json.Marshal(reqFileSize)
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/d4ata" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err = server.GetListData(context.TODO())
	assert.Error(t, err)
}
func TestAgentServer_GetListData400(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data" {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err := server.GetListData(context.TODO())
	assert.Error(t, err)
}
func TestAgentServer_GetListData500(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	body, err := json.Marshal(RespError{})
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data" {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err = server.GetListData(context.TODO())
	assert.Error(t, err)
}

func TestAgentServer_CheckUpdate(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	var response struct {
		Status bool `json:"status"`
	}

	body, err := json.Marshal(response)
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/CheckUpdate/1" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	tt := time.Now()
	_, err = server.CheckUpdate(context.TODO(), 1, &tt)
	assert.NoError(t, err)
}
func TestAgentServer_CheckUpdateBadReq(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	var response struct {
		Status bool `json:"status"`
	}

	body, err := json.Marshal(response)
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/CheckUp2date/1" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	tt := time.Now()
	_, err = server.CheckUpdate(context.TODO(), 1, &tt)
	assert.Error(t, err)
}
func TestAgentServer_CheckUpdate400(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/CheckUpdate/1" {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	tt := time.Now()
	_, err := server.CheckUpdate(context.TODO(), 1, &tt)
	assert.Error(t, err)
}

func TestAgentServer_CheckUpdate500(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	body, err := json.Marshal(RespError{})
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/CheckUpdate/1" {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	tt := time.Now()
	_, err = server.CheckUpdate(context.TODO(), 1, &tt)
	assert.Error(t, err)
}
func TestAgentServer_CheckUpdateBadJson(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/CheckUpdate/1" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	tt := time.Now()
	_, err := server.CheckUpdate(context.TODO(), 1, &tt)
	assert.Error(t, err)
}

func TestAgentServer_PostUpdateData(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	body, err := json.Marshal(RespData{})
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/update/1" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err = server.PostUpdateData(context.TODO(), 1, []byte("data"))
	assert.NoError(t, err)
}

func TestAgentServer_PostUpdateDataBadBody(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/update/1" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err := server.PostUpdateData(context.TODO(), 1, []byte("data"))
	assert.Error(t, err)
}
func TestAgentServer_PostUpdateDataBadReq(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	body, err := json.Marshal(RespData{})
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/updae/1" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err = server.PostUpdateData(context.TODO(), 1, []byte("data"))
	assert.Error(t, err)
}
func TestAgentServer_PostUpdateData400(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/update/1" {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err := server.PostUpdateData(context.TODO(), 1, []byte("data"))
	assert.Error(t, err)
}
func TestAgentServer_PostUpdateData500(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	body, err := json.Marshal(RespError{})
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/update/1" {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err = server.PostUpdateData(context.TODO(), 1, []byte("data"))
	assert.Error(t, err)
}

func TestAgentServer_PostUpdateBinaryFile(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	body, err := json.Marshal(RespData{})
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/updateBinary/1" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, _, err = server.PostUpdateBinaryFile(context.TODO(), []byte("data"), "data", "", 0, 0, 0, []byte("Rdata"), 1)
	assert.NoError(t, err)
}
func TestAgentServer_PostUpdateBinaryFileBadReq(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	body, err := json.Marshal(RespData{})
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/updateBi2nary/1" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, _, err = server.PostUpdateBinaryFile(context.TODO(), []byte("data"), "data", "", 0, 0, 0, []byte("Rdata"), 1)
	assert.Error(t, err)
}
func TestAgentServer_PostUpdateBinaryFileSuccess(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	body, err := json.Marshal(RespData{})
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/updateBinary/1" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, _, err = server.PostUpdateBinaryFile(context.TODO(), []byte("data"), "data", "", 0, 0, 1, []byte("Rdata"), 1)
	assert.NoError(t, err)
}
func TestAgentServer_PostUpdateBinaryFileSuccess2(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	body, err := json.Marshal(RespData{})
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/updateBinary/1" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, _, err = server.PostUpdateBinaryFile(context.TODO(), []byte("data"), "data", "112", 0, 0, 1, []byte("Rdata"), 1)
	assert.NoError(t, err)
}
func TestAgentServer_PostUpdateBinaryFile400(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/updateBinary/1" {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, _, err := server.PostUpdateBinaryFile(context.TODO(), []byte("data"), "data", "", 0, 0, 1, []byte("Rdata"), 1)
	assert.Error(t, err)
}
func TestAgentServer_PostUpdateBinaryFile500(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	body, err := json.Marshal(RespError{})
	assert.NoError(t, err)
	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/updateBinary/1" {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, _, err = server.PostUpdateBinaryFile(context.TODO(), []byte("data"), "data", "", 0, 0, 1, []byte("Rdata"), 1)
	assert.Error(t, err)
}
func TestAgentServer_PostUpdateBinaryFileBadJson(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/updateBinary/1" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer([]byte("122"))),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, _, err := server.PostUpdateBinaryFile(context.TODO(), []byte("data"), "data", "112", 0, 0, 1, []byte("Rdata"), 1)
	assert.Error(t, err)
}
func TestAgentServer_PostUpdateBinaryFileBadJson2(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/data/updateBinary/1" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer([]byte(""))),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, _, err := server.PostUpdateBinaryFile(context.TODO(), []byte("data"), "data", "112", 0, 0, 1, []byte("Rdata"), 1)
	assert.NoError(t, err)
}
