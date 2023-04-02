package wordhadnler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	"github.com/google/go-cmp/cmp"
)

func Test_decodeCollection(t *testing.T) {
	h, _ := setupWordHandler(t)

	type args struct {
		r *http.Request
	}
	tests := []struct {
		wantColl entity.Collection
		name     string
		wantErr  bool
		args     args
	}{
		{
			name: "Valid json",
			wantColl: entity.Collection{
				Word: "some_word",
				Name: "some_coll",
			},
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/decode",
					bytes.NewReader(
						[]byte(
							`
							{
								"word": "some_word",
								"collection_name": "some_coll"
							}
							`,
						),
					),
				),
			},
		},
		{
			name: "Invalid json",
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/decode",
					bytes.NewReader(
						[]byte(
							`
							{
								"word": "some_word,
								"col!@#!#@
							}
						`,
						),
					)),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotColl, err := h.decodeCollection(tt.args.r)
			if tt.wantErr && err == nil {
				t.Fatalf("want err but got: %v", err)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("want nil but got: %v", err)
			}
			if diff := cmp.Diff(gotColl, tt.wantColl); diff != "" {
				t.Fatalf("wanted: %v but got: %v diff: %v", tt.wantColl, gotColl, diff)
			}
		})
	}
}

func Test_encodeResopnse(t *testing.T) {
	h, _ := setupWordHandler(t)

	type args struct {
		w        *httptest.ResponseRecorder
		response httpResponse
	}
	tests := []struct {
		wantRes httpResponse
		name    string
		args    args
	}{
		{
			name: "Valid response",
			args: args{
				w: httptest.NewRecorder(),
				response: httpResponse{
					Path:    "/some_path",
					Status:  200,
					Message: "test",
				},
			},
			wantRes: httpResponse{
				Path:    "/some_path",
				Status:  200,
				Message: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h.encodeResponse(tt.args.w, tt.args.response)
			var gotResponse httpResponse
			err := json.Unmarshal(tt.args.w.Body.Bytes(), &gotResponse)
			if err != nil {
				t.Fatalf("%v - json.Unmarshal: %v", tt.name, err)
			}
			if diff := cmp.Diff(tt.wantRes, gotResponse); diff != "" {
				t.Fatalf("wanted: %v got: %v dif: %v", tt.wantRes, gotResponse, diff)
			}
		})
	}
}
