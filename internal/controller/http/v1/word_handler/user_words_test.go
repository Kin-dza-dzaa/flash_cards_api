package wordhadnler

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	"github.com/stretchr/testify/mock"
)

func (s *wordHandler_Suite) Test_wordHandler_userWords() {
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
			Name: "Without user_id in ctx error",
			Args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/getWords", nil),
			},
			Want:      `{"message":"Unauthorized", "path":"/getWords", "status":401}`,
			setupMock: func() {},
		},
		{
			Name: "Internal error",
			Args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/getWords", nil)
					ctx := inCtx(r.Context(), "user_id", "12345")
					return r.WithContext(ctx)
				}(),
			},
			Want: `{"message":"Internal Server Error", "path":"/getWords", "status":500}`,
			setupMock: func() {
				s.srv.On("UserWords", mock.Anything, mock.Anything).Once().Return(
					nil, errors.New("some internal error"),
				)
			},
		},
		{
			Name: "Valid request",
			Args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/getWords", nil)
					ctx := inCtx(r.Context(), "user_id", "12345")
					return r.WithContext(ctx)
				}(),
			},
			Want: `{"path":"/getWords","status":200,"message":"success","user_words":{"words":{}}}`,
			setupMock: func() {
				s.srv.On("UserWords", mock.Anything, mock.Anything).Once().Return(
					&entity.UserWords{
						Words: make(map[entity.CollectionName][]entity.WordData),
					}, nil,
				)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.Name, func() {
			tc.setupMock()
			s.h.userWords(tc.Args.w, tc.Args.r)
			s.Assert().JSONEq(tc.Want, tc.Args.w.Body.String(),
				"Json response must be as expected")
		})
	}
}
