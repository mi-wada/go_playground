package httpbinclient_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mi-wada/go_playground/httpbinclient"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	type test struct {
		name    string
		opts    []httpbinclient.Opt
		wantErr bool
	}

	tests := []test{
		{
			name:    "with no options",
			opts:    nil,
			wantErr: false,
		},
		{
			name: "with valid WithBaseURL",
			opts: []httpbinclient.Opt{
				httpbinclient.WithBaseURL("https://example.com"),
			},
			wantErr: false,
		},
		{
			name: "with invalid WithBaseURL",
			opts: []httpbinclient.Opt{
				httpbinclient.WithBaseURL("://invalid"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, err := httpbinclient.NewClient(nil, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && client == nil {
				t.Errorf("NewClient() client is nil")
			}
		})
	}
}

func TestClient_Get(t *testing.T) {
	t.Parallel()

	httpbinMockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"origin": "hoge", "url": "http://example.com"}`)
	}))
	baseURL := httpbinMockServer.URL

	client, err := httpbinclient.NewClient(http.DefaultClient, httpbinclient.WithBaseURL(baseURL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	resp, err := client.Get(context.Background())

	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if resp.Origin != "hoge" {
		t.Errorf("Get() resp.Origin = %v, want hoge", resp.Origin)
	}
	if resp.URL != "http://example.com" {
		t.Errorf("Get() resp.URL = %v, want http://example.com", resp.URL)
	}
}
