package wordhadnler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/mock"
)

func Test_addWord(t *testing.T) {
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
			name: "Valid request",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/addWord",
						bytes.NewReader(
							[]byte(
								`
									{
										"word": "some_word",
										"collection_name": "some collection"
									}
								`,
							),
						))
					ctx := inCtx(r.Context(), "user_id", "12345")
					return r.WithContext(ctx)
				}(),
			},
			wantRes: httpResponse{
				Path:    "/addWord",
				Status:  http.StatusOK,
				Message: http.StatusText(http.StatusOK),
			},
			setupMock: func(args args) {
				srvMock.On("AddWord", args.r.Context(), mock.Anything).Once().Return(nil)
			},
		},
		{
			name: "Invalid json",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/addWord", nil)
					ctx := inCtx(r.Context(), "user_id", "12345")
					return r.WithContext(ctx)
				}(),
			},
			wantRes: httpResponse{
				Path:    "/addWord",
				Status:  http.StatusBadRequest,
				Message: "wrong json format",
			},
			setupMock: func(args args) {},
		},
		{
			name: "Invalid user_id",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/addWord",
						bytes.NewReader(
							[]byte(
								`
									{
										"word": "some_word",
										"collection_name": "some collection"
									}
								`,
							),
						))
					ctx := inCtx(r.Context(), "user_id", "")
					return r.WithContext(ctx)
				}(),
			},
			wantRes: httpResponse{
				Path:    "/addWord",
				Status:  http.StatusBadRequest,
				Message: http.StatusText(http.StatusBadRequest),
			},
			setupMock: func(args args) {},
		},
		{
			name: "Invalid word",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/addWord",
						bytes.NewReader(
							[]byte(
								`
									{
										"word": "",
										"collection_name": "some collection"
									}
								`,
							),
						))
					ctx := inCtx(r.Context(), "user_id", "12345")
					return r.WithContext(ctx)
				}(),
			},
			wantRes: httpResponse{
				Path:    "/addWord",
				Status:  http.StatusBadRequest,
				Message: http.StatusText(http.StatusBadRequest),
			},
			setupMock: func(args args) {},
		},
		{
			name: "Invalid collection name",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/addWord",
						bytes.NewReader(
							[]byte(
								`
									{
										"word": "some_word",
										"collection_name": ""
									}
								`,
							),
						))
					ctx := inCtx(r.Context(), "user_id", "")
					return r.WithContext(ctx)
				}(),
			},
			wantRes: httpResponse{
				Path:    "/addWord",
				Status:  http.StatusBadRequest,
				Message: http.StatusText(http.StatusBadRequest),
			},
			setupMock: func(args args) {},
		},
		{
			name: "Not supported word error",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/addWord",
						bytes.NewReader(
							[]byte(
								`
									{
										"word": "bad word 12123@#!@!@3",
										"collection_name": "valid_coll"
									}
								`,
							),
						))
					ctx := inCtx(r.Context(), "user_id", "12345")
					return r.WithContext(ctx)
				}(),
			},
			wantRes: httpResponse{
				Path:    "/addWord",
				Status:  http.StatusBadRequest,
				Message: "word not supported",
			},
			setupMock: func(args args) {
				srvMock.On("AddWord", mock.Anything, mock.Anything).Once().
					Return(entity.ErrWordNotSupported)
			},
		},
		{
			name: "Without user_id in ctx error",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/addWord",
					bytes.NewReader(
						[]byte(
							`
								{
									"word": "some_word",
									"collection_name": "valid_coll"
								}
							`,
						),
					)),
			},
			wantRes: httpResponse{
				Path:    "/addWord",
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
					r := httptest.NewRequest(http.MethodGet, "/addWord",
						bytes.NewReader(
							[]byte(
								`
									{
										"word": "some_word",
										"collection_name": "valid_coll"
									}
								`,
							),
						))
					ctx := inCtx(r.Context(), "user_id", "12345")
					return r.WithContext(ctx)
				}(),
			},
			wantRes: httpResponse{
				Path:    "/addWord",
				Status:  http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
			},
			setupMock: func(args args) {
				srvMock.On("AddWord", args.r.Context(), mock.Anything).Once().
					Return(errors.New("deep test repo internal error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(tt.args)
			h.addWord(tt.args.w, tt.args.r)
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
