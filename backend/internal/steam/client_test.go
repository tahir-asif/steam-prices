package steam

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchGameDetails(t *testing.T) {
	tests := []struct {
		name           string
		appID          int
		mockResponse   string
		mockStatusCode int
		wantErr        bool
		wantName       string
		wantPrice      int
	}{
		{
			name:         "free game",
			appID:        730,
			mockResponse: `{"730":{"success":true,"data":{"name":"CS2","is_free":true}}}`,
			wantName:     "CS2",
			wantPrice:    0,
		},
		{
			name:         "paid game",
			appID:        292030,
			mockResponse: `{"292030":{"success":true,"data":{"name":"Witcher 3","is_free":false,"price_overview":{"final":3999,"currency":"USD"}}}}`,
			wantName:     "Witcher 3",
			wantPrice:    3999,
		},
		{
			name:           "non-200 status",
			appID:          123,
			mockStatusCode: http.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name:         "success false",
			appID:        999999,
			mockResponse: `{"999999":{"success":false}}`,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.mockStatusCode != 0 {
					w.WriteHeader(tt.mockStatusCode)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, tt.mockResponse)
			}))
			defer server.Close()

			client := &Client{
				HTTPClient: server.Client(),
				BaseURL:    server.URL,
			}

			game, err := client.FetchGameDetails(tt.appID)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
			if err == nil {
				if game.Name != tt.wantName {
					t.Errorf("expected name %q, got %q", tt.wantName, game.Name)
				}
				if game.Price != tt.wantPrice {
					t.Errorf("expected price %d, got %d", tt.wantPrice, game.Price)
				}
			}
		})
	}
}
