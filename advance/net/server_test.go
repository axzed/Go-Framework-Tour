package net

import (
	"testing"
)

func TestServer_StartAndServe(t *testing.T) {
	s := &Server{
		addr: ":8080",
	}
	_ = s.StartAndServe()
}

func TestServe(t *testing.T) {
	_ = Serve(":8080")
}
