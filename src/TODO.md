# Gopa # Check List
[狗爬],A Spider Written in GO.

核心

存储
    使用weedfs存储
    domain or url hash，sharding
    hash to kafka channel
    未保存文件，优先解析url，本地已存储文件，入解析url队列，使用本地路径作为url
    使用一个channel，处理多个事件，实现kafka及urlextract gorouting的优雅关闭


内容处理
    文字block抓取
    处理跳转：<meta http-equiv="refresh" content="0;url=http://www.baidu.com/">