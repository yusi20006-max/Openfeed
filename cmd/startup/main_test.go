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
		{name: "native relative openfeed binary", cmd: "./openfeed", want: true},
		{name: "native absolute openfeed binary", cmd: "/home/user/openfeed", want: true},
		{name: "foreign server", cmd: "/home/user/other-service server", want: false},
		{name: "foreign openfeed-like name", cmd: "/home/user/openfeed-helper", want: false},
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

func TestIsManagedOpenFeedProcess(t *testing.T) {
	root := "/data/data/com.termux/files/home/yasineco/openfeed"
	otherRoot := "/data/data/com.termux/files/home/Openfeed-clean/Openfeed-new"

	tests := []struct {
		name string
		info processInfo
		want bool
	}{
		{
			name: "native binary from same checkout",
			info: processInfo{cmd: "./openfeed", cwd: root},
			want: true,
		},
		{
			name: "native binary from different checkout",
			info: processInfo{cmd: "./openfeed", cwd: otherRoot},
			want: true,
		},
		{
			name: "absolute native binary from different checkout",
			info: processInfo{cmd: otherRoot + "/openfeed", cwd: otherRoot},
			want: true,
		},
		{
			name: "go run from same checkout",
			info: processInfo{cmd: "go run ./cmd/server", cwd: root},
			want: true,
		},
		{
			name: "go run from different checkout",
			info: processInfo{cmd: "go run ./cmd/server", cwd: otherRoot},
			want: false,
		},
		{
			name: "foreign binary from different checkout",
			info: processInfo{cmd: "/home/user/other-service", cwd: otherRoot},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isManagedOpenFeedProcess(root, tt.info); got != tt.want {
				t.Fatalf("isManagedOpenFeedProcess(%q) = %v, want %v", tt.info.cmd, got, tt.want)
			}
		})
	}
}
