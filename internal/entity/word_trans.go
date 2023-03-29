// Package entity describes business logic entity.
package entity

type (
	PartOfSpeech string

	WordDefinition struct {
		Definition string `json:"definition"`
		Example    string `json:"example,omitempty"`
	}

	WordTrans struct {
		Word            string                            `json:"word"`
		SrcLang         string                            `json:"source_language"`
		TrgtLang        string                            `json:"target_language"`
		Examples        []string                          `json:"examples,omitempty"`
		Definitions     map[PartOfSpeech][]WordDefinition `json:"definitions_with_examples,omitempty"`
		Translations    map[PartOfSpeech][]string         `json:"transltions"`
		MainTranslation string                            `json:"main_translation"`
	}
)
