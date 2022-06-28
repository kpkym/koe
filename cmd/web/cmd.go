package web

import (
	"fmt"
	"github.com/kpkym/koe/cmd/web/config"
	"github.com/kpkym/koe/colly"
	"github.com/kpkym/koe/dao/cache"
	"github.com/kpkym/koe/global"
	"github.com/kpkym/koe/model/domain"
	"github.com/kpkym/koe/model/others"
	"github.com/kpkym/koe/router"
	"github.com/kpkym/koe/utils/koe"
	"github.com/spf13/cobra"
)

// //go:embed dist
// var dist embed.FS

var (
	configFile string
)

var (
	Cmd = &cobra.Command{
		Use:   "web",
		Short: "启动web服务",
		Run: func(_ *cobra.Command, _ []string) {
			config.Init(configFile)
			koe.BuildTree()
			trees, _ := cache.NewMapCache[string, []*others.Node]().Get("trees")
			// koe.RemoveIncomplete()
			koe.RemoveNotInTrees(trees, func(needCrawlCodes []string) {
				colly.C(needCrawlCodes, func(workDomain *domain.WorkDomain) {
					global.GetServiceContext().DB.Create(workDomain)
				})
			})
			web()
		},
	}
)

func init() {
	Cmd.Flags().StringVarP(&configFile, "config", "c", "config.toml", "配置文件")
}

func web() {
	serve := router.GetGinServe()
	// serve.NoRoute(func(context *gin.Context) {
	// 	context.FileFromFS("/"+context.Request.RequestURI, http.FS(utils.IgnoreErr(fs.Sub(dist, "dist"))))
	// })

	serve.Run(fmt.Sprintf(":%d", global.GetServiceContext().Config.CommonConfig.Port))
}
