// Package googletransclient implements HTTP 2.0 connection to google.translate.com
package googletransclient

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/net/http2"
)

var queries = url.Values{
	"soc-platform": {"1"},
	"soc-device":   {"1"},
	"soc-app":      {"1"},
	"rpcids":       {"MkEWBc"},
	"rt":           {"c"},
	"bl":           {"boq_translate-webserver_20201207.13_p0"},
}

type TranlateClient struct {
	client *http.Client
	url    string
}

// Translates word from srcLang to trgtLang.
func (t *TranlateClient) Translate(text, srcLang, trgtLang string) ([]byte, error) {
	resp, err := t.client.PostForm(t.url, t.getPostForm(text, srcLang, trgtLang))
	if err != nil {
		return nil, fmt.Errorf("TranlateClient - Translate - PostForm: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("TranlateClient - Translate - ReadAll: %w", err)
	}

	return data, err
}

func (t *TranlateClient) getPostForm(text, srcLang, trgtLang string) url.Values {
	const RPCMethod = `[[["MkEWBc","[[\"%v\",\"%v\",\"%v\",true],[null]]",null,"generic"]]]`

	return url.Values{
		"f.req": {fmt.Sprintf(RPCMethod, text, srcLang, trgtLang)},
	}
}

func New(apiurl string) (*TranlateClient, error) {
	u, err := url.Parse(apiurl)
	if err != nil {
		return nil, err
	}

	u.RawQuery = queries.Encode()
	t := &http2.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: t,
	}

	transCLient := new(TranlateClient)
	transCLient.client = client
	transCLient.url = u.String()

	return transCLient, nil
}
