package config

import (
	"github.com/kpkym/koe/cmd/web/config/db"
	"github.com/kpkym/koe/global"
	"github.com/kpkym/koe/model/domain"
	"github.com/kpkym/koe/utils"
	"io/ioutil"
)

func Init(configFile string) {
	configString := string(utils.IgnoreErr(ioutil.ReadFile(utils.GetFileBaseOnPwd(configFile))))
	global.AddConfig(&global.Config{
		CommonConfig:     utils.GetTomlConfig[global.CommonConfig](configString, "common"),
		FixConfig:        getCommonConfig(),
		SqliteConfig:     utils.GetTomlConfig[global.SqliteConfig](configString, "sqlite"),
		PostgresqlConfig: utils.GetTomlConfig[global.PostgresqlConfig](configString, "postgresql"),
	})

	global.AddDB(db.Init())
	global.AddSettings(loadSettingsFromDB())
}

func loadSettingsFromDB() *domain.Settings {
	settings := &domain.Settings{ID: 1, ScanDirs: []byte("[]")}
	global.GetServiceContext().DB.FirstOrCreate(settings)
	return settings
}
