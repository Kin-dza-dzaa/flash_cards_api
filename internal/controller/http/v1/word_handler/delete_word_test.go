package wordhadnler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/mock"
)

func (s *wordHandler_Suite) Test_wordHandler_deleteWord() {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		Name      string
		Args      args
		Want      string
		setupMock func()
	}{
		{
			Name: "Invalid json",
			Args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/deleteWord", nil)
					ctx := inCtx(r.Context(), "user_id", "12345")
					return r.WithContext(ctx)
				}(),
			},
			Want:      `{"message":"wrong json format", "path":"/deleteWord", "status":400}`,
			setupMock: func() {},
		},
		{
			Name: "Invalid user_id",
			Args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/deleteWord",
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
			Want:      `{"message":"Bad Request", "path":"/deleteWord", "status":400}`,
			setupMock: func() {},
		},
		{
			Name: "Invalid word",
			Args: args{
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
					ctx := inCtx(r.Context(), "user_id", "12345")
					return r.WithContext(ctx)
				}(),
			},
			Want:      `{"message":"Bad Request", "path":"/deleteWord", "status":400}`,
			setupMock: func() {},
		},
		{
			Name: "Invalid collecton name",
			Args: args{
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
					ctx := inCtx(r.Context(), "user_id", "")
					return r.WithContext(ctx)
				}(),
			},
			Want:      `{"message":"Bad Request", "path":"/deleteWord", "status":400}`,
			setupMock: func() {},
		},
		{
			Name: "Without user_id in ctx error",
			Args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/deleteWord",
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
			Want:      `{"message":"Unauthorized", "path":"/deleteWord", "status":401}`,
			setupMock: func() {},
		},
		{
			Name: "Internal error",
			Args: args{
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
					ctx := inCtx(r.Context(), "user_id", "12345")
					return r.WithContext(ctx)
				}(),
			},
			Want: `{"message":"Internal Server Error", "path":"/deleteWord", "status":500}`,
			setupMock: func() {
				s.srv.On("DeleteWord", mock.Anything, mock.Anything).Once().Return(
					errors.New("some internal error"),
				)
			},
		},
		{
			Name: "Valid request",
			Args: args{
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
					ctx := inCtx(r.Context(), "user_id", "12345")
					return r.WithContext(ctx)
				}(),
			},
			Want: `{"message":"success", "path":"/deleteWord", "status":200}`,
			setupMock: func() {
				s.srv.On("DeleteWord", mock.Anything, mock.Anything).Once().Return(nil)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.Name, func() {
			tc.setupMock()
			s.h.deleteWord(tc.Args.w, tc.Args.r)
			s.Assert().JSONEq(tc.Want, tc.Args.w.Body.String(),
				"Json response must be as expected")
		})
	}
}
