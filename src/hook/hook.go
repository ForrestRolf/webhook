package hook

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"net/textproto"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	EnvNamespace string = "WEBHOOK_"
)
const (
	SourceHeader         string = "header"
	SourceQuery          string = "url"
	SourceQueryAlias     string = "query"
	SourcePayload        string = "payload"
	SourceRawRequestBody string = "raw-request-body"
	SourceRequest        string = "request"
	SourceString         string = "string"
	SourceEntirePayload  string = "entire-payload"
	SourceEntireQuery    string = "entire-query"
	SourceEntireHeaders  string = "entire-headers"
)
const (
	ActionShellDriver string = "shell"
	ActionHttpDriver  string = "http"
)

type ParameterNodeError struct {
	key string
}

func (e *ParameterNodeError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("parameter node not found: %s", e.key)
}

func IsParameterNodeError(err error) bool {
	switch err.(type) {
	case *ParameterNodeError:
		return true
	default:
		return false
	}
}

type ArgumentError struct {
	Argument Argument
}

func (e *ArgumentError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("couldn't retrieve argument for %+v", e.Argument)
}

type SourceError struct {
	Argument Argument
}

func (e *SourceError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("invalid source for argument %+v", e.Argument)
}

type ParseError struct {
	Err error
}

func (e *ParseError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return e.Err.Error()
}

func GetParameter(s string, params interface{}) (interface{}, error) {
	if params == nil {
		return nil, errors.New("no parameters")
	}

	paramsValue := reflect.ValueOf(params)

	switch paramsValue.Kind() {
	case reflect.Slice:
		paramsValueSliceLength := paramsValue.Len()
		if paramsValueSliceLength > 0 {

			if p := strings.SplitN(s, ".", 2); len(p) > 1 {
				index, err := strconv.ParseUint(p[0], 10, 64)

				if err != nil || paramsValueSliceLength <= int(index) {
					return nil, &ParameterNodeError{s}
				}

				return GetParameter(p[1], params.([]interface{})[index])
			}

			index, err := strconv.ParseUint(s, 10, 64)

			if err != nil || paramsValueSliceLength <= int(index) {
				return nil, &ParameterNodeError{s}
			}

			return params.([]interface{})[index], nil
		}

		return nil, &ParameterNodeError{s}

	case reflect.Map:
		if v, ok := params.(map[string]interface{})[s]; ok {
			return v, nil
		}

		p := strings.SplitN(s, ".", 2)
		if pValue, ok := params.(map[string]interface{})[p[0]]; ok {
			if len(p) > 1 {
				return GetParameter(p[1], pValue)
			}

			return pValue, nil
		}
	}

	return nil, &ParameterNodeError{s}
}

func ExtractParameterAsString(s string, params interface{}) (string, error) {
	pValue, err := GetParameter(s, params)
	if err != nil {
		return "", err
	}

	switch v := reflect.ValueOf(pValue); v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice:
		r, err := json.Marshal(pValue)
		if err != nil {
			return "", err
		}

		return string(r), nil

	default:
		return fmt.Sprintf("%v", pValue), nil
	}
}

type Argument struct {
	Source  string `json:"source,omitempty"`
	Name    string `json:"name,omitempty"`
	EnvName string `json:"envname,omitempty"`
}

func (arg *Argument) Get(r *Request) (string, error) {
	var source *map[string]interface{}
	key := arg.Name

	switch arg.Source {
	case SourceHeader:
		source = &r.Headers
		key = textproto.CanonicalMIMEHeaderKey(arg.Name)

	case SourceQuery, SourceQueryAlias:
		source = &r.Query

	case SourcePayload:
		source = &r.Payload

	case SourceString:
		return arg.Name, nil

	case SourceRawRequestBody:
		return string(r.Body), nil

	case SourceRequest:
		if r == nil || r.RawRequest == nil {
			return "", errors.New("request is nil")
		}

		switch strings.ToLower(arg.Name) {
		case "remote-addr":
			return r.RawRequest.RemoteAddr, nil
		case "method":
			return r.RawRequest.Method, nil
		default:
			return "", fmt.Errorf("unsupported request key: %q", arg.Name)
		}

	case SourceEntirePayload:
		res, err := json.Marshal(&r.Payload)
		if err != nil {
			return "", err
		}

		return string(res), nil

	case SourceEntireHeaders:
		res, err := json.Marshal(&r.Headers)
		if err != nil {
			return "", err
		}

		return string(res), nil

	case SourceEntireQuery:
		res, err := json.Marshal(&r.Query)
		if err != nil {
			return "", err
		}

		return string(res), nil
	}

	if source != nil {
		return ExtractParameterAsString(key, *source)
	}

	return "", errors.New("no source for value retrieval")
}

type Hook struct {
	ID                    string     `json:"id,omitempty"`
	Name                  string     `json:"name,omitempty"`
	TriggerRule           *Rules     `json:"trigger_rule,omitempty"`
	Actions               *[]Action  `json:"actions,omitempty"`
	PassArgumentsToAction []Argument `json:"pass_arguments_to_action,omitempty"`
}
type Action struct {
	Driver     string                 `json:"driver,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}
type ShellAction struct {
	WorkingDirectory string `bson:"workingDirectory" json:"workingDirectory"`
	Scripts          string `bson:"scripts" json:"scripts,omitempty"`
}
type HttpAction struct {
	Method       string `json:"method,omitempty" bson:"method"`
	Url          string `json:"url,omitempty" bson:"url"`
	ContentType  string `json:"contentType,omitempty" bson:"contentType"`
	Payload      string `json:"payload,omitempty" bson:"payload"`
	AuthToken    string `json:"authToken,omitempty" bson:"authToken"`
	Timeout      int    `json:"timeout,omitempty" bson:"timeout"`
	SaveResponse bool   `json:"saveResponse,omitempty" bson:"saveResponse"`
}

type Hooks []Hook

type Rules struct {
	And   *AndRule   `json:"and,omitempty"`
	Or    *OrRule    `json:"or,omitempty"`
	Match *MatchRule `json:"match,omitempty"`
}
type AndRule []Rules
type OrRule []Rules
type MatchRule struct {
	Type      string   `json:"type,omitempty"`
	Regex     string   `json:"regex,omitempty"`
	Secret    string   `json:"secret,omitempty"`
	Value     string   `json:"value,omitempty"`
	Parameter Argument `json:"parameter,omitempty"`
}

const (
	MatchValue string = "value"
	MatchRegex string = "regex"
)

func (r Rules) Evaluate(req *Request) (bool, error) {
	switch {
	case r.And != nil:
		return r.And.Evaluate(req)
	case r.Or != nil:
		return r.Or.Evaluate(req)
	case r.Match != nil:
		return r.Match.Evaluate(req)
	}
	return false, nil
}

func (r AndRule) Evaluate(req *Request) (bool, error) {
	res := true
	for _, v := range r {
		rv, err := v.Evaluate(req)
		if err != nil {
			return false, err
		}

		res = res && rv
		if !res {
			return res, nil
		}
	}
	return res, nil
}

func (r OrRule) Evaluate(req *Request) (bool, error) {
	res := false
	for _, v := range r {
		rv, err := v.Evaluate(req)
		if err != nil {
			return false, err
		}

		res = res || rv
		if res {
			return res, nil
		}
	}

	return res, nil
}

func (r MatchRule) Evaluate(req *Request) (bool, error) {
	arg, err := r.Parameter.Get(req)
	if err == nil {
		switch r.Type {
		case MatchValue:
			return compare(arg, r.Value), nil
		case MatchRegex:
			return regexp.MatchString(r.Regex, arg)
		}
	}
	return false, err
}

func (h *Hook) ExtractArgumentsForEnv(r *Request) ([]string, []error) {
	args := make([]string, 0)
	errors := make([]error, 0)
	for i := range h.PassArgumentsToAction {
		arg, err := h.PassArgumentsToAction[i].Get(r)
		if err != nil {
			errors = append(errors, &ArgumentError{h.PassArgumentsToAction[i]})
			continue
		}
		if h.PassArgumentsToAction[i].EnvName != "" {
			args = append(args, EnvNamespace+h.PassArgumentsToAction[i].EnvName+"="+arg)
		}
	}
	if len(errors) > 0 {
		return args, errors
	}
	return args, nil
}

func (h *Hook) ExtractArgumentsAsMap(r *Request) (map[string]string, []error) {
	args := make(map[string]string, 0)
	errors := make([]error, 0)
	for i := range h.PassArgumentsToAction {
		arg, err := h.PassArgumentsToAction[i].Get(r)
		if err != nil {
			errors = append(errors, &ArgumentError{h.PassArgumentsToAction[i]})
			continue
		}
		if h.PassArgumentsToAction[i].Name != "" {
			args[h.PassArgumentsToAction[i].Name] = arg
		}
	}
	if len(errors) > 0 {
		return args, errors
	}
	return args, nil
}

func compare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

func getenv(s string) string {
	return os.Getenv(s)
}
