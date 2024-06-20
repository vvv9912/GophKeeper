package handler

import (
	mock_service "GophKeeper/internal/Server/service/mocks"
	"context"
	"testing"
)

// todo
func TestHandler_HandlerSignUp(t *testing.T) {
	type mockFunc func(s *mock_service.MockData, ctx context.Context, login, password string)
	type args struct {
		mockFunc
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "ok",
			args: args{
				mockFunc: func(s *mock_service.MockAuth, ctx context.Context, login, password string) {
					s.EXPECT().SignUp(ctx, login, password).Return(nil)
				},
			},
		},
		{
			name: "fail",
			args: args{
				mockFunc: func(s *mock_service.MockAuth, ctx context.Context, login, password string) {
					s.EXPECT().SignUp(ctx, login, password).Return(errors.New("fail"))
				},
			},
		},
	}
}
