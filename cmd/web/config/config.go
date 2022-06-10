package config

type Config struct {
	*FlagConfig
	*CommonConfig
}

type FlagConfig struct {
	Port string `mapstructure:"port"`
}

type CommonConfig struct {
	DownloadPattern1 string `mapstructure:"downloadPattern1"`
	DownloadPattern2 string `mapstructure:"downloadPattern2"`
}
