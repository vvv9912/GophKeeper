package logger

import (
	"testing"
)

func TestInitialize(t *testing.T) {
	type args struct {
		level string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantLevel string
	}{
		{
			name:      "#1",
			args:      args{level: "info"},
			wantLevel: "info",
		},
		{
			name:      "#2",
			args:      args{level: "info"},
			wantLevel: "debug",
			wantErr:   true,
		},
		{
			name:      "#3",
			args:      args{level: ""},
			wantLevel: "info",
			wantErr:   false,
		},
		{
			name:      "#3",
			args:      args{level: ""},
			wantLevel: "debug",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("test:", tt.name)
			if err := Initialize(tt.args.level); err != nil {
				t.Errorf("Initialize() error = %v, wantErr %v", err, tt.wantErr)
			}

			if Log == nil {
				t.Error("Log is not initialized")
			}
			if Log.Level().String() != tt.wantLevel {
				if !tt.wantErr {
					t.Error("Log level is not set correctly\n", "\t\t\tcurrent log: ", Log.Level().String(), "\n\t\t\t", "want log", tt.args.level)
				}
			}
		})
	}
}
