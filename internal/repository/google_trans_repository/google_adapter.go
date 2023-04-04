// Package googletransrepo represents adapter layer for google.translate.com.
package googletransrepository

import (
	"bytes"
	"fmt"

	"github.com/Kin-dza-dzaa/flash_cards_api/internal/entity"
	googletransclient "github.com/Kin-dza-dzaa/flash_cards_api/pkg/google_trans_client"
	"github.com/tidwall/gjson"
)

type GoogleTranslate struct {
	client          *googletransclient.TranlateClient
	defaultSrcLang  string
	defaultTrgtLang string
}

func (t *GoogleTranslate) Translate(word string) (entity.WordTrans, error) {
	response, err := t.client.Translate(word, t.defaultSrcLang, t.defaultTrgtLang)
	if err != nil {
		return entity.WordTrans{}, fmt.Errorf("GoogleTranslate - Translate - client.Translate: %w", err)
	}
	wordTrans := t.unmarshal(response)
	if len(wordTrans.Translations) == 0 {
		return entity.WordTrans{}, entity.ErrWordNotSupported
	}

	return wordTrans, nil
}

func (t *GoogleTranslate) getValidJSON(data []byte) []byte {
	const validJSONPartIndex = 3
	return bytes.Split(data, []byte{'\n'})[validJSONPartIndex]
}

func (t *GoogleTranslate) getTrans(data []byte) gjson.Result {
	data = t.getValidJSON(data)
	const transPath = "0.2"
	return gjson.Parse(gjson.GetBytes(data, transPath).String())
}

func (t *GoogleTranslate) getWord(wordTransJRes gjson.Result) string {
	const (
		wordPath = "1.4.0"
	)

	return wordTransJRes.Get(wordPath).String()
}

func (t *GoogleTranslate) getLangs(wordTransJRes gjson.Result) (srcLang, trgtLang string) {
	const (
		srcLangPath  = "1.3"
		trgtLangPath = "1.1"
	)
	srcLang = wordTransJRes.Get(srcLangPath).String()
	trgtLang = wordTransJRes.Get(trgtLangPath).String()
	return srcLang, trgtLang
}

func (t *GoogleTranslate) getMainTranslation(wordTransJRes gjson.Result) string {
	const (
		mainTransPath = "1.0.0.5.0.0"
	)
	return wordTransJRes.Get(mainTransPath).String()
}

func (t *GoogleTranslate) getExamples(wordTransJRes gjson.Result) []string {
	const (
		examplesPath = "3.2.0"
		examplePath  = "1"
	)
	examplesJRes := wordTransJRes.Get(examplesPath).Array()
	examples := make([]string, 0, len(examplesJRes))

	for _, exampleJRes := range examplesJRes {
		example := exampleJRes.Get(examplePath).String()
		examples = append(examples, example)
	}

	return examples
}

func (t *GoogleTranslate) getPOSDefs(defsJRes []gjson.Result) []entity.WordDefinition {
	const (
		defPath     = "0"
		examplePath = "1"
	)
	defsWithExamples := make([]entity.WordDefinition, 0, len(defsJRes))

	for _, defJRes := range defsJRes {
		def := defJRes.Get(defPath).String()
		example := defJRes.Get(examplePath).String()

		defsWithExamples = append(
			defsWithExamples,
			entity.WordDefinition{Definition: def, Example: example},
		)
	}

	return defsWithExamples
}

func (t *GoogleTranslate) getDefs(
	wordTransJRes gjson.Result,
) map[entity.PartOfSpeech][]entity.WordDefinition {
	const (
		defsPath = "3.1.0"
		POSPath  = "0"
	)
	POSDefsJRes := wordTransJRes.Get(defsPath).Array()
	defs := make(map[entity.PartOfSpeech][]entity.WordDefinition, len(POSDefsJRes))

	for _, defsJRes := range POSDefsJRes {
		POS := entity.PartOfSpeech(defsJRes.Get(POSPath).String())
		defs[POS] = t.getPOSDefs(defsJRes.Get("1").Array())
	}

	return defs
}

func (t *GoogleTranslate) getPOSTransltions(transJRes []gjson.Result) []string {
	const (
		transPath = "0"
	)
	trans := make([]string, 0, len(transJRes))

	for _, tJres := range transJRes {
		t := tJres.Get(transPath).String()
		trans = append(trans, t)
	}

	return trans
}

func (t *GoogleTranslate) getTranslations(
	wordTransJRes gjson.Result,
) map[entity.PartOfSpeech][]string {
	const (
		POSPath          = "0"
		translationsPath = "3.5.0"
	)
	POStransJRes := wordTransJRes.Get(translationsPath).Array()
	trans := make(map[entity.PartOfSpeech][]string, len(POStransJRes))

	for _, transJRes := range POStransJRes {
		POS := entity.PartOfSpeech(transJRes.Get(POSPath).String())
		trans[POS] = t.getPOSTransltions(transJRes.Get("1").Array())
	}

	return trans
}

func (t *GoogleTranslate) unmarshal(data []byte) entity.WordTrans {
	wordTransJRes := t.getTrans(data)
	var wordTrans entity.WordTrans

	wordTrans.Translations = t.getTranslations(wordTransJRes)
	wordTrans.Word = t.getWord(wordTransJRes)
	wordTrans.SrcLang, wordTrans.TrgtLang = t.getLangs(wordTransJRes)
	wordTrans.MainTranslation = t.getMainTranslation(wordTransJRes)
	wordTrans.Examples = t.getExamples(wordTransJRes)
	wordTrans.Definitions = t.getDefs(wordTransJRes)

	return wordTrans
}

func New(client *googletransclient.TranlateClient, defaultSrcLang, defaultTrgtLang string) *GoogleTranslate {
	return &GoogleTranslate{
		client:          client,
		defaultSrcLang:  defaultSrcLang,
		defaultTrgtLang: defaultTrgtLang,
	}
}
