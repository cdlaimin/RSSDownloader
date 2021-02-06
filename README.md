- [RSSDownloader](#rssdownloader)
  * [Why？](#why-)
  * [目的](#--)
  * [特性](#--)
  * [程序设计思路](#------)
    + [为什么会出现 RSSProxy 以及两个 RSSHub？](#-------rssproxy------rsshub-)
    + [为何下载部分要分开使用额外的 docker 容器？](#---------------docker----)
  * [如何使用](#----)
    + [RSSProxy 设置](#rssproxy---)
    + [RSSDownloader 设置](#rssdownloader---)
      - [基础配置](#----)
        * [ListenPort](#listenport)
        * [DownloadHttpProxy](#downloadhttpproxy)
        * [RSSHubAddress](#rsshubaddress)
        * [RSSProxyAddress](#rssproxyaddress)
        * [ReadRSSTimeOut](#readrsstimeout)
        * [EveryTime](#everytime)
      - [RSSProxyInfos 设置](#rssproxyinfos---)
        * [DefaultDownloaderName](#defaultdownloadername)
        * [DefaultDownloadRoot](#defaultdownloadroot)
        * [DefaultUseProxy](#defaultuseproxy)
        * [RSSInfos](#rssinfos)
      - [BiliBiliInfos 设置](#bilibiliinfos---)
        * [DefaultDownloaderName](#defaultdownloadername-1)
        * [DefaultDownloadRoot](#defaultdownloadroot-1)
        * [DefaultUseProxy](#defaultuseproxy-1)
        * [BiliBiliUserInfos](#bilibiliuserinfos)
      - [DockerDownloaderInfos 设置](#dockerdownloaderinfos---)
        * [ContainerProxy](#containerproxy)
        * [PhysicalMachineDownloadRootPath](#physicalmachinedownloadrootpath)
        * [OutSideAPPOrFolderLocation](#outsideapporfolderlocation)
        * [httpadd](#httpadd)
        * [downloadURL](#downloadurl)
        * [nowFileName](#nowfilename)
    + [Docker 部署](#docker---)
      - [1. 部署用的 [docker-compose.yaml](https://github.com/allanpk716/RSSDownloader/blob/master/DockerThings/docker-compose.yaml)](#1-------docker-composeyaml--https---githubcom-allanpk716-rssdownloader-blob-master-dockerthings-docker-composeyaml-)
      - [2. RSSProxy 的 [Config.yaml](https://github.com/allanpk716/RSSDownloader/blob/master/DockerThings/rssproxy-config.yaml)](#2-rssproxy----configyaml--https---githubcom-allanpk716-rssdownloader-blob-master-dockerthings-rssproxy-configyaml-)
      - [3. RSSDownloader 的 [Config.yaml](https://github.com/allanpk716/RSSDownloader/blob/master/DockerThings/rssdownloader-config.yaml)](#3-rssdownloader----configyaml--https---githubcom-allanpk716-rssdownloader-blob-master-dockerthings-rssdownloader-configyaml-)
      - [4.下载 [youtube-dl](https://github.com/ytdl-org/youtube-dl/releases)](#4----youtube-dl--https---githubcom-ytdl-org-youtube-dl-releases-)
  * [风险](#--)
  * [可能遇到的问题](#-------)
    + [1. 部分模块失效导致无法正常使用](#1---------------)
    + [2. 代理问题](#2-----)
  * [项目规划、进度](#-------)
  * [更新记录](#----)
  * [感谢](#--)

<small><i><a href='http://ecotrust-canada.github.io/markdown-toc/'>Table of contents generated with markdown-toc</a></i></small>

# RSSDownloader

## Why？

初始的目标是代替 [BiliBiliDownloader](https://github.com/allanpk716/BiliBiliDownloader) 这个项目，之前做的很随意，效率低，但是也足够用。后因为[巫师财经](https://www.youtube.com/channel/UC55ahPQ7m5iJdVWcOfmuE6g)跟B站闹掰了，导致本来一个脚本也能搞定的时候，现在得跑 Youtube 去专门下载。既然如此，那么就打算把订阅的视频下载功能给重构，顺带练习下 golang 的基本使用。

## 目的

为了把关注的博主的视频收集下载，方便在家离线观看，同时也为娃提前构建知识库。

后续的重构，会把 RSS 订阅的 BT 以及图片的下载都支持，目标是做到家庭内部的订阅的 All in One 下载。

## 特性

**注意**！现在还是初期开发阶段，虽然已经可以基本正常的使用了，但是怎么便于一般人使用还有待磨合，如果不是很急，可以观望以下。后续会提供 docker-compose 的全家桶版本，便于各位使用。

目前示例上提供了两个站点的支持：

* Youtube
* BiliBili

目前仅仅支持 docker 部署，Windows 暂时没有列入支持计划。

如果理解了整套的逻辑，其实支不支持下载，就是 RSS 订阅获取到的是什么网站的视频页面地址，然后 [Youtube-dl](https://github.com/ytdl-org/youtube-dl) 是否支持下载的问题了。同理，其实可以扩展到你实现的下载器 docker 与对应相应网站的 RSS 订阅下载。

比如后续其实可以很容易扩展，pixiv 等网站的图片下载，又或者某些视频的 RSS 订阅 BT 下载。

## 程序设计思路

本程序需要依赖一下几个部分：

* [RSSHub](https://rsshub.app)（官方的地址）
* RSSHub（SelfHost 自己搭建的）
* [RSSProxy](https://github.com/allanpk716/RSSProxy)（解决需要 Proxy 才能使用的 RSS 源）
* [RSSDownloader](https://github.com/allanpk716/RSSDownloader)（本程序）

![all-struct](Pics/all-struct.png)

### 为什么会出现 RSSProxy 以及两个 RSSHub？

对于 RSSProxy，起因还是因为 TTRSS 需要订阅一些信息，但是 TTRSS 的代理设置遇到问题，那么就做了一个中转的程序进行过度。主要还是依赖 RSSHub 进行订阅，但是去年 RSSHub 已经无法正常访问了，需要代理，所以···其中如果是 Instagram 或者 tweet ，除去需要账号密码，还有就是 RSSHub 默认能够读取的条目数是有限的，那么比如想从某个博主的第一条开始制作 RSS 就会遇到问题，也就没法很好的让 TTRSS 同步所有的信息。所以其实 RSSProxy 附带的做了  Instagram 和 tweet 的订阅 RSS 功能。

对于两个 RSSHub，还是因为有一些网站的 RSS 订阅，需要使用自己的账号密码，或者就是用官方 RSSHub 也会在部分订阅的 RSS 建议做 SelfHost 。所以就存在了 SelfHost 的 RSSHub。

### 为何下载部分要分开使用额外的 docker 容器？

主要还是主要功能的分离，这样就主体逻辑部分一般来说容易定型，不怎么需要经常改动，但是 RSS 订阅和下载器这类，就会很可能需要经常跟随着目标网站来更新。（自己维护不现实）

## 如何使用

下面会按本人在家部署的订阅 RSS 下载服务来举例。

目前有以下几个视频博主想要保存他们的视频进行存档：

* 巫师财经
* 李永乐
* 回形针PaperClip
* 柴知道
* 吟游诗人基德
* 讲解员河森堡
* 沙盘上的战争

每一个博主的视频希望能够按发布的时间加标题进行存储。比如：

```
Auther/2021-02-05_Title.mp4
```

除了[巫师财经](https://www.youtube.com/channel/UC55ahPQ7m5iJdVWcOfmuE6g)外，他们上传视频的网站是 B 站。那么就需要分两个方向去下载视频：

* Youtube 类型的需要代理下载，走 RSSProxy 代理
* BiliBili 无需代理下载，直接用 RSSHub 订阅

如果想快速尝鲜，那么就直接跳到 [Docker 部署](###Docker 部署)

### RSSProxy 设置

默认是添加的所有都走代理，所以如果不走代理的就不用在这设置。

RSSHub 的使用请去看对应的官网文档。

### <img src="Pics/RSSProxy-Setting.png" alt="RSSProxy-Setting" style="zoom:50%;" />

建议看 [RSSProxy](https://github.com/allanpk716/RSSProxy) 的文档···后续会慢慢补···

这里主要是设置了“巫师财经”的订阅地址，后续会用到。主意这里的 Key 就是 “巫师财经”。

### RSSDownloader 设置

#### 基础配置

基本的配置信息如下图：

<img src="Pics/RSSDownloader-Setting.png" alt="RSSDownloader-Setting" style="zoom:50%;" />

##### ListenPort

本程序的监听端口，如果是docker部署不要改。

##### DownloadHttpProxy

下载使用的代理服务器地址

##### RSSHubAddress

SelfHost 的 RSSHub 地址

##### RSSProxyAddress

RSSProxy 的地址

##### ReadRSSTimeOut

读取 RSS 的超时设置

##### EveryTime

下载轮询的间隔

#### RSSProxyInfos 设置

这个需要配合着  [RSSProxy 设置](###RSSProxy 设置)  的设置信息来设置。

<img src="Pics/RSSDownloader-Setting-RSSProxyInfo.png" alt="RSSDownloader-Setting-RSSProxyInfo" style="zoom:50%;" />

##### DefaultDownloaderName

默认使用的下载器名称，比如 YouTube-dl，但是着其实是你下面设置的下载器来决定的哈。

##### DefaultDownloadRoot

默认的下载路径，注意，这个是你物理机器的路径，因为是要传递给 docker 下载器使用的。

##### DefaultUseProxy

是否使用代理，当然也可以在下面每个具体的 RSS key 中再次设置

##### RSSInfos

这需要配合 RSSProxy 来设置如上图

巫师财经：

* RSSInfosName：这里一定要跟 RSSProxy 中 RSSInfos 设置的 Key 一致。
* DownloadRoot：不填写这个字段，就是使用 RSSProxyInfos 的设置。这里可以再次指定这个 Key  单独的下载位置，注意一定是物理机的路径。
* UseProxy：不填写这个字段，就是使用 RSSProxyInfos 的设置。这里可以再次指定这个 Key  是否使用代理。

#### BiliBiliInfos 设置

这里直接走的是自建的 RSSHub，默认无需走代理。如果你想走代理，那么就 DefaultUseProxy: true

<img src="Pics/RSSDownloader-Setting-BiliBiliInfo.png" alt="RSSDownloader-Setting-BiliBiliInfo" style="zoom:50%;" />

##### DefaultDownloaderName

默认使用的下载器名称，比如 YouTube-dl，但是着其实是你下面设置的下载器来决定的哈。

##### DefaultDownloadRoot

默认的下载路径，注意，这个是你物理机器的路径，因为是要传递给 docker 下载器使用的。

##### DefaultUseProxy

是否使用代理，当然也可以在下面每个具体的 RSS key 中再次设置

##### BiliBiliUserInfos

李永乐:

* UserID: 这个 BiliBili 对应用户的 ID
* DownloadRoot：不填写这个字段，就是使用 RSSProxyInfos 的设置。这里可以再次指定这个 Key  单独的下载位置，注意一定是物理机的路径。
* UseProxy：不填写这个字段，就是使用 RSSProxyInfos 的设置。这里可以再次指定这个 Key  是否使用代理。

#### DockerDownloaderInfos 设置

这里指定了可供选择的 Docker 下载器

![RSSDownloader-Setting-DockerDownloaderInfos](Pics/RSSDownloader-Setting-DockerDownloaderInfos.png)

默认给出了一个示例的下载器 [config.yaml.sample](https://github.com/allanpk716/RSSDownloader/blob/master/config.yaml.sample) ，使用的是 [youtube-dl docker](https://hub.docker.com/r/qmcgaw/youtube-dl-alpine/)。

这里的 Key 是 Youtube-dl，上面的设置使用默认下载器的时候用到了。

* DockSSHAddress：因为使用的是 SSH 去启动物理机的 docker 的，所以需要设置物理机的 SSH 连接

* DockerUserName：有 docker 权限的账户，且你的下载地址也得有对应的权限

* DockerPassword：密码

* UpdateCommands：连接 SSH 执行的更新用的命令，为了更新你的下载器

* DownloadCommands：下载器 docker 使用的命令

* OutSideAPPOrFolderLocation：这里可以动态的指定多个需要外部映射可执行程序或者是文件夹路径。

这里的设置用到的 $xxx$ 的变量是 RSSDownloader 使用的内置变量，不要改。

注意，这里设置 docker 的 name 为 youtube-dl-runner，是为了，如果本程序故障的时候，重复启动的时候，如果发现上一个 docker 下载器启动了没有推出就不重复启动了，不然会出问题。

**有几个特殊的内置字段需要注意下**：

##### ContainerProxy

这个与 DownloadHttpProxy 一致，会被替换，一般是为了 docker 内程序更新使用的，因为国内嘛，你懂。

##### PhysicalMachineDownloadRootPath

跟相应资源的 DefaultDownloadRoot 一致，会自动替换。

##### OutSideAPPOrFolderLocation

这个可以设置多个值，然后替换的时候是以每一个 command 中，指定对应的顺序。举例：

设置了：

```
OutSideAPPOrFolderLocation:
	- /mnt/user/appdata/rssdownloader/abc
	- /mnt/user/appdata/rssdownloader/haha
```

然后执行的命令有：

```
- docker run --rm -v $OutSideAPPOrFolderLocation0$:/downloads -v $OutSideAPPOrFolderLocation1$:/usr/local/bin/youtube-dl qmcgaw/youtube-dl-alpine
```

那么对应的：

* $OutSideAPPOrFolderLocation0$ = /mnt/user/appdata/rssdownloader/abc
* $OutSideAPPOrFolderLocation1$ = /mnt/user/appdata/rssdownloader/haha

如果执行的命令是：

```
- docker run --rm -v $OutSideAPPOrFolderLocation1$:/downloads -v $OutSideAPPOrFolderLocation0$:/usr/local/bin/youtube-dl qmcgaw/youtube-dl-alpine
```

那么其实对应的值也**一样**的，只不过在命令中换了个位置。

##### httpadd

这个是 youtube-dl 的代理设置，这个与 DownloadHttpProxy 一致，会被替换，一般是为了 docker 内程序更新使用的，因为国内嘛，你懂。

##### downloadURL

需要下载的资源的 URL，这个是 RSS 中解析出来的，无需改动。

##### nowFileName

需要将下载的资源重命名为什么名称，不包含后缀名，因为会有对应的下载程序决定。

这里默认是 **2021-02-06_VideoTitle** 这样的格式。无需修改。

注意，这里一定要用这个命名的格式，不然前面的检测是否下载过的逻辑会出问题，虽然问题不大，正常会有下载程序判断是否下载过了会跳过。

### Docker 部署

> 如果你一开始就跳到这里来看了，那么理想情况，你根据下面的提示修改基本的信息是能够直接跑起来的。如果你想知道这些设置参数有啥子用，那么建议你把上面的如何使用给看了，看不懂的话，一定是我描述的问题，不是你的问题，希望提 [ISSUS](https://github.com/allanpk716/RSSDownloader/issues) 帮后续人的能看懂（逃。

可以参考 RSSDownloader 项目中 [DockerThings](https://github.com/allanpk716/RSSDownloader/tree/master/DockerThings) 中的几个文件。

下面给出的几个文件都是以物理机 IP <u>192.168.50.135</u> 举例。

#### 1. 部署用的 [docker-compose.yaml](https://github.com/allanpk716/RSSDownloader/blob/master/DockerThings/docker-compose.yaml)

这里的 RSSHub 其实就是自行部署的 Selfhost 。如果你已经部署有了，那么就无需再整一个出来。如果不想自己部署一个，那么就用官方的 RSSHub，效果嘛，自行判断。

```yaml
version: "3"
services:
  service_rssdownloader:
    image: allanpk716/rssdownloader:latest
    container_name: rssdownloader
    ports:
      - 1202:1200
    volumes:
      - /mnt/user/appdata/rssdownloader/config.yaml:/app/config.yaml
  service_rssproxy:
    image: allanpk716/rssproxy:latest
    container_name: rssproxy
    ports:
      - 1201:1200
    volumes:
      - /mnt/user/appdata/rssproxy/config.yaml:/app/config.yaml
  service_rsshub:
    image: diygod/rsshub:latest
    container_name: rsshub
    ports:
      - 1200:1200
    environment:
      - TZ=Asia/Shanghai
```

#### 2. RSSProxy 的 [Config.yaml](https://github.com/allanpk716/RSSDownloader/blob/master/DockerThings/rssproxy-config.yaml)

```yaml
ListenPort: 1200
HttpProxy: http://192.168.50.252:20171
EveryTime: 4h

RSSInfos:
  巫师财经: https://rsshub.app/youtube/channel/UC55ahPQ7m5iJdVWcOfmuE6g
```

#### 3. RSSDownloader 的 [Config.yaml](https://github.com/allanpk716/RSSDownloader/blob/master/DockerThings/rssdownloader-config.yaml)

```yaml
ListenPort: 1200
DownloadHttpProxy: http://192.168.50.252:20171
RSSHubAddress: http://192.168.50.135:1200
RSSProxyAddress: http://192.168.50.135:1201
ReadRSSTimeOut: 30
EveryTime: 4h

# 这里订阅的RSS其实是 RSSProxy 中转的，这些 RSS 需要代理才能访问
RSSProxyInfos:
  DefaultDownloaderName: Youtube-dl
  DefaultDownloadRoot: /mnt/remotes/Video/科普
  # 使用代理下载
  DefaultUseProxy: true
  RSSInfos:
    巫师财经:
      # 相应 RSSProxy中 RSSInfos 的 key
      RSSInfosName: 巫师财经

# 这里直接走的是自建的 RSSHub，默认无需走代理
BiliBiliInfos:
  DefaultDownloaderName: Youtube-dl
  DefaultDownloadRoot: /mnt/remotes/Video/科普
  # 使用代理下载
  DefaultUseProxy: false
  BiliBiliUserInfos:
    # 这里的 Key 会直接作为下载文件夹
    李永乐:
      UserID: 9458053
    回形针PaperClip:
      UserID: 258150656
    柴知道:
      UserID: 26798384
    吟游诗人基德:
      UserID: 510856133
    讲解员河森堡:
      UserID: 483884702
    沙盘上的战争:
      UserID: 612194373

DockerDownloaderInfos:
  Youtube-dl:
    DockSSHAddress: 192.168.50.135:22
    DockerUserName: dockeruser
    DockerPassword: password
    OutSideAPPOrFolderLocation:
      - /mnt/user/appdata/rssdownloader/youtube-dl
    UpdateCommands:
      - docker pull qmcgaw/youtube-dl-alpine
      - docker run --rm --name=youtube-dl-runner -e LOG=no -e AUTOUPDATE=yes $ContainerProxy$ -v $PhysicalMachineDownloadRootPath$:/downloads -v $OutSideAPPOrFolderLocation0$:/usr/local/bin/youtube-dl qmcgaw/youtube-dl-alpine
    DownloadCommands:
      - docker run --rm --name=youtube-dl-runner -e LOG=no -v $PhysicalMachineDownloadRootPath$:/downloads -v $OutSideAPPOrFolderLocation0$:/usr/local/bin/youtube-dl qmcgaw/youtube-dl-alpine --proxy=$httpadd$ $downloadURL$ -o "/downloads/$nowFileName$.%(ext)s"
```

以上是默认的配置，下图会标记出你需要改的地方。这个就是根据你的物理机 **IP** 以及**存储路径**来调整了。

![RSSDownloader-Setting-modify](Pics/RSSDownloader-Setting-modify.png)

#### 4.下载 [youtube-dl](https://github.com/ytdl-org/youtube-dl/releases)

[youtube-dl](https://github.com/ytdl-org/youtube-dl/releases) 这个文件需要自己也下载好，不然 youtube-dl-docker 启动后会提示找不到文件的。放到的目录需要与你设置的 RSSDownloader -- config.yaml -- DockerDownloaderInfos -- Youtube-dl -- OutSideAPPOrFolderLocation

```yaml
   - /mnt/user/appdata/rssdownloader/youtube-dl
```

一定要注意，这个位置放对，也写对。后续自动更新会根据这里来更改的。

<img src="Pics/youtube-dl-download.png" alt="youtube-dl-download" style="zoom:50%;" />

## 风险

从项目的结构来看，是使用了 docker 去启动 docker 的方式来做到的，所以用到了 ssh 去操作你的 docker，风险你懂的，一定是在内网玩，别暴露到外网哈。

## 可能遇到的问题

### 1. 部分模块失效导致无法正常使用

本系统引用了基础外部的库，如果他们针对相应的网站更新了，那么一般就无法正确的获取 RSS 信息，或者正确去下载了。如果发现了，希望能提 [ISSUS ](https://github.com/allanpk716/RSSDownloader/issues)。

### 2. 代理问题

本项目主要目标是在国内使用，所以代理是必须的。如果你的代理不稳定，这个项目是没法解决的···暂时不想整公用代理池，因为这个你可以自行解决。

## 项目规划、进度

详细见。Project：[RSSDownloadHub](https://github.com/users/allanpk716/projects/1)

不排除如果没有新的个人需求，不会进行重构（主要是懒，写文档很累···）

## 更新记录

> 本来是设置了 docker hub 的自动根据 tag 制作镜像的，但是···配置没动过，第三次 tag push 上去发现没得反应，只能整一个手动根据 master 分支输出 latest ···
>
> 所以对应的 tag 没用上，直接用 latest 来获取吧···

* 2021年2月6日 -- tag:v0.0.3 -- 支持 docker-compose 版本

## 感谢

* [youtube-dl](https://github.com/ytdl-org/youtube-dl)
* [youtube-dl docker](https://hub.docker.com/r/qmcgaw/youtube-dl-alpine/)
* [instafeed](https://github.com/falzm/instafeed)
* [goinsta](https://github.com/ahmdrz/goinsta)
* [twitter2rss](https://github.com/n0madic/twitter2rss)