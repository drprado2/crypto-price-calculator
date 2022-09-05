package configs

import (
	"log"
	"strings"
	"sync"

	"github.com/gosidekick/goconfig"
)

const (
	DeveloperEnvironment  = "dev"
	TestEnvironment       = "test"
	ProductionEnvironment = "prod"
)

var (
	doOnce sync.Once
	env    *Configuration
)

type Configuration struct {
	ServerPort            int    `cfg:"SERVER_PORT" cfgDefault:"5050" cfgRequired:"true"`
	ServerHost            string `cfg:"SERVER_HOST" cfgDefault:"localhost" cfgRequired:"true"`
	ServiceName           string `cfg:"SERVICE_NAME" cfgDefault:"Crypto Price Calculator" cfgRequired:"true"`
	ServerEndpointTimeout string `cfg:"SERVER_ENDPOINT_TIMEOUT" cfgDefault:"600s" cfgRequired:"true"`
	ServerEnvironment     string `cfg:"SERVER_ENVIRONMENT" cfgDefault:"dev" cfgRequired:"true"`
	SystemVersion         string `cfg:"SYSTEM_VERSION" cfgDefault:"UNKNOWN" cfgRequired:"true"`
	AppName               string `cfg:"APP_NAME" cfgDefault:"crypto-calculator" cfgRequired:"true"`
	DbUser                string `cfg:"DB_USER" cfgDefault:"admin" cfgRequired:"true"`
	DbPass                string `cfg:"DB_PASS" cfgDefault:"Postgres2019!" cfgRequired:"true"`
	DbName                string `cfg:"DB_NAME" cfgDefault:"crypto-price-calculator" cfgRequired:"true"`
	DbSchema              string `cfg:"DB_SCHEMA" cfgDefault:"crypto-price-calculator" cfgRequired:"true"`
	DbHost                string `cfg:"DB_HOST" cfgDefault:"localhost" cfgRequired:"true"`
	DbPort                int    `cfg:"DB_PORT" cfgDefault:"5432" cfgRequired:"true"`
	JaegerEndpoint        string `cfg:"JAEGER_ENDPOINT" cfgDefault:"http://localhost:14250/api/traces" cfgRequired:"false"`
	EnableTracing         bool   `cfg:"ENABLE_TRACING" cfgDefault:"true" cfgRequired:"false"`
	VwapWindowSize        int    `cfg:"VWAP_WINDOW_SIZE" cfgDefault:"200" cfgRequired:"false"`

	// coinbase
	CoinbaseWsEndpoint string `cfg:"COINBASE_WS_ENDPOINT" cfgDefault:"wss://ws-feed.exchange.coinbase.com" cfgRequired:"true"`
	ProductIds         string `cfg:"PRODUCT_IDS" cfgDefault:"BTC-USD,ETH-USD,ETH-BTC" cfgRequired:"true"`
	MatchesChannel     string `cfg:"MATCHES_CHANNEL" cfgDefault:"matches" cfgRequired:"true"`
}

func (c *Configuration) GetProductIds() []string {
	const separator = ","
	return strings.Split(c.ProductIds, separator)
}

func Get() *Configuration {
	doOnce.Do(func() {
		env = &Configuration{}
		err := goconfig.Parse(env)
		if err != nil {
			log.Fatal(err)
		}
	})
	return env
}

func Reset() *Configuration {
	doOnce = sync.Once{}
	return Get()
}
