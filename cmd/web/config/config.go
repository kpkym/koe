package config

type Config struct {
	Common struct {
		Port           string `mapstructure:"port"`
		NasCacheFile   string `mapstructure:"nascachefile"`
		SqliteDataFile string `mapstructure:"sqlitedatafile"`
		ScanDir        string `mapstructure:"scandir"`
	} `mapstructure:"common"`

	DownloadPattern1 string `mapstructure:"downloadpattern1"`
	DownloadPattern2 string `mapstructure:"downloadpattern2"`
}
