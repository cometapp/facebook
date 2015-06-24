// A facebook graph api client in go.
// https://github.com/huandu/facebook/
//
// Copyright 2012 - 2015, Huan Du
// Licensed under the MIT license
// https://github.com/huandu/facebook/blob/master/LICENSE

package facebook

import (
	"io"
	"net/http"
	"reflect"
)

type AppService interface {
	AppAccessToken() string
	ParseSignedRequest(signedRequest string) (res Result, err error)
	ParseCode(code string) (token string, err error)
	ParseCodeInfo(code, machineId string) (token string, expires int, newMachineId string, err error)
	ExchangeToken(accessToken string) (token string, expires int, err error)
	GetCode(accessToken string) (code string, err error)
	Session(accessToken string) *SessionService
	SessionFromSignedRequest(signedRequest string) (session *Session, err error)
}

type SessionService interface {
	Api(path string, method Method, params Params) (ResultService, error)
	Get(path string, params Params) (ResultService, error)
	Post(path string, params Params) (ResultService, error)
	Delete(path string, params Params) (ResultService, error)
	Put(path string, params Params) (ResultService, error)
	BatchApi(params ...Params) ([]Result, error)
	Batch(batchParams Params, params ...Params) ([]Result, error)
	FQL(query string) ([]Result, error)
	MultiFQL(queries Params) (Result, error)
	Request(request *http.Request) (res Result, err error)
	User() (id string, err error)
	Validate() (err error)
	Inspect(token string) (result Result, err error)
	AccessToken() string
	SetAccessToken(token string)
	SetVersion(version string)
	AppsecretProof() string
	EnableAppsecretProof(enabled bool) error
	App() *App
	Debug() DebugMode
	SetDebug(debug DebugMode) DebugMode
	// graph(path string, method Method, params Params) (res ResultService, err error)
	// graphBatch(batchParams Params, params ...Params) ([]Result, error)
	// graphFQL(params Params) (res Result, err error)
	// prepareParams(params Params)
	// sendGetRequest(uri string, res interface{}) (*http.Response, error)
	// sendPostRequest(uri string, params Params, res interface{}) (*http.Response, error)
	// sendOauthRequest(uri string, params Params) (Result, error)
	// sendRequest(request *http.Request) (response *http.Response, data []byte, err error)
	// isVideoPost(path string, method Method) bool
	// getUrl(name, path string, params Params) string
	// addDebugInfo(res ResultService, response *http.Response) ResultService
}

type ParamsService interface {
	Encode(writer io.Writer) (mime string, err error)
	encodeFormUrlEncoded(writer io.Writer) (mime string, err error)
	encodeMultipartForm(writer io.Writer) (mime string, err error)
}

type ResultService interface {
	Get(field string) interface{}
	Set(field string, value interface{})
	GetField(fields ...string) interface{}
	get(fields []string) interface{}
	Decode(v interface{}) (err error)
	DecodeField(field string, v interface{}) error
	Err() error
	Paging(session *SessionService) (*PagingResult, error)
	Batch() (*BatchResult, error)
	DebugInfo() *DebugInfo
	decode(v reflect.Value, fullName string) error
}

// Holds facebook application information.
type App struct {
	// Facebook app id
	AppId string

	// Facebook app secret
	AppSecret string

	// Facebook app redirect URI in the app's configuration.
	RedirectUri string

	// Enable appsecret proof in every API call to facebook.
	// Facebook document: https://developers.facebook.com/docs/graph-api/securing-requests
	EnableAppsecretProof bool
}

// An interface to send http request.
// This interface is designed to be compatible with type `*http.Client`.
type HttpClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
	Get(url string) (resp *http.Response, err error)
	Post(url string, bodyType string, body io.Reader) (resp *http.Response, err error)
}

// Holds a facebook session with an access token.
// Session should be created by App.Session or App.SessionFromSignedRequest.
type Session struct {
	HttpClient HttpClient
	Version    string // facebook versioning.

	accessToken string // facebook access token. can be empty.
	app         *App
	id          string

	enableAppsecretProof bool   // add "appsecret_proof" parameter in every facebook API call.
	appsecretProof       string // pre-calculated "appsecret_proof" value.

	debug DebugMode // using facebook debugging api in every request.
}

// API HTTP method.
// Can be GET, POST or DELETE.
type Method string

// Graph API debug mode.
// See https://developers.facebook.com/docs/graph-api/using-graph-api/v2.3#graphapidebugmode
type DebugMode string

// API params.
//
// For general uses, just use Params as a ordinary map.
//
// For advanced uses, use MakeParams to create Params from any struct.
type Params map[string]interface{}

// Facebook API call result.
type Result map[string]interface{}

// Represents facebook API call result with paging information.
type PagingResult struct {
	session  *SessionService
	paging   pagingData
	previous string
	next     string
}

// Represents facebook batch API call result.
// See https://developers.facebook.com/docs/graph-api/making-multiple-requests/#multiple_methods.
type BatchResult struct {
	StatusCode int         // HTTP status code.
	Header     http.Header // HTTP response headers.
	Body       string      // Raw HTTP response body string.
	Result     Result      // Facebook api result parsed from body.
}

// Facebook API error.
type Error struct {
	Message      string
	Type         string
	Code         int
	ErrorSubcode int // subcode for authentication related errors.
}

// Binary data.
type binaryData struct {
	Filename string    // filename used in multipart form writer.
	Source   io.Reader // file data source.
}

// Binary file.
type binaryFile struct {
	Filename string // filename used in multipart form writer.
	Path     string // path to file. must be readable.
}

// DebugInfo is the debug information returned by facebook when debug mode is enabled.
type DebugInfo struct {
	Messages []DebugMessage // debug messages. it can be nil if there is no message.
	Header   http.Header    // all HTTP headers for this response.
	Proto    string         // HTTP protocol name for this response.

	// Facebook debug HTTP headers.
	FacebookApiVersion string // the actual graph API version provided by facebook-api-version HTTP header.
	FacebookDebug      string // the X-FB-Debug HTTP header.
	FacebookRev        string // the x-fb-rev HTTP header.
}

// DebugMessage is one debug message in "__debug__" of graph API response.
type DebugMessage struct {
	Type    string
	Message string
	Link    string
}
