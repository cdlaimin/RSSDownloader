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

### RSSProxy 设置

默认是添加的所有都走代理，所以如果不走代理的就不用在这设置。

RSSHub 的使用请去看对应的官网文档。

### ![RSSProxy-Setting](Pics/RSSProxy-Setting.png)

建议看 [RSSProxy](https://github.com/allanpk716/RSSProxy) 的文档···后续会慢慢补···

这里主要是设置了“巫师财经”的订阅地址，后续会用到。主意这里的 Key 就是 “巫师财经”。

### RSSDownloader 设置

#### 基础配置

![RSSDownloader-Setting](Pics/RSSDownloader-Setting.png)

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

![RSSDownloader-Setting-RSSProxyInfo](Pics/RSSDownloader-Setting-RSSProxyInfo.png)

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

![RSSDownloader-Setting-BiliBiliInfo](Pics/RSSDownloader-Setting-BiliBiliInfo.png)

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

这里的设置用到的 $xxx$ 的变量是 RSSDownloader 使用的内置变量，不要改。

## 项目规划、进度

详细见。Project：[RSSDownloadHub](https://github.com/users/allanpk716/projects/1)

## 感谢

* [youtube-dl](https://github.com/ytdl-org/youtube-dl)
* [youtube-dl docker](https://hub.docker.com/r/qmcgaw/youtube-dl-alpine/)