package server

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestAgentServer_SetJWTToken(t *testing.T) {
	a := AgentServer{}
	a.SetJWTToken("token")

	if a.JWTToken != "token" {
		t.Errorf("got %s, want %s", a.JWTToken, "token")
	}
}
func TestAgentServer_GetJWTToken(t *testing.T) {
	a := AgentServer{}
	a.JWTToken = "token"

	if a.GetJWTToken() != "token" {
		t.Errorf("got %s, want %s", a.JWTToken, "token")
	}
}

type MockTransport struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

// Метод RoundTrip для удовлетворения интерфейса http.RoundTripper
func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}
func TestSignIn_SuccessfulSignInReturnsUser(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/signIn" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"login": "testuser", "jwt": "testtoken"}`)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	user, err := server.SignIn(context.TODO(), "testuser", "testpassword")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Login)
}
func TestSignIn_SuccessfulSignInReturnsBadUser(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/signIn" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`nil`)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err := server.SignIn(context.TODO(), "testuser", "testpassword")
	assert.Error(t, err)
}
func TestSignIn_InvalidCredentialsReturnsError(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	// Mock the response
	client.SetTransport(&MockTransport{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.String() == "http://localhost:8080/api/signIn" {
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"code": 401, "message": "invalid credentials"}`)),
					Header:     make(http.Header),
				}, nil
			}
			return nil, fmt.Errorf("unexpected request")
		},
	})

	user, err := server.SignIn(context.Background(), "wronguser", "wrongpassword")
	assert.Error(t, err)
	assert.Nil(t, user)
}
func TestSignIn_ServerReturnsNon200StatusCode(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	// Mock the response
	client.SetTransport(&MockTransport{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.String() == "http://localhost:8080/api/signIn" {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"code": 500, "message": "internal server error"}`)),
					Header:     make(http.Header),
				}, nil
			}
			return nil, fmt.Errorf("unexpected request")
		},
	})

	user, err := server.SignIn(context.Background(), "testuser", "testpassword")
	assert.Error(t, err)
	assert.Nil(t, user)
}
func TestSignIn_ServerReturnsNon404(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	// Mock the response
	client.SetTransport(&MockTransport{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.String() == "http://localhost:8080/api/l/signIn" {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"code": 500, "message": "internal server error"}`)),
					Header:     make(http.Header),
				}, nil
			}
			return nil, fmt.Errorf("unexpected request")
		},
	})

	user, err := server.SignIn(context.Background(), "testuser", "testpassword")
	assert.Error(t, err)
	assert.Nil(t, user)
}
func TestSignIn_ServerReturns404(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	// Mock the response
	client.SetTransport(&MockTransport{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.String() == "http://localhost:8080/api/signIn" {
				return &http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`400`)),
					Header:     make(http.Header),
				}, nil
			}
			return nil, fmt.Errorf("unexpected request")
		},
	})

	user, err := server.SignIn(context.Background(), "testuser", "testpassword")
	assert.Error(t, err)
	assert.Nil(t, user)
}
func TestAgentServer_SignUp(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/signUp" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"login": "testuser", "jwt": "testtoken"}`)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	user, err := server.SignUp(context.TODO(), "testuser", "testpassword")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Login)
}
func TestAgentServer_SignUpBadUser(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/signUp" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`er", "jwt": esttoken"}`)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err := server.SignUp(context.TODO(), "testuser", "testpassword")
	assert.Error(t, err)
}
func TestAgentServer_SignUpBadReq(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/signUp" {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`er", "jwt": esttoken"}`)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err := server.SignUp(context.TODO(), "testuser", "testpassword")
	assert.Error(t, err)
}
func TestAgentServer_SignUp404(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/a/signUp" {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"login": "testuser", "jwt": "testtoken"}`)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err := server.SignUp(context.TODO(), "testuser", "testpassword")
	assert.Error(t, err)
}
func TestAgentServer_SignUp500(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/signUp" {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"login": "testuser", "jwt": "testtoken"}`)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	_, err := server.SignUp(context.TODO(), "testuser", "testpassword")
	assert.Error(t, err)
}

func TestAgentServer_Ping400(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/ping" {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"login": "testuser", "jwt": "testtoken"}`)),
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	err := server.Ping(context.TODO())
	assert.Error(t, err)
}
func TestAgentServer_Ping(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/ping" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	err := server.Ping(context.TODO())
	assert.NoError(t, err)
}
func TestAgentServer_Ping404(t *testing.T) {
	client := resty.New()
	server := &AgentServer{
		host:   "http://localhost:8080",
		client: client,
	}

	rr := func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == "http://localhost:8080/api/pingg" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
			}, nil
		}
		return nil, fmt.Errorf("unexpected request")
	}

	client.SetTransport(&MockTransport{RoundTripFunc: rr})

	err := server.Ping(context.TODO())
	assert.Error(t, err)
}
