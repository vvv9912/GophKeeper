package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_getUserId(t *testing.T) {
	type args struct {
		r *http.Request
	}

	tests := []struct {
		name       string
		args       args
		wantUserId int64
		wantErr    bool
	}{
		{
			name: "get user id",
			args: args{
				r: nil,
			},
			wantUserId: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("POST", "/u", nil)
			ctx := context.WithValue(r.Context(), "UserId", tt.wantUserId)
			tt.args.r = r.WithContext(ctx)

			gotUserId, err := getUserId(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("getUserId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserId != tt.wantUserId {
				t.Errorf("getUserId() gotUserId = %v, want %v", gotUserId, tt.wantUserId)
			}
		})
	}
}
