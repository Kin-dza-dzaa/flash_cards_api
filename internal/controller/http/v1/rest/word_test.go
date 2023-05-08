package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/controller/http/v1/srvmock"
	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	"github.com/Kin-dza-dzaa/flash_cards_api/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/mock"
	"golang.org/x/exp/slog"
)

func setupWordHandler(t *testing.T) (*WordHandler, *srvmock.WordService) {
	t.Helper()
	srvMock := srvmock.NewWordService(t)
	h := &WordHandler{
		wordService: srvMock,
		logger:      logger.New(slog.LevelDebug),
		v:           validator.New(),
	}
	return h, srvMock
}

func Test_userWords(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name      string
		args      args
		wantRes   httpResponse
		setupMock func(srvMock *srvmock.WordService, args args)
	}{
		{
			name: "Without user_id in ctx",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/getWords", nil),
			},
			wantRes: httpResponse{
				Path:    "/getWords",
				Message: http.StatusText(http.StatusUnauthorized),
			},
			setupMock: func(srvMock *srvmock.WordService, args args) {},
		},
		{
			name: "Internal error",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/getWords", nil)
					ctx := inCtx(r.Context(), userIDCtxKey, "12345")
					return r.WithContext(ctx)
				}(),
			},
			wantRes: httpResponse{
				Path:    "/getWords",
				Message: http.StatusText(http.StatusInternalServerError),
			},
			setupMock: func(srvMock *srvmock.WordService, args args) {
				srvMock.On("UserWords", args.r.Context(), mock.Anything).Once().Return(
					nil, errors.New("some internal error"),
				)
			},
		},
	}

	for _, tt := range tests {
		h, srvMock := setupWordHandler(t)
		tt.setupMock(srvMock, tt.args)

		t.Run(tt.name, func(t *testing.T) {
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

func Test_updateLearnInterval(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name      string
		args      args
		wantRes   httpResponse
		setupMock func(srvMock *srvmock.WordService, args args)
	}{
		{
			name: "Invalid json",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/updateLearnInterval", nil)
					ctx := inCtx(r.Context(), userIDCtxKey, "12345")
					return r.WithContext(ctx)
				}(),
			},
			wantRes: httpResponse{
				Path:    "/updateLearnInterval",
				Message: "wrong json format",
			},
			setupMock: func(srvMock *srvmock.WordService, args args) {},
		},
		{
			name: "Empty word",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/updateLearnInterval",
						bytes.NewReader(
							[]byte(
								`
									{
										"word": "",
										"collection_name": "some collection",
										"last_repeat": "2012-04-23T18:25:43.511Z",
										"time_diff": 12351213
									}
								`,
							),
						))
					ctx := inCtx(r.Context(), userIDCtxKey, "12345")
					return r.WithContext(ctx)
				}(),
			},
			wantRes: httpResponse{
				Path:    "/updateLearnInterval",
				Message: http.StatusText(http.StatusBadRequest),
			},
			setupMock: func(srvMock *srvmock.WordService, args args) {},
		},
		{
			name: "Empty last_repeat",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/updateLearnInterval",
						bytes.NewReader(
							[]byte(
								`
									{
										"word": "some word",
										"collection_name": "some collection",
										"time_diff": 12351213
									}
								`,
							),
						))
					ctx := inCtx(r.Context(), userIDCtxKey, "12345")
					return r.WithContext(ctx)
				}(),
			},
			wantRes: httpResponse{
				Path:    "/updateLearnInterval",
				Message: http.StatusText(http.StatusBadRequest),
			},
			setupMock: func(srvMock *srvmock.WordService, args args) {},
		},
		{
			name: "Empty time_diff",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/updateLearnInterval",
						bytes.NewReader(
							[]byte(
								`
									{
										"word": "some word",
										"collection_name": "some collection",
										"last_repeat": "2012-04-23T18:25:43.511Z"
									}
								`,
							),
						))
					ctx := inCtx(r.Context(), userIDCtxKey, "12345")
					return r.WithContext(ctx)
				}(),
			},
			wantRes: httpResponse{
				Path:    "/updateLearnInterval",
				Message: http.StatusText(http.StatusBadRequest),
			},
			setupMock: func(srvMock *srvmock.WordService, args args) {},
		},
		{
			name: "Empty collection_name",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/updateLearnInterval",
						bytes.NewReader(
							[]byte(
								`
									{
										"word": "some_word",
										"collection_name": "",
										"last_repeat": "2012-04-23T18:25:43.511Z",
										"time_diff": 12351213
									}
								`,
							),
						))
					ctx := inCtx(r.Context(), userIDCtxKey, "12345")
					return r.WithContext(ctx)
				}(),
			},
			wantRes: httpResponse{
				Path:    "/updateLearnInterval",
				Message: http.StatusText(http.StatusBadRequest),
			},
			setupMock: func(srvMock *srvmock.WordService, args args) {},
		},
		{
			name: "Without user_id in ctx error",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/updateLearnInterval",
					bytes.NewReader(
						[]byte(
							`
								{
									"word": "some_word",
									"collection_name": "valid_coll",
									"last_repeat": "2012-04-23T18:25:43.511Z",
									"time_diff": 12351213
								}
							`,
						),
					)),
			},
			wantRes: httpResponse{
				Path:    "/updateLearnInterval",
				Message: http.StatusText(http.StatusUnauthorized),
			},
			setupMock: func(srvMock *srvmock.WordService, args args) {},
		},
		{
			name: "Internal error",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/updateLearnInterval",
						bytes.NewReader(
							[]byte(
								`
									{
										"word": "some_word",
										"collection_name": "valid_coll",
										"last_repeat": "2012-04-23T18:25:43.511Z",
										"time_diff": 12351213
									}
								`,
							),
						))
					ctx := inCtx(r.Context(), userIDCtxKey, "12345")
					return r.WithContext(ctx)
				}(),
			},
			wantRes: httpResponse{
				Path:    "/updateLearnInterval",
				Message: http.StatusText(http.StatusInternalServerError),
			},
			setupMock: func(srvMock *srvmock.WordService, args args) {
				srvMock.On("UpdateLearnInterval", args.r.Context(), mock.Anything).Once().Return(
					errors.New("some internal error"),
				)
			},
		},
		{
			name: "Valid request",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/updateLearnInterval",
						bytes.NewReader(
							[]byte(
								`
									{
										"word": "some_word",
										"collection_name": "valid_coll",
										"last_repeat": "2012-04-23T18:25:43.511Z",
										"time_diff": 12351213
									}
								`,
							),
						))
					ctx := inCtx(r.Context(), userIDCtxKey, "12345")
					return r.WithContext(ctx)
				}(),
			},
			wantRes: httpResponse{
				Path:    "/updateLearnInterval",
				Message: http.StatusText(http.StatusOK),
			},
			setupMock: func(srvMock *srvmock.WordService, args args) {
				srvMock.On("UpdateLearnInterval", args.r.Context(), mock.Anything).Once().Return(nil)
			},
		},
	}

	for _, tt := range tests {
		h, srvMock := setupWordHandler(t)
		tt.setupMock(srvMock, tt.args)

		t.Run(tt.name, func(t *testing.T) {
			h.updateLearnInterval(tt.args.w, tt.args.r)
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

func Test_deleteWord(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name      string
		args      args
		wantRes   httpResponse
		setupMock func(srvMock *srvmock.WordService, args args)
	}{
		{
			name: "Invalid json",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/deleteWord", nil)
					ctx := inCtx(r.Context(), userIDCtxKey, "12345")
					return r.WithContext(ctx)
				}(),
			},
			wantRes: httpResponse{
				Path:    "/deleteWord",
				Message: "wrong json format",
			},
			setupMock: func(srvMock *srvmock.WordService, args args) {},
		},
		{
			name: "Empty word",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/deleteWord",
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
					ctx := inCtx(r.Context(), userIDCtxKey, "12345")
					return r.WithContext(ctx)
				}(),
			},
			wantRes: httpResponse{
				Path:    "/deleteWord",
				Message: http.StatusText(http.StatusBadRequest),
			},
			setupMock: func(srvMock *srvmock.WordService, args args) {},
		},
		{
			name: "Empty collection name",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/deleteWord",
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
					ctx := inCtx(r.Context(), userIDCtxKey, "12345")
					return r.WithContext(ctx)
				}(),
			},
			wantRes: httpResponse{
				Path:    "/deleteWord",
				Message: http.StatusText(http.StatusBadRequest),
			},
			setupMock: func(srvMock *srvmock.WordService, args args) {},
		},
		{
			name: "Without user_id in ctx error",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/deleteWord",
					bytes.NewReader(
						[]byte(
							`
								{
									"word": "some_word",
									"collection_name": "valid_coll",
									"last_repeat": "2012-04-23T18:25:43.511Z",
									"time_diff": 12351213
								}
							`,
						),
					)),
			},
			wantRes: httpResponse{
				Path:    "/deleteWord",
				Message: http.StatusText(http.StatusUnauthorized),
			},
			setupMock: func(srvMock *srvmock.WordService, args args) {},
		},
		{
			name: "Internal error",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/deleteWord",
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
					ctx := inCtx(r.Context(), userIDCtxKey, "12345")
					return r.WithContext(ctx)
				}(),
			},
			wantRes: httpResponse{
				Path:    "/deleteWord",
				Message: http.StatusText(http.StatusInternalServerError),
			},
			setupMock: func(srvMock *srvmock.WordService, args args) {
				srvMock.On("DeleteWord", args.r.Context(), mock.Anything).Once().Return(
					errors.New("some internal error"),
				)
			},
		},
		{
			name: "Valid request",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/deleteWord",
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
					ctx := inCtx(r.Context(), userIDCtxKey, "12345")
					return r.WithContext(ctx)
				}(),
			},
			wantRes: httpResponse{
				Path:    "/deleteWord",
				Message: http.StatusText(http.StatusOK),
			},
			setupMock: func(srvMock *srvmock.WordService, args args) {
				srvMock.On("DeleteWord", args.r.Context(), mock.Anything).Once().Return(nil)
			},
		},
	}

	for _, tt := range tests {
		h, srvMock := setupWordHandler(t)
		tt.setupMock(srvMock, tt.args)

		t.Run(tt.name, func(t *testing.T) {
			h.deleteWord(tt.args.w, tt.args.r)
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

func Test_addWord(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name      string
		args      args
		wantRes   httpResponse
		setupMock func(srvMock *srvmock.WordService, args args)
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
										"collection_name": "some collection",
										"last_repeat": "2012-04-23T18:25:43.511Z",
										"time_diff": 12351213
									}
								`,
							),
						))
					ctx := inCtx(r.Context(), userIDCtxKey, "12345")
					return r.WithContext(ctx)
				}(),
			},
			wantRes: httpResponse{
				Path:    "/addWord",
				Message: http.StatusText(http.StatusCreated),
			},
			setupMock: func(srvMock *srvmock.WordService, args args) {
				srvMock.On("AddWord", args.r.Context(), mock.Anything).Once().Return(nil)
			},
		},
		{
			name: "Invalid json",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/addWord", nil)
					ctx := inCtx(r.Context(), userIDCtxKey, "12345")
					return r.WithContext(ctx)
				}(),
			},
			wantRes: httpResponse{
				Path:    "/addWord",
				Message: "wrong json format",
			},
			setupMock: func(srvMock *srvmock.WordService, args args) {},
		},
		{
			name: "Empty last_repeat",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/addWord",
						bytes.NewReader(
							[]byte(
								`
									{
										"word": "some_word",
										"collection_name": "some_coll",
										"time_diff": 12351213
									}
								`,
							),
						))
					ctx := inCtx(r.Context(), userIDCtxKey, "12345")
					return r.WithContext(ctx)
				}(),
			},
			wantRes: httpResponse{
				Path:    "/addWord",
				Message: http.StatusText(http.StatusBadRequest),
			},
			setupMock: func(srvMock *srvmock.WordService, args args) {},
		},
		{
			name: "Empty word",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/addWord",
						bytes.NewReader(
							[]byte(
								`
									{
										"word": "",
										"collection_name": "some collection",
										"last_repeat": "2012-04-23T18:25:43.511Z",
										"time_diff": 12351213
									}
								`,
							),
						))
					ctx := inCtx(r.Context(), userIDCtxKey, "12345")
					return r.WithContext(ctx)
				}(),
			},
			wantRes: httpResponse{
				Path:    "/addWord",
				Message: http.StatusText(http.StatusBadRequest),
			},
			setupMock: func(srvMock *srvmock.WordService, args args) {},
		},
		{
			name: "Empty collection name",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/addWord",
						bytes.NewReader(
							[]byte(
								`
									{
										"word": "some_word",
										"collection_name": "",
										"last_repeat": "2012-04-23T18:25:43.511Z",
										"time_diff": 12351213
									}
								`,
							),
						))
					ctx := inCtx(r.Context(), userIDCtxKey, "12345")
					return r.WithContext(ctx)
				}(),
			},
			wantRes: httpResponse{
				Path:    "/addWord",
				Message: http.StatusText(http.StatusBadRequest),
			},
			setupMock: func(srvMock *srvmock.WordService, args args) {},
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
										"collection_name": "valid_coll",
										"last_repeat": "2012-04-23T18:25:43.511Z",
										"time_diff": 12351213
									}
								`,
							),
						))
					ctx := inCtx(r.Context(), userIDCtxKey, "12345")
					return r.WithContext(ctx)
				}(),
			},
			wantRes: httpResponse{
				Path:    "/addWord",
				Message: "word not supported",
			},
			setupMock: func(srvMock *srvmock.WordService, args args) {
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
									"collection_name": "valid_coll",
									"last_repeat": "2012-04-23T18:25:43.511Z",
									"time_diff": 12351213
								}
							`,
						),
					)),
			},
			wantRes: httpResponse{
				Path:    "/addWord",
				Message: http.StatusText(http.StatusUnauthorized),
			},
			setupMock: func(srvMock *srvmock.WordService, args args) {},
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
										"collection_name": "valid_coll",
										"last_repeat": "2012-04-23T18:25:43.511Z",
										"time_diff": 12351213
									}
								`,
							),
						))
					ctx := inCtx(r.Context(), userIDCtxKey, "12345")
					return r.WithContext(ctx)
				}(),
			},
			wantRes: httpResponse{
				Path:    "/addWord",
				Message: http.StatusText(http.StatusInternalServerError),
			},
			setupMock: func(srvMock *srvmock.WordService, args args) {
				srvMock.On("AddWord", args.r.Context(), mock.Anything).Once().
					Return(errors.New("deep test repo internal error"))
			},
		},
	}

	for _, tt := range tests {
		h, srvMock := setupWordHandler(t)
		tt.setupMock(srvMock, tt.args)

		t.Run(tt.name, func(t *testing.T) {
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
