package config

type Config struct {
	*FlagConfig
	*CommonConfig
}

type FlagConfig struct {
	Port         string `mapstructure:"port"`
	Proxy        string `mapstructure:"proxy"`
	NasCacheFile string `mapstructure:"nascachefile"`
	DataDir      string `mapstructure:"dataDir"`
	ScanDir      string `mapstructure:"scandir"`
}

type CommonConfig struct {
	DownloadPattern1 string `mapstructure:"downloadpattern1"`
	DownloadPattern2 string `mapstructure:"downloadpattern2"`
}
