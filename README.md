# RSSDownloader

## Why？

初始的目标是代替 [BiliBiliDownloader](https://github.com/allanpk716/BiliBiliDownloader) 这个项目，之前做的很随意，效率低，但是也足够用。后因为[巫师财经](https://www.youtube.com/channel/UC55ahPQ7m5iJdVWcOfmuE6g)跟B站闹掰了，导致本来一个脚本也能搞定的时候，现在得跑 Youtube 去专门下载。既然如此，那么就打算把订阅的视频下载功能给重构，顺带练习下 golang 的基本使用。

## 目的

为了把关注的博主的视频收集下载，方便在家离线观看，同时也为娃提前构建知识库。

后续的重构，会把 RSS 订阅的 BT 以及图片的下载都支持，目标是做到家庭内部的订阅的 All in One 下载。

## 特性

目前示例上提供了两个站点的支持：

* Youtube
* BiliBili

如果理解了整套的逻辑，其实支不支持下载，就是 RSS 订阅获取到的是什么网站的视频页面地址，然后 [Youtube-dl](https://github.com/ytdl-org/youtube-dl) 是否支持下载的问题了。同理，其实可以扩展到你实现的下载器 docker 与对应相应网站的 RSS 订阅下载。

## 程序结构、组成

本程序需要依赖一下几个部分：

* [RSSHub](https://rsshub.app)（官方的地址）
* RSSHub（SelfHost 自己搭建的）
* [RSSProxy](https://github.com/allanpk716/RSSProxy)（解决需要 Proxy 才能使用的 RSS 源）
* [RSSDownloader](https://github.com/allanpk716/RSSDownloader)（本程序）

![all-struct](Pics/all-struct.png)

### 为什么会出现 RSSProxy 以及两个 RSSHub？



## 如何使用

## 项目规划、进度

详细见。Project：[RSSDownloadHub](https://github.com/users/allanpk716/projects/1)

## 感谢