package v1

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
)

func (s *wordHandler_Suite) Test_wordHandler_decodeCollection() {
	tests := []struct {
		Want    entity.Collection
		Req     *http.Request
		Name    string
		WantErr bool
	}{
		{
			Name: "Valid json",
			Req: httptest.NewRequest(http.MethodGet, "/decode",
				bytes.NewReader(
					[]byte(
						`
							{
								"word": "some_word",
								"collection_name": "some_coll"
							}
						`,
					),
				)),
			WantErr: false,
			Want: entity.Collection{
				Word: "some_word",
				Name: "some_coll",
			},
		},
		{
			Name: "Invalid json",
			Req: httptest.NewRequest(http.MethodGet, "/decode",
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
			WantErr: true,
			Want:    entity.Collection{},
		},
	}
	for _, tc := range tests {
		s.Run(tc.Name, func() {
			coll, err := s.h.decodeCollection(tc.Req)
			if tc.WantErr {
				s.Assert().Error(err, "Err must be not nill")
			} else {
				s.Assert().Nil(err, "Err must be nil")
			}
			s.Assert().Equal(tc.Want, coll, "Collcetions must be equal")
		})
	}
}

func (s *wordHandler_Suite) Test_wordHandler_encodeResopnse() {
	tests := []struct {
		ResRec      *httptest.ResponseRecorder
		ResponsHTTP *httpResponse
		Want        string
		Name        string
	}{
		{
			Name:   "Valid response",
			ResRec: httptest.NewRecorder(),
			ResponsHTTP: &httpResponse{
				Path:    "/some_path",
				Status:  200,
				Message: "test",
			},
			Want: `{"message":"test", "path":"/some_path", "status":200}`,
		},
	}
	for _, tc := range tests {
		s.Run(tc.Name, func() {
			s.h.encodeResponse(tc.ResRec, *tc.ResponsHTTP)
			s.Assert().JSONEq(tc.Want, tc.ResRec.Body.String(), "Json responses must be equal")
		})
	}
}
