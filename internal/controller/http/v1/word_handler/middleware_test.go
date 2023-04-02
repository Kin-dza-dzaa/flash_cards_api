package wordhadnler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_jwtAuthenticator(t *testing.T) {
	h, _ := setupWordHandler(t)
	dmyHand := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		wantRes httpResponse
	}{
		{
			name: "Without jwt in header",
			wantRes: httpResponse{
				Path:    "/jwt",
				Status:  http.StatusUnauthorized,
				Message: http.StatusText(http.StatusUnauthorized),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/jwt", nil),
			},
		},
		{
			name: "With invalid jwt in header",
			wantRes: httpResponse{
				Path:    "/jwt",
				Status:  http.StatusUnauthorized,
				Message: http.StatusText(http.StatusUnauthorized),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/jwt", nil)
					r.Header.Add("Authorization", "BEARER invalid.jwt.bad_jwt")
					return r
				}(),
			},
		},
	}

	for _, tt := range tests {
		h.jwtAuthenticator(dmyHand).ServeHTTP(tt.args.w, tt.args.r)
		var gotResponse httpResponse
		err := json.Unmarshal(tt.args.w.Body.Bytes(), &gotResponse)
		if err != nil {
			t.Fatalf("%v - json.Unmarshal: %v", tt.name, err)
		}
		if diff := cmp.Diff(tt.wantRes, gotResponse); diff != "" {
			t.Fatalf("wanted: %v got: %v dif: %v", tt.wantRes, gotResponse, diff)
		}
	}
}
