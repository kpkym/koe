package global

type Config struct {
	*CommonConfig
	*FixConfig
	*SqliteConfig
	*PostgresqlConfig
}

type CommonConfig struct {
	Port    int    `mapstructure:"port"`
	DataDir string `mapstructure:"datadir"`
}

type FixConfig struct {
	DownloadPattern1 string `mapstructure:"downloadPattern1"`
	DownloadPattern2 string `mapstructure:"downloadPattern2"`
}

type SqliteConfig struct {
	Filename string `mapstructure:"filename"`
}

type PostgresqlConfig struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Dbname   string `mapstructure:"dbname"`
	Port     int    `mapstructure:"port"`
	Sslmode  string `mapstructure:"sslmode"`
	Timezone string `mapstructure:"timezone"`
}
