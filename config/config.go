package config

import (
	"bytes"
	"path/filepath"
	"strings"

	"github.com/0xPolygonHermez/zkevm-aggregator/aggregator"
	"github.com/0xPolygonHermez/zkevm-aggregator/etherman"
	"github.com/0xPolygonHermez/zkevm-aggregator/event"
	"github.com/0xPolygonHermez/zkevm-aggregator/log"
	"github.com/0xPolygonHermez/zkevm-aggregator/metrics"
	"github.com/0xPolygonHermez/zkevm-ethtx-manager/ethtxmanager"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

const (
	// FlagYes is the flag for yes.
	FlagYes = "yes"
	// FlagCfg is the flag for cfg.
	FlagCfg = "cfg"
	// FlagNetwork is the flag for the network name. Valid values: ["testnet", "mainnet", "custom"].
	FlagNetwork = "network"
	// FlagCustomNetwork is the flag for the custom network file. This is required if --network=custom
	FlagCustomNetwork = "custom-network-file"
	// FlagKeyStorePath is the path of the key store file containing the private key of the account going to sing and approve the tokens
	FlagKeyStorePath = "key-store-path"
	// FlagPassword is the password needed to decrypt the key store
	FlagPassword = "password"
	// FlagMigrations is the flag for migrations.
	FlagMigrations = "migrations"
	// FlagDocumentationFileType is the flag for the choose which file generate json-schema
	FlagDocumentationFileType = "config-file"
)

/*
Config represents the configuration of the entire Hermez Node
The file is [TOML format]
You could find some examples:
  - `config/environments/local/local.node.config.toml`: running a permisionless node
  - `config/environments/mainnet/node.config.toml`
  - `config/environments/public/node.config.toml`
  - `test/config/test.node.config.toml`: configuration for a trusted node used in CI

[TOML format]: https://en.wikipedia.org/wiki/TOML
*/
type Config struct {
	// Configuration of the etherman (client for access L1)
	Etherman etherman.Config
	// Configuration for ethereum transaction manager
	EthTxManager ethtxmanager.Config
	// Configuration of the aggregator service
	Aggregator aggregator.Config `mapstructure:"Aggregator"`
	// Configuration of the genesis of the network. This is used to known the initial state of the network
	NetworkConfig NetworkConfig
	// Configuration of the metrics service, basically is where is going to publish the metrics
	Metrics metrics.Config
	// Configuration of the event database connection
	EventLog event.Config
}

// Default parses the default configuration values.
func Default() (*Config, error) {
	var cfg Config
	viper.SetConfigType("toml")

	err := viper.ReadConfig(bytes.NewBuffer([]byte(DefaultValues)))
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&cfg, viper.DecodeHook(mapstructure.TextUnmarshallerHookFunc()))
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Load loads the configuration
func Load(ctx *cli.Context, loadNetworkConfig bool) (*Config, error) {
	cfg, err := Default()
	if err != nil {
		return nil, err
	}
	configFilePath := ctx.String(FlagCfg)
	if configFilePath != "" {
		dirName, fileName := filepath.Split(configFilePath)

		fileExtension := strings.TrimPrefix(filepath.Ext(fileName), ".")
		fileNameWithoutExtension := strings.TrimSuffix(fileName, "."+fileExtension)

		viper.AddConfigPath(dirName)
		viper.SetConfigName(fileNameWithoutExtension)
		viper.SetConfigType(fileExtension)
	}
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("ZKEVM_AGGREGATOR")
	err = viper.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if ok {
			log.Infof("config file not found")
		} else {
			log.Infof("error reading config file: ", err)
			return nil, err
		}
	}

	decodeHooks := []viper.DecoderConfigOption{
		// this allows arrays to be decoded from env var separated by ",", example: MY_VAR="value1,value2,value3"
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(mapstructure.TextUnmarshallerHookFunc(), mapstructure.StringToSliceHookFunc(","))),
	}

	err = viper.Unmarshal(&cfg, decodeHooks...)
	if err != nil {
		return nil, err
	}

	if loadNetworkConfig {
		// Load genesis parameters
		cfg.loadNetworkConfig(ctx)
	}
	return cfg, nil
}
