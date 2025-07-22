package config

import (
	"fmt"
	"os"
	"time"
	
	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	API      APIConfig      `mapstructure:"api"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Tracker  TrackerConfig  `mapstructure:"tracker"`
}

// AppConfig holds application configuration
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
	LogLevel    string `mapstructure:"log_level"`
	LogFormat   string `mapstructure:"log_format"`
}

// DatabaseConfig holds Neo4j database configuration
type DatabaseConfig struct {
	URI      string `mapstructure:"uri"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

// APIConfig holds blockchain API configuration
type APIConfig struct {
	Provider        string            `mapstructure:"provider"`
	Keys            map[string]string `mapstructure:"keys"`
	RateLimit       int               `mapstructure:"rate_limit"`
	Timeout         time.Duration     `mapstructure:"timeout"`
	MaxRetries      int               `mapstructure:"max_retries"`
	RetryDelay      time.Duration     `mapstructure:"retry_delay"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string        `mapstructure:"host"`
	Port     int           `mapstructure:"port"`
	Password string        `mapstructure:"password"`
	DB       int           `mapstructure:"db"`
	TTL      time.Duration `mapstructure:"ttl"`
}

// TrackerConfig holds tracker specific configuration
type TrackerConfig struct {
	DefaultNetwork      string   `mapstructure:"default_network"`
	SupportedNetworks   []string `mapstructure:"supported_networks"`
	BatchSize           int      `mapstructure:"batch_size"`
	MaxDepth            int      `mapstructure:"max_depth"`
	ConcurrentWorkers   int      `mapstructure:"concurrent_workers"`
}

// Load loads configuration from files and environment variables
func Load(configPath string) (*Config, error) {
	// Set config file
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
		viper.AddConfigPath("/etc/wallet-tracker")
	}
	
	// Set environment variable prefix
	viper.SetEnvPrefix("WALLET_TRACKER")
	viper.AutomaticEnv()
	
	// Set defaults
	setDefaults()
	
	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		// It's ok if config file doesn't exist, we'll use defaults and env vars
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}
	
	// Unmarshal config
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}
	
	// Override with environment variables
	overrideFromEnv(&config)
	
	// Validate config
	if err := validate(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}
	
	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// App defaults
	viper.SetDefault("app.name", "wallet-tracker")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.environment", "development")
	viper.SetDefault("app.log_level", "info")
	viper.SetDefault("app.log_format", "json")
	
	// Database defaults
	viper.SetDefault("database.uri", "neo4j://localhost:7687")
	viper.SetDefault("database.username", "neo4j")
	viper.SetDefault("database.password", "letmein")
	viper.SetDefault("database.database", "neo4j")
	
	// API defaults
	viper.SetDefault("api.provider", "blockchain.info")
	viper.SetDefault("api.rate_limit", 10)
	viper.SetDefault("api.timeout", "30s")
	viper.SetDefault("api.max_retries", 3)
	viper.SetDefault("api.retry_delay", "1s")
	
	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.ttl", "1h")
	
	// Tracker defaults
	viper.SetDefault("tracker.default_network", "BTC")
	viper.SetDefault("tracker.supported_networks", []string{"BTC", "ETH"})
	viper.SetDefault("tracker.batch_size", 100)
	viper.SetDefault("tracker.max_depth", 10)
	viper.SetDefault("tracker.concurrent_workers", 5)
}

// overrideFromEnv overrides config values from environment variables
func overrideFromEnv(config *Config) {
	// Neo4j compatibility with existing .env
	if username := os.Getenv("NEO4J_USERNAME"); username != "" {
		config.Database.Username = username
	}
	if password := os.Getenv("NEO4J_PASS"); password != "" {
		config.Database.Password = password
	}
}

// validate validates the configuration
func validate(config *Config) error {
	if config.Database.URI == "" {
		return fmt.Errorf("database URI is required")
	}
	
	if config.Database.Username == "" {
		return fmt.Errorf("database username is required")
	}
	
	if config.Database.Password == "" {
		return fmt.Errorf("database password is required")
	}
	
	if config.Tracker.BatchSize <= 0 {
		return fmt.Errorf("tracker batch size must be positive")
	}
	
	if config.Tracker.ConcurrentWorkers <= 0 {
		return fmt.Errorf("tracker concurrent workers must be positive")
	}
	
	return nil
}
