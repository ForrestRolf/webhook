package hook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/clbanning/mxj"
	"net/http"
	"net/url"
	"unicode"
)

type Request struct {
	ID string

	ContentType string

	Body []byte

	Headers map[string]interface{}

	Query map[string]interface{}

	Payload map[string]interface{}

	RawRequest *http.Request
}

type DebugRequest struct {
	Body    string              `json:"body,omitempty" bson:"body"`
	Headers map[string][]string `json:"headers,omitempty" bson:"headers"`
	Query   map[string][]string `json:"query,omitempty" bson:"query"`
}

func (r *Request) ParseJSONPayload() error {
	decoder := json.NewDecoder(bytes.NewReader(r.Body))
	decoder.UseNumber()

	var firstChar byte
	for i := 0; i < len(r.Body); i++ {
		if unicode.IsSpace(rune(r.Body[i])) {
			continue
		}
		firstChar = r.Body[i]
		break
	}

	if firstChar == byte('[') {
		var arrayPayload interface{}
		err := decoder.Decode(&arrayPayload)
		if err != nil {
			return fmt.Errorf("error parsing JSON array payload %+v", err)
		}

		r.Payload = make(map[string]interface{}, 1)
		r.Payload["root"] = arrayPayload
	} else {
		err := decoder.Decode(&r.Payload)
		if err != nil {
			return fmt.Errorf("error parsing JSON payload %+v", err)
		}
	}

	return nil
}

func (r *Request) ParseHeaders(headers map[string][]string) {
	r.Headers = make(map[string]interface{}, len(headers))

	for k, v := range headers {
		if len(v) > 0 {
			r.Headers[k] = v[0]
		}
	}
}

func (r *Request) ParseQuery(query map[string][]string) {
	r.Query = make(map[string]interface{}, len(query))

	for k, v := range query {
		if len(v) > 0 {
			r.Query[k] = v[0]
		}
	}
}

func (r *Request) ParseFormPayload() error {
	fd, err := url.ParseQuery(string(r.Body))
	if err != nil {
		return fmt.Errorf("error parsing form payload %+v", err)
	}

	r.Payload = make(map[string]interface{}, len(fd))

	for k, v := range fd {
		if len(v) > 0 {
			r.Payload[k] = v[0]
		}
	}

	return nil
}

func (r *Request) ParseXMLPayload() error {
	var err error

	r.Payload, err = mxj.NewMapXmlReader(bytes.NewReader(r.Body))
	if err != nil {
		return fmt.Errorf("error parsing XML payload: %+v", err)
	}

	return nil
}

func (r *DebugRequest) ToJson() string {
	str, _ := json.Marshal(r)
	return string(str)
}
