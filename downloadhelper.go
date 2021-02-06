package main

import (
	"context"
	"errors"
	"github.com/allanpk716/rssdownloader.common/model"
	"github.com/bitfield/script"
	"github.com/mmcdole/gofeed"
	"github.com/prometheus/common/log"
	"golang.org/x/crypto/ssh"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func UpdateDockerDownloader(configs model.Configs, downloadRoot string, dockerDownloaderInfo model.DockerDownloaderInfo) (string, error) {

	var err error
	config := &ssh.ClientConfig{
		Timeout:         time.Second,
		User:            dockerDownloaderInfo.DockerUserName,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	config.Auth = []ssh.AuthMethod{ssh.Password(dockerDownloaderInfo.DockerPassword)}
	sshClient, err := ssh.Dial("tcp", dockerDownloaderInfo.DockSSHAddress, config)
	if err != nil {
		return "", err
	}
	defer sshClient.Close()

	session, err := sshClient.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	assemblyCommand := ""
	for _, oneCommand := range dockerDownloaderInfo.UpdateCommands {
		// 替换关键信息
		tmpCommand := strings.ReplaceAll(oneCommand, ConstPhysicalmachinedownloadrootpath, downloadRoot)
		ContainerProxy := "-e HTTP_PROXY=" + configs.DownloadHttpProxy
		tmpCommand = strings.ReplaceAll(tmpCommand, ConstContainerProxy, ContainerProxy)

		tmpCommand = replaceOutSideAPPOrFolderLocation(tmpCommand, dockerDownloaderInfo)

		assemblyCommand += tmpCommand
		assemblyCommand += ";"
	}

	combo, err := session.CombinedOutput(assemblyCommand)
	if err != nil {
		return string(combo), err
	}

	return "", nil
}

func OneDownload(nowFileName, downloadURL string, downloadInfo model.DownloadInfo,
	dockerDownloaderInfo model.DockerDownloaderInfo) (string, error) {
	var err error
	config := &ssh.ClientConfig{
		Timeout:         time.Second,
		User:            dockerDownloaderInfo.DockerUserName,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	config.Auth = []ssh.AuthMethod{ssh.Password(dockerDownloaderInfo.DockerPassword)}
	sshClient, err := ssh.Dial("tcp", dockerDownloaderInfo.DockSSHAddress, config)
	if err != nil {
		return "", err
	}
	defer sshClient.Close()

	session, err := sshClient.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	// TODO 先实现 Youtube-dl 的功能，后续需分析这个部分，支持多种 Docker 下载器
	assemblyCommand := ""
	for _, oneCommand := range dockerDownloaderInfo.DownloadCommands {
		// 替换关键信息
		nowDesPath := path.Join(downloadInfo.DownloadRoot, downloadInfo.FolderName)
		tmpCommand := strings.ReplaceAll(oneCommand, ConstPhysicalmachinedownloadrootpath, nowDesPath)

		httpAdd := downloadInfo.DownloadHttpProxy
		if downloadInfo.UseProxy == false {
			httpAdd = ""
		}
		tmpCommand = strings.ReplaceAll(tmpCommand, ConstHttpadd, httpAdd)
		tmpCommand = strings.ReplaceAll(tmpCommand, ConstDownloadurl, downloadURL)
		tmpCommand = strings.ReplaceAll(tmpCommand, ConstNowfilename, nowFileName)

		tmpCommand = replaceOutSideAPPOrFolderLocation(tmpCommand, dockerDownloaderInfo)

		assemblyCommand += tmpCommand
		assemblyCommand += ";"
	}

	combo, err := session.CombinedOutput(assemblyCommand)
	if err != nil {
		return string(combo), err
	}

	return "", nil
}

func MainDownloader(configs model.Configs, rssProxyInfos model.RSSProxyInfos, biliBiliInfos model.BiliBiliInfos)  {

	// 先进行 downloader 的统一更新
	log.Infoln("Docker Downloader Update Start")

	updateTmpDownloadRoot := ""
	if rssProxyInfos.DefaultDownloadRoot != "" {
		updateTmpDownloadRoot = rssProxyInfos.DefaultDownloadRoot
	}
	if biliBiliInfos.DefaultDownloadRoot != "" {
		updateTmpDownloadRoot = biliBiliInfos.DefaultDownloadRoot
	}

	for _, dockerDownloaderInfo := range dockerDownloaderInfos {
		log.Infoln("Update Docker Downloader:", dockerDownloaderInfo.Name, "Start")
		outstring, err := UpdateDockerDownloader(configs, updateTmpDownloadRoot, dockerDownloaderInfo)
		if err != nil {
			log.Errorln("UpdateDockerDownloader OutString:", outstring)
			log.Errorln("UpdateDockerDownloader:", err)
			continue
		}
		log.Infoln("Update Docker Downloader:", dockerDownloaderInfo.Name, "End")
	}
	log.Infoln("Docker Downloader Update End")

	log.Infoln("Download RSS From RSSProxyInfos Start")
	for _, oneRSSInfos := range rssProxyInfos.RSSInfos {
		DownloadFromOneFeed(configs, oneRSSInfos)
	}
	log.Infoln("Download RSS From RSSProxyInfos End")

	log.Infoln("Download RSS From BiliBiliInfos Start")
	for _, oneBiliBiliUserInfos := range biliBiliInfos.BiliBiliUserInfos {
		DownloadFromOneFeed(configs, oneBiliBiliUserInfos)
	}
	log.Infoln("Download RSS From BiliBiliInfos End")
}

func SelectDownloadInfo(rssInfo interface{}) (model.DownloadInfo, error) {
	var downloadInfo model.DownloadInfo
	// TODO 后续新增的订阅类型，需要在这里新增对应的 Switch 语句
	switch rssInfo.(type) {
	case model.RSSInfo:
		// RSSProxyInfos
		// http://127.0.0.1:1201/rss?key=巫师财经
		downloadInfo = model.DownloadInfo{
			FolderName: rssInfo.(model.RSSInfo).FolderName,
			RSSUrl: configs.RSSProxyAddress + "/rss?key=" + rssInfo.(model.RSSInfo).RSSInfosName,
			DownloadRoot: rssInfo.(model.RSSInfo).DownloadRoot,
			UseProxy: rssInfo.(model.RSSInfo).UseProxy,
			DownloaderName: rssProxyInfos.DefaultDownloaderName,
		}

	case model.BiliBiliUserInfo:
		// BiliBiliInfos
		// http://192.168.50.135:1200/bilibili/user/video/258150656
		downloadInfo = model.DownloadInfo{
			FolderName: rssInfo.(model.BiliBiliUserInfo).FolderName,
			RSSUrl: configs.RSSHubAddress + "/bilibili/user/video/" + rssInfo.(model.BiliBiliUserInfo).UserID,
			DownloadRoot: rssInfo.(model.BiliBiliUserInfo).DownloadRoot,
			UseProxy: rssInfo.(model.BiliBiliUserInfo).UseProxy,
			DownloaderName: biliBiliInfos.DefaultDownloaderName,
		}
	default:
		return downloadInfo, errors.New("RSS info type not support")
	}

	return downloadInfo, nil
}

func DownloadFromOneFeed(configs model.Configs, rssInfo interface{}) {
	var downloadInfo model.DownloadInfo
	var err error
	downloadInfo, err = SelectDownloadInfo(rssInfo)
	if err != nil {
		log.Errorln("SelectDownloadInfo:", err, rssInfo)
		return
	}
	// 设置代理信息
	downloadInfo.DownloadHttpProxy = configs.DownloadHttpProxy
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

func StartDownload(item *gofeed.Item, downloadInfo model.DownloadInfo) {
	// 还需要拼接具体某人的目录
	var tmpDownloadPath string
	tmpDownloadPath = path.Join(downloadInfo.DownloadRoot, downloadInfo.FolderName)
	// 如果目录不存在则创建
	if Exists(tmpDownloadPath) == false {
		err := os.MkdirAll(tmpDownloadPath, os.ModePerm)
		if err != nil {
			log.Errorln("os.MkdirAll:", tmpDownloadPath, err)
			return
		}
	}
	// 当前路径下面的是否有下载好的文件，或者是否有 .part 的正在下载的文件
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
			if path.Ext(outName) == ConstYoutudlfileextension {
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
	nowDockerDownloader, find := dockerDownloaderInfos[strings.ToLower(downloadInfo.DownloaderName)]
	if find == false {
		log.Errorln("DockerDownloaderInfos", downloadInfo.DownloaderName, "Not Found")
		return
	}
	outstring, err := OneDownload(nowVideoName, item.Link, downloadInfo, nowDockerDownloader)
	if err != nil {
		log.Errorln("OneDownload OutString:", outstring)
		log.Errorln("OneDownload:", err)
		return
	}
	log.Infoln("Download", nowVideoName, "End")
}

// 替换 OutSideAPPOrFolderLocation
func replaceOutSideAPPOrFolderLocation(tmpCommand string, dockerDownloaderInfo model.DockerDownloaderInfo) string {
	for index, one := range dockerDownloaderInfo.OutSideAPPOrFolderLocation {
		nowReplaceKeyWord := strings.ReplaceAll(ConstOutSideAPPOrFolderLocation, "n$", "n" + strconv.Itoa(index) + "$")
		tmpCommand = strings.ReplaceAll(tmpCommand, nowReplaceKeyWord, one)
	}

	return tmpCommand
}

const ConstYoutudlfileextension = ".part"
const ConstPhysicalmachinedownloadrootpath = "$PhysicalMachineDownloadRootPath$"
const ConstHttpadd = "$httpadd$"
const ConstContainerProxy = "$ContainerProxy$"
const ConstDownloadurl = "$downloadURL$"
const ConstNowfilename = "$nowFileName$"
const ConstOutSideAPPOrFolderLocation = "$OutSideAPPOrFolderLocation$"