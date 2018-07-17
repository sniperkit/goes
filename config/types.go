package config

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type Config struct {
	App       App       `json:"app" yaml:"app" toml:"app"`
	Api       Api       `json:"api" yaml:"rest" toml:"rest"`
	Websocket Websocket `json:"websocket" yaml:"websocket" toml:"websocket"`
	Store     Database  `json:"database" yaml:"database" toml:"database"`
	Server    Server    `json:"server" yaml:"server" toml:"server"`
}

type App struct {
	Name      string `default:"Your app" json:"name" yaml:"name" toml:"name"`
	Verbose   bool   `default:"false" json:"verbose" yaml:"verbose" toml:"verbose"`
	EnableLog bool   `default:"true" json:"enable_log" yaml:"enable_log" toml:"enable_log"`
	Dir       struct {
		Store  string `default:"./shared/data/storage" json:"storage" yaml:"storage" toml:"storage"`
		Cache  string `default:"./shared/data/cache" json:"cache" yaml:"cache" toml:"cache"`
		Export string `default:"./shared/export" json:"export" yaml:"export" toml:"export"`
	} `json:"data" yaml:"data" toml:"data"`
	Limits     Limits `json:"limits" yaml:"limits" toml:"limits"`
	Web        Web    `json:"web" yaml:"web" toml:"web"`
	Env        string `default:"dev" json:"env" yaml:"env" toml:"env"`
	EnvFile    string `default:".env" json:"env_file" yaml:"env_file" toml:"env_file"`
	UseProxy   bool   `json:"use_proxy" yaml:"use_proxy" toml:"use_proxy"`
	ProxyUri   string `json:"proxy_uri" yaml:"proxy_uri" toml:"proxy_uri"`
	Port       uint   `json:"port" yaml:"port" toml:"port"`
	StaticPort uint   `json:"static_port" yaml:"static_port" toml:"static_port"`
}

type Limits struct {
	PageSize            int `default:"20" json:"page_size" yaml:"page_size" toml:"page_size"`
	MaxPageSize         int `default:"200" json:"max_page_size" yaml:"max_page_size" toml:"max_page_size"`
	MinPageSize         int `default:"10" json:"min_page_size" yaml:"min_page_size" toml:"min_page_size"`
	MinOrder            int `default:"0" json:"min_order" yaml:"min_order" toml:"min_order"`
	MaxOrder            int `default:"10000" json:"max_order" yaml:"max_order" toml:"max_order"`
	MaxNameLength       int `default:"100" json:"max_name_length" yaml:"max_name_length" toml:"max_name_length"`
	MaxContentLength    int `default:"10000" json:"max_content_length" yaml:"max_content_length" toml:"max_content_length"`
	MaxArticleCateCount int `default:"6" json:"max_article_cate_count" yaml:"max_article_cate_count" toml:"max_article_cate_count"`
	MaxCommentLength    int `default:"5000" json:"max_comment_length" yaml:"max_comment_length" toml:"max_comment_length"`
}

type Web struct {
	Static         StaticWeb `json:"static" yaml:"static" toml:"static"`
	Title          string    `json:"title" yaml:"title" toml:"title"`
	JavascriptPath string    `json:"javascript_path" yaml:"javascript_path" toml:"javascript_path"`
	ImagePath      string    `json:"image_path" yaml:"image_path" toml:"image_path"`
	CssPath        string    `json:"css_path" yaml:"css_path" toml:"css_path"`
}

type StaticWeb struct {
	PrefixPath string `json:"prefix_path" yaml:"prefix_path" toml:"prefix_path"`
}

type Database struct {
	Dialect      string `json:"dialect" yaml:"dialect" toml:"dialect"`
	Database     string `json:"database" yaml:"database" toml:"database"`
	User         string `json:"user" yaml:"user" toml:"user"`
	Password     string `json:"password" yaml:"password" toml:"password"`
	Charset      string `json:"charset" yaml:"charset" toml:"charset"`
	Host         string `json:"host" yaml:"host" toml:"host"`
	Port         int    `json:"port" yaml:"port" toml:"port"`
	DSN          string `json:"dsn" yaml:"dsn" toml:"dsn"`
	MaxIdleConns int    `json:"max_idle_conns" yaml:"max_idle_conns" toml:"max_idle_conns"`
	MaxOpenConns int    `json:"max_open_conns" yaml:"max_open_conns" toml:"max_open_conns"`
}

// Config type represent configuration json.
type Api struct {
	Port      int        `json:"port" yaml:"port" toml:"port"`
	Delay     int        `json:"delay" yaml:"delay" toml:"delay"`
	Auth      *Auth      `json:"auth" yaml:"auth" toml:"auth"`
	JWT       *JWTData   `json:"jwt" yaml:"jwt" toml:"jwt"`
	Static    *Static    `json:"static" yaml:"static" toml:"static"`
	Resources []Resource `json:"resources" yaml:"resources" toml:"resources"`
	URLs      []URL      `json:"urls" yaml:"urls" toml:"urls"`
	EnableLog bool       `json:"enable_log" yaml:"enable_log" toml:"enable_log"`
	Path      string     `json:"-" yaml:"-" toml:"-"`
}

// Resource respresent a single resource in rest api.
type Resource struct {
	Name    string            `json:"name" yaml:"name" toml:"name"`
	Headers map[string]string `json:"headers,omitempty" yaml:"headers,omitempty" toml:"headers,omitempty"`
}

// JWTData represents jwt token details.
type JWTData struct {
	URL    string        `json:"url" yaml:"url" toml:"url"`
	EXP    int           `json:"exp" yaml:"exp" toml:"exp"`
	Secret string        `json:"secret" yaml:"secret" toml:"secret"`
	Data   jwt.MapClaims `json:"data" yaml:"data" toml:"data"`
}

// Endpoint defines a single rest endpoint(route) in the server
type Endpoint struct {
	URL     string       `json:"url" yaml:"url" toml:"url"`
	Method  string       `json:"method" yaml:"method" toml:"method"`
	Handler http.Handler `json:"-" yaml:"-" toml:"-"`
}

type URL struct {
	URL         string            `json:"url" yaml:"url" toml:"url"`
	Method      string            `json:"method,omitempty" yaml:"method,omitempty" toml:"method,omitempty"`
	ContentType string            `json:"content_type,omitempty" yaml:"content_type,omitempty" toml:"content_type,omitempty"`
	File        string            `json:"file,omitempty" yaml:"file,omitempty" toml:"file,omitempty"`
	StatusCode  int               `json:"status,omitempty" yaml:"status,omitempty" toml:"status,omitempty"`
	Headers     map[string]string `json:"headers,omitempty" yaml:"headers,omitempty" toml:"headers,omitempty"`
}

type Static struct {
	URL  string `json:"url" yaml:"url" toml:"url"`
	Path string `json:"path" yaml:"path" toml:"path"`
}

type Server struct {
	Env       string `default:"dev" json:"env" yaml:"env" toml:"env"`
	SessionID string `default:"snk" json:"session_id" yaml:"session_id" toml:"session_id"`
	Port      int    `default:"8080" json:"port" yaml:"port" toml:"port"`
	Engine    string `default:"iris" json:"engine" yaml:"engine" toml:"engine"`
	Settings  struct {
		Iris *Iris `json:"iris,omitempty" yaml:"iris,omitempty" toml:"iris,omitempty"`
		Gin  *Gin  `json:"gin,omitempty" yaml:"gin,omitempty" toml:"gin,omitempty"`
	} `json:"settings,omitempty" yaml:"settings,omitempty" toml:"settings,omitempty"`
}

type Gin struct {
	TimeFormat string `default:"Mon, 02 Jan 2006 15:04:05 GMT" json:"time_format,omitempty" yaml:"time_format,omitempty" toml:"time_format,omitempty"`
	Charset    string `default:"UTF-8" json:"charset,omitempty" yaml:"charset,omitempty" toml:"charset,omitempty"`
}

type Iris struct {
	DisableStartupLog                 bool   `json:"disable_startup_log,omitempty" yaml:"disable_startup_log,omitempty" toml:"disable_startup_log,omitempty"`
	DisableInterruptHandler           bool   `json:"disable_interrupt_handler,omitempty" yaml:"disable_interrupt_handler,omitempty" toml:"disable_interrupt_handler,omitempty"`
	DisablePathCorrection             bool   `json:"disable_path_correction,omitempty" yaml:"disable_path_correction,omitempty" toml:"disable_path_correction,omitempty"`
	EnablePathEscape                  bool   `json:"enable_path_escape" yaml:"enable_path_escape" toml:"enable_path_escape"`
	FireMethodNotAllowed              bool   `json:"fire_method_not_allowed" yaml:"fire_method_not_allowed" toml:"fire_method_not_allowed"`
	DisableBodyConsumptionOnUnmarshal bool   `json:"disable_body_consumption_on_unmarshal" yaml:"disable_body_consumption_on_unmarshal" toml:"disable_body_consumption_on_unmarshal"`
	DisableAutoFireStatusCode         bool   `json:"disable_auto_fire_status_code" yaml:"disable_auto_fire_status_code" toml:"disable_auto_fire_status_code"`
	TimeFormat                        string `default:"Mon, 02 Jan 2006 15:04:05 GMT" json:"time_format,omitempty" yaml:"time_format,omitempty" toml:"time_format,omitempty"`
	Charset                           string `default:"UTF-8" json:"charset,omitempty" yaml:"charset,omitempty" toml:"charset,omitempty"`
}

type Websocket struct {
	Port  int  `json:"port" yaml:"port" toml:"port"`
	Delay int  `json:"delay" yaml:"delay" toml:"delay"`
	Auth  Auth `json:"auth" yaml:"auth" toml:"auth"`
	// JWT       *rest.JWTData   `json:"jwt" yaml:"jwt" toml:"jwt"`
	// Static    *rest.Static    `json:"static" yaml:"static" toml:"static"`
	Path string `json:"-" yaml:"-" toml:"-"`
	// Resources []rest.Resource `json:"resources" yaml:"resources" toml:"resources"`
	// URLs      []rest.URL      `json:"urls" yaml:"urls" toml:"urls"`
	EnableLog bool `json:"enable_log" yaml:"enable_log" toml:"enable_log"`
}

type Auth struct {
	Name     string `json:"username" yaml:"username" toml:"username"`
	Password string `json:"password" yaml:"password" toml:"password"`
}
