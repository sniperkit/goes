package config

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
	Web        Web    `json:"web" yaml:"web" toml:"web"`
	Env        string `default:"dev" json:"env" yaml:"env" toml:"env"`
	EnvFile    string `default:".env" json:"env_file" yaml:"env_file" toml:"env_file"`
	UseProxy   bool   `json:"use_proxy" yaml:"use_proxy" toml:"use_proxy"`
	ProxyUri   string `json:"proxy_uri" yaml:"proxy_uri" toml:"proxy_uri"`
	Port       uint   `json:"port" yaml:"port" toml:"port"`
	StaticPort uint   `json:"static_port" yaml:"static_port" toml:"static_port"`
}

type Web struct {
	Static         Static `json:"static" yaml:"static" toml:"static"`
	Title          string `json:"title" yaml:"title" toml:"title"`
	JavascriptPath string `json:"javascript_path" yaml:"javascript_path" toml:"javascript_path"`
	ImagePath      string `json:"image_path" yaml:"image_path" toml:"image_path"`
	CssPath        string `json:"css_path" yaml:"css_path" toml:"css_path"`
}

type Static struct {
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
	Port  int   `json:"port" yaml:"port" toml:"port"`
	Delay int   `json:"delay" yaml:"delay" toml:"delay"`
	Auth  *Auth `json:"auth" yaml:"auth" toml:"auth"`
	// JWT       *rest.JWTData   `json:"jwt" yaml:"jwt" toml:"jwt"`
	// Static    *rest.Static    `json:"static" yaml:"static" toml:"static"`
	Path string `json:"-" yaml:"-" toml:"-"`
	// Resources []rest.Resource `json:"resources" yaml:"resources" toml:"resources"`
	// URLs      []rest.URL      `json:"urls" yaml:"urls" toml:"urls"`
	EnableLog bool `json:"enable_log" yaml:"enable_log" toml:"enable_log"`
}

type Server struct {
	Env                 string `json:"env" yaml:"env" toml:"env"`
	SessionID           string `json:"sessionID" yaml:"sessionID" toml:"sessionID"`
	Port                int    `json:"port" yaml:"port" toml:"port"`
	PageSize            int    `default:"20" json:"page_size" yaml:"page_size" toml:"page_size"`
	MaxPageSize         int    `default:"200" json:"max_page_size" yaml:"max_page_size" toml:"max_page_size"`
	MinPageSize         int    `default:"10" json:"min_page_size" yaml:"min_page_size" toml:"min_page_size"`
	MinOrder            int    `default:"0" json:"min_order" yaml:"min_order" toml:"min_order"`
	MaxOrder            int    `default:"10000" json:"max_order" yaml:"max_order" toml:"max_order"`
	MaxNameLength       int    `default:"100" json:"max_name_length" yaml:"max_name_length" toml:"max_name_length"`
	MaxContentLength    int    `default:"10000" json:"max_content_length" yaml:"max_content_length" toml:"max_content_length"`
	MaxArticleCateCount int    `default:"6" json:"max_article_cate_count" yaml:"max_article_cate_count" toml:"max_article_cate_count"`
	MaxCommentLength    int    `default:"5000" json:"max_comment_length" yaml:"max_comment_length" toml:"max_comment_length"`
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
