package main

import (
	"github.com/allanpk716/rssdownloader.common/model"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/log"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"strconv"
)

func main() {
	var (
		err       error
		vipConfig *viper.Viper
	)
	// -------------------------------------------------------------
	// 加载配置
	vipConfig, err = InitConfigure()
	if err != nil {
		log.Errorln("InitConfigure:", err)
		return
	}
	// 缓存配置
	dockerDownloaderInfos = make(map[string]model.DockerDownloaderInfo)
	ViperConfig2Cache(vipConfig, &configs, &rssProxyInfos, &biliBiliInfos, dockerDownloaderInfos)
	// -------------------------------------------------------------
	// 开启一个协程做定时的更新
	// 这里直接使用 cron ，他会自动开启一个协程
	//任务还没执行完，下一次执行时间到来，默认任务异步执行，也就是说可能同一个任务在当前有2个在执行中
	//c := cron.New()
	//任务还没执行完，下一次执行时间到来，下一次执行必须要等这一次执行完才能执行
	//c := cron.New(cron.WithChain(cron.DelayIfStillRunning(cron.DefaultLogger)))
	//任务还没执行完，下一次执行时间到来，下一次执行就跳过不执行
	c := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))
	// 定时器
	entryID, err := c.AddFunc("@every " + configs.EveryTime, func() {
		MainDownloader(configs, rssProxyInfos, biliBiliInfos)
	})
	if err != nil {
		log.Errorln("cron entryID:", entryID, "Error:", err)
		return
	}
	// 立即触发第一次的更新
	MainDownloader(configs, rssProxyInfos, biliBiliInfos)
	c.Start()
	defer c.Stop()
	// -------------------------------------------------------------
	// 初始化 gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// 默认所有都通过
	r.Use(cors.Default())
	// 注册路由事件
	//..
	// 启动 gin
	log.Infoln("RSSDownloader Start at port " + strconv.Itoa(configs.ListenPort) + " ...")
	err = r.Run(":" + strconv.Itoa(configs.ListenPort))
	if err != nil {
		log.Fatal("Start RSSProxy Server At Port", configs.ListenPort, "Fatal", err)
	}
}

var (
	configs               model.Configs
	rssProxyInfos         model.RSSProxyInfos
	biliBiliInfos         model.BiliBiliInfos
	dockerDownloaderInfos map[string]model.DockerDownloaderInfo
)