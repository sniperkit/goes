package config

import (
	// "net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sniperkit/iris"
)

type Config struct {
	App       *App       `json:"application" yaml:"application" toml:"application"`
	Api       *Api       `json:"api" yaml:"api" toml:"api"`
	Websocket *Websocket `json:"websocket" yaml:"websocket" toml:"websocket"`
	Store     *Database  `json:"database" yaml:"database" toml:"database"`
	Server    *Server    `json:"server" yaml:"server" toml:"server"`
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
	Port      int        `json:"port,omitempty" yaml:"port,omitempty" toml:"port,omitempty"`
	Path      string     `default:"/" json:"path,omitempty" yaml:"path,omitempty" toml:"path,omitempty"`
	Static    *Static    `json:"static,omitempty" yaml:"static,omitempty" toml:"static,omitempty"`
	Resources []Resource `json:"resources,omitempty" yaml:"resources,omitempty" toml:"resources,omitempty"`
	URLs      []URL      `json:"urls,omitempty" yaml:"urls,omitempty" toml:"urls,omitempty"`
	Delay     int        `json:"delay,omitempty" yaml:"delay,omitempty" toml:"delay,omitempty"`
	Auth      *Auth      `json:"auth,omitempty" yaml:"auth,omitempty" toml:"auth,omitempty"`
	JWT       *JWTData   `json:"jwt,omitempty" yaml:"jwt,omitempty" toml:"jwt,omitempty"`
	Cors      *Cors      `json:"cors,omitempty" yaml:"cors,omitempty" toml:"cors,omitempty"`
	Docs      *Docs      `json:"docs,omitempty" yaml:"docs,omitempty" toml:"docs,omitempty"`
	EnableLog bool       `json:"enable_log,omitempty" yaml:"enable_log,omitempty" toml:"enable_log,omitempty"`
}

type Docs struct {
	Enabled  bool              `json:"enable,omitempty" yaml:"enable,omitempty" toml:"enable,omitempty"`
	BaseUrls map[string]string `json:"urls,omitempty" yaml:"urls,omitempty" toml:"urls,omitempty"`
	DocTitle string            `default:"Api Documentation" json:"title,omitempty" yaml:"title,omitempty" toml:"title,omitempty"`
	DocPath  string            `default:"/docs" json:"path,omitempty" yaml:"path,omitempty" toml:"path,omitempty"`
	DocFile  string            `default:"apidoc.html" json:"filename,omitempty" yaml:"filename,omitempty" toml:"filename,omitempty"`
}

// Resource respresent a single resource in rest api.
type Resource struct {
	Name    string            `json:"name" yaml:"name" toml:"name"`
	Headers map[string]string `json:"headers,omitempty" yaml:"headers,omitempty" toml:"headers,omitempty"`
}

type Cors struct {
	AllowedOrigins     []string `json:"allowed_origins,omitempty" yaml:"allowed_origins,omitempty" toml:"allowed_origins,omitempty"`
	AllowedMethods     []string `json:"allowed_methods,omitempty" yaml:"allowed_methods,omitempty" toml:"allowed_methods,omitempty"`
	AllowedHeaders     []string `json:"allowed_headers,omitempty" yaml:"allowed_headers,omitempty" toml:"allowed_headers,omitempty"`
	AllowCredentials   bool     `default:"false" json:"allow_credentials,omitempty" yaml:"allow_credentials,omitempty" toml:"allow_credentials,omitempty"`
	OptionsPassthrough bool     `default:"false" json:"options_pass_through,omitempty" yaml:"options_pass_through,omitempty" toml:"options_pass_through,omitempty"`
	Debug              bool     `default:"false" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
	MaxAge             int      `default:"3600" json:"max_age,omitempty" yaml:"max_age,omitempty" toml:"max_age,omitempty"`
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
	URL     string                 `json:"url" yaml:"url" toml:"url"`
	Method  string                 `json:"method" yaml:"method" toml:"method"`
	Handler func(ctx iris.Context) `json:"-" yaml:"-" toml:"-"`
}

type URL struct {
	URL         string `json:"url" yaml:"url" toml:"url"`
	Method      string `json:"method,omitempty" yaml:"method,omitempty" toml:"method,omitempty"`
	ContentType string `json:"content_type,omitempty" yaml:"content_type,omitempty" toml:"content_type,omitempty"`
	File        string `json:"file,omitempty" yaml:"file,omitempty" toml:"file,omitempty"`
	Engine      string `default:"file" json:"engine,omitempty" yaml:"engine,omitempty" toml:"engine,omitempty"`
	// PrefixPath  string            `json:"file,omitempty" yaml:"file,omitempty" toml:"file,omitempty"`
	StatusCode int               `json:"status,omitempty" yaml:"status,omitempty" toml:"status,omitempty"`
	Headers    map[string]string `json:"headers,omitempty" yaml:"headers,omitempty" toml:"headers,omitempty"`
}

type Static struct {
	URL  string `json:"url" yaml:"url" toml:"url"`
	Path string `json:"path" yaml:"path" toml:"path"`
}

type Server struct {
	Env       string   `default:"dev" json:"env" yaml:"env" toml:"env"`
	SessionID string   `default:"snk" json:"session_id" yaml:"session_id" toml:"session_id"`
	Port      int      `default:"8080" json:"port" yaml:"port" toml:"port"`
	Engine    string   `default:"iris" json:"engine" yaml:"engine" toml:"engine"`
	Session   *Session `json:"session,omitempty" yaml:"session,omitempty" toml:"session,omitempty"`
	Settings  struct {
		Gin     *Gin     `json:"gin,omitempty" yaml:"gin,omitempty" toml:"gin,omitempty"`
		Iris    *Iris    `json:"iris,omitempty" yaml:"iris,omitempty" toml:"iris,omitempty"`
		Echo    *Echo    `json:"echo,omitempty" yaml:"echo,omitempty" toml:"echo,omitempty"`
		Gorilla *Gorilla `json:"gorilla,omitempty" yaml:"gorilla,omitempty" toml:"gorilla,omitempty"`
	} `json:"settings,omitempty" yaml:"settings,omitempty" toml:"settings,omitempty"`
	Debug bool `default:"debug" json:"debug" yaml:"debug" toml:"debug"`
}

type Session struct {
	ID                          string        `default:"snk" json:"id,omitempty" yaml:"id,omitempty" toml:"id,omitempty"`
	Expires                     time.Duration `default:"3600" json:"expires,omitempty" yaml:"expires,omitempty" toml:"expires,omitempty"` // <=0 means unlimited life. Defaults to 0.
	AllowReclaim                bool          `default:"false" json:"allow_reclaim,omitempty" yaml:"allow_reclaim,omitempty" toml:"allow_reclaim,omitempty"`
	CookieSecureTLS             bool          `default:"false" json:"cookie_secure_tls,omitempty" yaml:"cookie_secure_tls,omitempty" toml:"cookie_secure_tls,omitempty"`
	DisableSubdomainPersistence bool          `default:"false" json:"disable_subdomain_persistence,omitempty" yaml:"disable_subdomain_persistence,omitempty" toml:"disable_subdomain_persistence,omitempty"`
	DataDir                     string        `default:"./shared/data/sessions" json:"data_dir,omitempty" yaml:"data_dir,omitempty" toml:"data_dir,omitempty"`
}

type Gin struct {
	TimeFormat string `default:"Mon, 02 Jan 2006 15:04:05 GMT" json:"time_format,omitempty" yaml:"time_format,omitempty" toml:"time_format,omitempty"`
	Charset    string `default:"UTF-8" json:"charset,omitempty" yaml:"charset,omitempty" toml:"charset,omitempty"`
}

type Echo struct {
	TimeFormat string `default:"Mon, 02 Jan 2006 15:04:05 GMT" json:"time_format,omitempty" yaml:"time_format,omitempty" toml:"time_format,omitempty"`
	Charset    string `default:"UTF-8" json:"charset,omitempty" yaml:"charset,omitempty" toml:"charset,omitempty"`
}

type Gorilla struct {
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
	Port              int            `default:"9001" json:"port" yaml:"port" toml:"port"`
	EnableLog         bool           `json:"enable_log,omitempty" yaml:"enable_log,omitempty" toml:"enable_log,omitempty"`
	Delay             int            `json:"delay,omitempty" yaml:"delay,omitempty" toml:"delay,omitempty"`
	Auth              Auth           `json:"auth,omitempty" yaml:"auth,omitempty" toml:"auth,omitempty"`
	Path              string         `default:"/ws" json:"path,omitempty" yaml:"path,omitempty" toml:"path,omitempty"`
	ReadBufferSize    int            `default:"1024" json:"read_buffer_size,omitempty" yaml:"read_buffer_size,omitempty" toml:"read_buffer_size,omitempty"`
	WriteBufferSize   int            `default:"1024" json:"write_buffer_size,omitempty" yaml:"write_buffer_size,omitempty" toml:"write_buffer_size,omitempty"`
	BinaryMessages    bool           `default:"false" json:"allow_binary_messages,omitempty" yaml:"allow_binary_messages,omitempty" toml:"allow_binary_messages,omitempty"`
	EnableCompression bool           `default:"false" json:"enable_compression,omitempty" yaml:"enable_compression,omitempty" toml:"enable_compression,omitempty"`
	Subprotocols      []Subprotocols `json:"sub_protocols,omitempty" yaml:"sub_protocols,omitempty" toml:"sub_protocols,omitempty"`
	HandshakeTimeout  time.Duration  `default:"10" json:"handshake_timeout,omitempty" yaml:"handshake_timeout,omitempty" toml:"handshake_timeout,omitempty"`
	WriteTimeout      time.Duration  `default:"10" json:"write_timeout,omitempty" yaml:"write_timeout,omitempty" toml:"write_timeout,omitempty"`
	ReadTimeout       time.Duration  `default:"10" json:"read_timeout,omitempty" yaml:"read_timeout,omitempty" toml:"read_timeout,omitempty"`
	PongTimeout       time.Duration  `default:"10" json:"pong_period,omitempty" yaml:"pong_period,omitempty" toml:"pong_period,omitempty"`
	PingPeriod        time.Duration  `default:"10" json:"ping_period,omitempty" yaml:"ping_period,omitempty" toml:"ping_period,omitempty"`
}

type Subprotocols struct {
	Name string `json:"name" yaml:"name" toml:"name"`
}

type Auth struct {
	Name     string `json:"username" yaml:"username" toml:"username"`
	Password string `json:"password" yaml:"password" toml:"password"`
}
