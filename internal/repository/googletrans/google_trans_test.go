package googletrans

import (
	"context"
	"testing"

	"github.com/Kin-dza-dzaa/flash_cards_api/config"
	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	"github.com/Kin-dza-dzaa/flash_cards_api/pkg/googletransclient"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func setupGoogleTrans(t *testing.T) *GoogleTranslate {
	t.Helper()
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Fatalf("setupGoogleTrans - config.ReadConfig: %v", err)
	}

	gc, err := googletransclient.New(cfg.GoogleAPI.URL)
	if err != nil {
		t.Fatalf("setupGoogleTrans - googletransclient.New: %v", err)
	}

	return New(gc, cfg.GoogleAPI.DefaultSrcLang, cfg.GoogleAPI.DefaultTrgtLang)
}

// Test makes real calls to google translate api.
func Test_Translate(t *testing.T) {
	tests := []struct {
		name    string
		word    string
		wantErr error
	}{
		{
			name:    "Unsupported word",
			word:    "bad_word!!!!!@#!@$#!%#",
			wantErr: entity.ErrWordNotSupported,
		},
		{
			name:    "Supported word",
			word:    "lead",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		ctx := context.Background()
		googletrans := setupGoogleTrans(t)

		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := googletrans.Translate(ctx, tt.word)
			if !cmp.Equal(gotErr, tt.wantErr, cmpopts.EquateErrors()) {
				t.Fatalf("wanted: %v but got: %v", tt.wantErr, gotErr)
			}
		})
	}
}
