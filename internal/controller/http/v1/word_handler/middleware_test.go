package wordhadnler

import (
	"net/http"
	"net/http/httptest"
)

func (s *wordHandler_Suite) Test_wordHandler_jwtAuthenticator() {
	dmyHand := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}

	tests := []struct {
		Name string
		Args args
		Want string
	}{
		{
			Name: "Without jwt in header",
			Want: `{"message":"Unauthorized", "path":"/jwt", "status":401}`,
			Args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/jwt", nil),
			},
		},
		{
			Name: "With invalid jwt in header",
			Want: `{"message":"Unauthorized", "path":"/jwt", "status":401}`,
			Args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/jwt", nil)
					r.Header.Add("Authorization", "BEARER invalid.jwt.bad_jwt")
					return r
				}(),
			},
		},
	}

	for _, tc := range tests {
		s.h.jwtAuthenticator(dmyHand).ServeHTTP(tc.Args.w, tc.Args.r)
		s.Assert().JSONEq(tc.Want, tc.Args.w.Body.String(), "Json response must be as expected")
	}
}
