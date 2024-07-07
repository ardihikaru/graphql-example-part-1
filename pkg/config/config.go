package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/spf13/viper"
)

type General struct {
	BuildMode string `mapstructure:"buildMode"`
}

type Log struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type LogPublisher struct {
	Notice bool `mapstructure:"notice"`
	Error  bool `mapstructure:"error"`
	Info   bool `mapstructure:"info"`
	Warn   bool `mapstructure:"warn"`
	Debug  bool `mapstructure:"debug"`
}

type Http struct {
	Address        string        `mapstructure:"address"`
	Port           int           `mapstructure:"port"`
	RequestTimeout time.Duration `mapstructure:"requestTimeout"`
	ReadTimeout    time.Duration `mapstructure:"readTimeout"`  // The maximum time to wait while writing data to the server
	WriteTimeout   time.Duration `mapstructure:"writeTimeout"` // The maximum time to wait while reading data from the server
	HttpClientTLS  bool          `mapstructure:"httpClientTLS"`
	LogHttpRequest bool          `mapstructure:"logHttpRequest"`
}

// Cors defines cors-related configurations
type Cors struct {
	AllowedOrigins   []string `mapstructure:"allowedOrigins"`
	AllowedMethods   []string `mapstructure:"allowedMethods"`
	AllowedHeaders   []string `mapstructure:"allowedHeaders"`
	ExposedHeaders   []string `mapstructure:"exposedHeaders"`
	AllowCredentials bool     `mapstructure:"allowCredentials"`
	MaxAge           int      `mapstructure:"maxAge"`
	Debug            bool     `mapstructure:"debug"`
}

// JwtAuth defines JWT authentication related
type JwtAuth struct {
	Secret       string                 `mapstructure:"jwtSecret"`
	Algorithm    jwa.SignatureAlgorithm `mapstructure:"jwtAlgorithm"`
	ExpiredInSec int                    `mapstructure:"jwtExpiredInSec"`
}

// DbMySQL defines the database connection for MySQL database
type DbMySQL struct {
	Host   string `mapstructure:"host"`
	Port   string `mapstructure:"port"`
	User   string `mapstructure:"user"`
	Pass   string `mapstructure:"pass"`
	DbName string `mapstructure:"dbName"`
}

// Enforcer defines the enforcer configuration
type Enforcer struct {
	ModelFile string `mapstructure:"modelFile"`
	TableName string `mapstructure:"tableName"`
}

// GraphQL defines the GraphQL configuration
type GraphQL struct {
	PublicFunctions []string `mapstructure:"publicFunctions"`
}

type Encryption struct {
	PrivateKey    string `mapstructure:"privateKey"`
	HashCost      int    `mapstructure:"hashCost"`
	PrivateKeyRSA *rsa.PrivateKey
	PublicKeyRSA  *rsa.PublicKey
}

type Config struct {
	General      General      `mapstructure:"general"`
	Log          Log          `mapstructure:"log"`
	LogPublisher LogPublisher `mapstructure:"logPublisher"`
	Http         Http         `mapstructure:"http"`
	Cors         Cors         `mapstructure:"cors"`
	JwtAuth      JwtAuth      `mapstructure:"jwtAuth"`
	Encryption   Encryption   `mapstructure:"encryption"`
	DbMySQL      DbMySQL      `mapstructure:"dbMysql"`
	Enforcer     Enforcer     `mapstructure:"enforcer"`
	GraphQL      GraphQL      `mapstructure:"graphQL"`
}

// Validate validates any miss configurations or missing configs
func (cfg *Config) Validate() error {
	if cfg.Encryption.PrivateKeyRSA == nil {
		return fmt.Errorf("missing value of private key RSA")
	}

	if cfg.Encryption.PublicKeyRSA == nil {
		return fmt.Errorf("missing value of public key RSA")
	}

	// TODO: adds more here

	return nil
}

// Get gets config object
func Get() (*Config, error) {
	cfg, err := Load()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// BuildEncryptionKeys builds private key and public key RSA
func (cfg *Config) BuildEncryptionKeys() error {
	pemData, err := os.ReadFile(cfg.Encryption.PrivateKey)
	if err != nil {
		return fmt.Errorf("read key file: %s", err.Error())
	}
	block, _ := pem.Decode(pemData)
	if block == nil {
		return fmt.Errorf("bad key data: not PEM-encoded")
	}
	if got, want := block.Type, "RSA PRIVATE KEY"; got != want {
		return fmt.Errorf("unknown key type %q, want %q", got, want)
	}
	// Decode the RSA private key
	privateKeyRSA, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("bad private key: %s", err.Error())
	}

	// sets the keys
	cfg.Encryption.PrivateKeyRSA = privateKeyRSA
	cfg.Encryption.PublicKeyRSA = &privateKeyRSA.PublicKey

	return nil
}

// Load loads config from the config.yaml
func Load() (*Config, error) {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	err := v.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, err
}
