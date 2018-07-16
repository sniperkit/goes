package config

import (
	"fmt"
	"os"

	// external
	"github.com/jinzhu/configor"
	"github.com/k0kubun/pp"
)

var (
	Global *Config
)

func PrettyPrint(msg interface{}) {
	pp.Println(msg)
}

func New(cfgFiles ...string) *Config {
	cfg := &Config{}
	configor.New(&configor.Config{ErrorOnUnmatchedKeys: false, Debug: false, Verbose: true, ENVPrefix: "SNK"}).Load(cfg, cfgFiles...)

	if cfg == nil {
		fmt.Println("error while loading configuration file")
		os.Exit(1)
	}

	var url string
	switch cfg.Store.Dialect {
	case "mysql", "postgres", "mssql":
		url = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			cfg.Store.User,
			cfg.Store.Password,
			cfg.Store.Host,
			cfg.Store.Port,
			cfg.Store.Database,
			cfg.Store.Charset)
	case "sqlite", "sqlite3":
		fallthrough
	default:
		cfg.Store.Dialect = "sqlite3"
		url = fmt.Sprintf("%s/%s.db", cfg.App.Dir.Store, cfg.Store.Database)
	}

	if cfg.Store.DSN == "" {
		cfg.Store.DSN = url
	}

	return cfg
}

func (c *Config) WithDB(dbCfg *Database) *Config {

	var url string
	switch dbCfg.Dialect {
	case "mysql", "postgres", "mssql":
		url = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			dbCfg.User,
			dbCfg.Password,
			dbCfg.Host,
			dbCfg.Port,
			dbCfg.Database,
			dbCfg.Charset)
	case "sqlite", "sqlite3":
		fallthrough
	default:
		dbCfg.Dialect = "sqlite3"
		url = fmt.Sprintf("%s/%s.db", c.App.Dir.Store, dbCfg.Database)
	}

	if dbCfg.DSN == "" {
		dbCfg.DSN = url
	}

	c.Store = *dbCfg
	return c
}

func (c *Config) WithApi(apiCfg *Api) *Config {
	c.Api = *apiCfg
	return c
}

func (c *Config) WithServer(srvCfg *Server) *Config {
	c.Server = *srvCfg
	return c
}

func (c *Config) WithWebsocket(wsCfg *Websocket) *Config {
	c.Websocket = *wsCfg
	return c
}
