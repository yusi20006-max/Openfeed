package main

import "testing"

func TestIsOpenFeedCommand(t *testing.T) {
	tests := []struct {
		name string
		cmd  string
		want bool
	}{
		{name: "go run server", cmd: "go run ./cmd/server", want: true},
		{name: "go build temporary server", cmd: "/tmp/go-build123/exe/server", want: true},
		{name: "openfeed server binary", cmd: "/home/user/openfeed server", want: true},
		{name: "foreign server", cmd: "/home/user/other-service server", want: false},
		{name: "go test", cmd: "go test ./...", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isOpenFeedCommand(tt.cmd); got != tt.want {
				t.Fatalf("isOpenFeedCommand(%q) = %v, want %v", tt.cmd, got, tt.want)
			}
		})
	}
}
