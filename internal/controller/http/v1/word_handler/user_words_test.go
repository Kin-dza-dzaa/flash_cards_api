package wordhadnler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/mock"
)

func Test_userWords(t *testing.T) {
	h, srvMock := setupWordHandler(t)

	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name      string
		args      args
		wantRes   httpResponse
		setupMock func(args args)
	}{
		{
			name: "Without user_id in ctx",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/getWords", nil),
			},
			wantRes: httpResponse{
				Path:    "/getWords",
				Status:  http.StatusUnauthorized,
				Message: http.StatusText(http.StatusUnauthorized),
			},
			setupMock: func(args args) {},
		},
		{
			name: "Internal error",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/getWords", nil)
					ctx := inCtx(r.Context(), "user_id", "12345")
					return r.WithContext(ctx)
				}(),
			},
			wantRes: httpResponse{
				Path:    "/getWords",
				Status:  http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
			},
			setupMock: func(args args) {
				srvMock.On("UserWords", args.r.Context(), mock.Anything).Once().Return(
					nil, errors.New("some internal error"),
				)
			},
		},
		{
			name: "Valid request",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/getWords", nil)
					ctx := inCtx(r.Context(), "user_id", "12345")
					return r.WithContext(ctx)
				}(),
			},
			wantRes: httpResponse{
				Path:    "/getWords",
				Status:  http.StatusOK,
				Message: http.StatusText(http.StatusOK),
				UserWords: &entity.UserWords{
					Words: make(map[entity.CollectionName][]entity.WordData),
				},
			},
			setupMock: func(args args) {
				srvMock.On("UserWords", args.r.Context(), mock.Anything).Once().Return(
					&entity.UserWords{
						Words: make(map[entity.CollectionName][]entity.WordData),
					}, nil,
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(tt.args)
			h.userWords(tt.args.w, tt.args.r)
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
