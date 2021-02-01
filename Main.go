package main

import (
	"context"
	"errors"
	"github.com/bitfield/script"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mmcdole/gofeed"
	"github.com/prometheus/common/log"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"os"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
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
	ViperConfig2Cache(vipConfig, &configs, &rssProxyInfos, &biliBiliInfos)
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

func MainDownloader(configs Configs, rssProxyInfos RSSProxyInfos, biliBiliInfos BiliBiliInfos)  {
	for _, oneRSSInfos := range rssProxyInfos.RSSInfos {
		DownloadFromOneFeed(configs, oneRSSInfos)
	}

	for _, oneBiliBiliUserInfos := range biliBiliInfos.BiliBiliUserInfos {
		DownloadFromOneFeed(configs, oneBiliBiliUserInfos)
	}
}

func SelectDownloadInfo(rssInfo interface{}) (DownloadInfo, error) {
	var downloadInfo DownloadInfo
	// TODO 后续新增的订阅类型，需要在这里新增对应的 Switch 语句
	switch rssInfo.(type) {
	case RSSInfo:
		// RSSProxyInfos
		// http://127.0.0.1:1201/rss?key=巫师财经
		downloadInfo = DownloadInfo{
			FolderName: rssInfo.(RSSInfo).FolderName,
			RSSUrl: configs.RSSProxyAddress + "/rss?key=" + rssInfo.(RSSInfo).RSSInfosName,
			DownloadRoot: rssInfo.(RSSInfo).DownloadRoot,
			UseProxy: rssInfo.(RSSInfo).UseProxy,
		}

	case BiliBiliUserInfo:
		// BiliBiliInfos
		// http://192.168.50.135:1200/bilibili/user/video/258150656
		downloadInfo = DownloadInfo{
			FolderName: rssInfo.(BiliBiliUserInfo).FolderName,
			RSSUrl: configs.RSSHubAddress + "/bilibili/user/video/" + rssInfo.(BiliBiliUserInfo).UserID,
			DownloadRoot: rssInfo.(BiliBiliUserInfo).DownloadRoot,
			UseProxy: rssInfo.(BiliBiliUserInfo).UseProxy,
		}
	default:
		return downloadInfo, errors.New("RSS info type not support")
	}

	return downloadInfo, nil
}

func DownloadFromOneFeed(configs Configs, rssInfo interface{}) {
	var downloadInfo DownloadInfo
	var err error
	downloadInfo, err = SelectDownloadInfo(rssInfo)
	if err != nil {
		log.Errorln("SelectDownloadInfo:", err, rssInfo)
		return
	}
	// 解析 RSS 信息
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configs.ReadRSSTimeOut) * time.Second)
	defer cancel()
	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(downloadInfo.RSSUrl, ctx)
	if err != nil {
		log.Errorln(" ")
		log.Errorln("gofeed:", downloadInfo.RSSUrl, err)
		return
	}
	// 下载找到的所有的信息
	for _, oneItem := range feed.Items {
		StartDownload(oneItem, downloadInfo)
	}
}

func StartDownload(item *gofeed.Item, downloadInfo DownloadInfo) {
	// 还需要拼接具体某人的目录
	var tmpDownloadPath string
	if runtime.GOOS == "windows" {
		tmpDownloadPath = downloadInfo.DownloadRoot + "\\\\" +downloadInfo.FolderName
	} else {
		tmpDownloadPath = path.Join(downloadInfo.DownloadRoot, downloadInfo.FolderName)
	}
	if Exists(tmpDownloadPath) == false {
		err := os.MkdirAll(tmpDownloadPath, os.ModePerm)
		if err != nil {
			log.Errorln("os.MkdirAll:", tmpDownloadPath, err)
			return
		}
	}
	// 当前路径下面的是否有下载好的文件，或者是否有 .download 的正在下载的文件
	nowVideoName := item.PublishedParsed.Format("2006-01-02") + "_" + item.Title
	nowVideoName = strings.TrimSpace(nowVideoName)
	// 去除 Windows 下不允许出现在文件名中的特殊字符
	replaceWindowErrorChar, _ := regexp.Compile(`[\\\\/:*?\"<>|]`)
	nowVideoName = replaceWindowErrorChar.ReplaceAllString(nowVideoName, "-")
	// 去除空格
	nowVideoName = strings.ReplaceAll(nowVideoName, " ", "")
	// 搜索的时候使用一个通用正则表达式来找文件
	filterFiles := regexp.MustCompile(nowVideoName + `[\s]*`)
	p := script.FindFiles(tmpDownloadPath).MatchRegexp(filterFiles)
	outName, err := p.String()
	if err != nil {
		log.Errorln("script.FindFiles:", nowVideoName, err)
		return
	}
	outName = strings.TrimSpace(outName)
	if outName == "" {
		// 为空就直接下载
	} else {
		// 可能读取到多个文件
		stringSlice := strings.Split(outName, "\n")
		if len(stringSlice) == 1 {
			// 只有一个文件的时候，需要判断后缀名，如果是  .download 那么就需要继续下载
			if path.Ext(outName) == annieFileExtension {
				// 继续下载
			} else {
				// 跳出，无需下载
				return
			}
		} else {
			// 其他情况也是直接下载
			// 比如大于一个个文件（不可能没有，上面排除了）
		}
	}

	log.Infoln("Download", nowVideoName, "Start")
	// 如果为空，则没找到那么就可以下载，注意，这里是单线程下载，所以用阻塞调用方法
	// 这里有个梗，annie 已经无法正常下载 youtube 的视频了···
	err = OneDownload(tmpDownloadPath, nowVideoName, item.Link, downloadInfo)
	if err != nil {
		log.Errorln("OneDownload:", err)
		return
	}
	log.Infoln("Download", nowVideoName, "End")
}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

var (
	configs Configs
	rssProxyInfos RSSProxyInfos
	biliBiliInfos BiliBiliInfos
)

const annieFileExtension = ".download"