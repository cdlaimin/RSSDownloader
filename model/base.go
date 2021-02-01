package model

type Test2 struct {
	K int
}

type Configs struct {
	ListenPort        int
	DownloadHttpProxy string
	RSSHubAddress     string
	RSSProxyAddress   string
	EveryTime         string
	ReadRSSTimeOut    int
}

type RSSProxyInfos struct {
	DefaultDownloadRoot string
	DefaultUseProxy bool
	RSSInfos []RSSInfo
}

type RSSInfo struct {
	FolderName string
	RSSInfosName string
	DownloadRoot string
	UseProxy bool
}

type BiliBiliInfos struct {
	DefaultDownloadRoot string
	DefaultUseProxy bool
	BiliBiliUserInfos []BiliBiliUserInfo
}

type BiliBiliUserInfo struct {
	FolderName string
	UserID string
	DownloadRoot string
	UseProxy bool
}

type DownloadInfo struct {
	FolderName string
	RSSUrl string
	DownloadRoot string
	UseProxy bool
}