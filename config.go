package main

import (
	"errors"
	"github.com/allanpk716/rssdownloader/model"
	"github.com/spf13/viper"
)

func InitConfigure() (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName("config") // 设置文件名称（无后缀）
	v.SetConfigType("yaml")   // 设置后缀名 {"1.6以后的版本可以不设置该后缀"}
	v.AddConfigPath(".")      // 设置文件所在路径

	err := v.ReadInConfig()
	if err != nil {
		return nil, errors.New("error reading config:" + err.Error())
	}

	return v, nil
}

func ViperConfig2Cache(config *viper.Viper, configs *model.Configs,
	rssProxyInfos *model.RSSProxyInfos, biliBiliInfos *model.BiliBiliInfos) {
	// ------------------------------------------------------------
	// 基础配置信息
	configs.ListenPort =  config.GetInt("ListenPort")
	configs.DownloadHttpProxy =  config.GetString("DownloadHttpProxy")
	configs.RSSHubAddress =  config.GetString("RSSHubAddress")
	configs.RSSProxyAddress =  config.GetString("RSSProxyAddress")
	configs.EveryTime =  config.GetString("EveryTime")
	configs.ReadRSSTimeOut =  config.GetInt("ReadRSSTimeOut")
	// ------------------------------------------------------------
	// 读取需要走 RSSProxy 的信息
	rssProxyInfos.DefaultDownloadRoot = config.GetString("RSSProxyInfos.DefaultDownloadRoot")
	rssProxyInfos.DefaultUseProxy = config.GetBool("RSSProxyInfos.DefaultUseProxy")
	rsshubInfos := config.GetStringMapString("RSSProxyInfos.RSSInfos")
	// 巫师财经 ： 具体的内容
	for k := range rsshubInfos {
		rssInfo := model.RSSInfo{
			FolderName: k,
			RSSInfosName: config.GetString("RSSProxyInfos.RSSInfos." + k + ".RSSInfosName"),
		}
		// 优先使用单独设置的是否使用代理
		keyRSSUseProxy := "RSSProxyInfos.RSSInfos.\" + k + \".UseProxy"
		if config.InConfig(keyRSSUseProxy) == true {
			rssInfo.UseProxy = config.GetBool(keyRSSUseProxy)
		} else {
			rssInfo.UseProxy = rssProxyInfos.DefaultUseProxy
		}
		// 优先使用单独设置的下载路径
		keyRSSDownloadRoot := "RSSProxyInfos.RSSInfos.\" + k + \".UseProxy"
		if config.InConfig(keyRSSDownloadRoot) == true {
			rssInfo.DownloadRoot = config.GetString(keyRSSDownloadRoot)
		} else {
			rssInfo.DownloadRoot = rssProxyInfos.DefaultDownloadRoot
		}

		rssProxyInfos.RSSInfos = append(rssProxyInfos.RSSInfos, rssInfo)

	}
	// ------------------------------------------------------------
	// 读取 BiliBiliInfos 的 User 信息
	biliBiliInfos.DefaultDownloadRoot = config.GetString("BiliBiliInfos.DefaultDownloadRoot")
	biliBiliInfos.DefaultUseProxy = config.GetBool("BiliBiliInfos.DefaultUseProxy")
	bilabialUserInfos := config.GetStringMapString("BiliBiliInfos.BiliBiliUserInfos")
	// 李永乐 ： 具体的内容
	for k := range bilabialUserInfos {
		 biliUserInfo := model.BiliBiliUserInfo{
			FolderName: k,
			UserID: config.GetString("BiliBiliInfos.BiliBiliUserInfos." + k + ".UserID"),
		}
		// 优先使用单独设置的是否使用代理
		keyRSSUseProxy := "BiliBiliInfos.BiliBiliUserInfos." + k + ".UseProxy"
		if config.InConfig(keyRSSUseProxy) == true {
			biliUserInfo.UseProxy = config.GetBool(keyRSSUseProxy)
		} else {
			biliUserInfo.UseProxy = biliBiliInfos.DefaultUseProxy
		}
		// 优先使用单独设置的下载路径
		keyRSSDownloadRoot := "BiliBiliInfos.BiliBiliUserInfos." + k + ".DownloadRoot"
		if config.InConfig(keyRSSDownloadRoot) == true {
			biliUserInfo.DownloadRoot = config.GetString(keyRSSDownloadRoot)
		} else {
			biliUserInfo.DownloadRoot = biliBiliInfos.DefaultDownloadRoot
		}

		biliBiliInfos.BiliBiliUserInfos = append(biliBiliInfos.BiliBiliUserInfos, biliUserInfo)
	}
	// ------------------------------------------------------------
}


