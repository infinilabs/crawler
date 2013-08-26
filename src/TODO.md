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

检查内存泄露的原因

taskItem任务未接收到新的，没有进行下载操作

配置文件支持多个参数，通过，分割，转换成fields

任务
    每个gopa能够注册一系列partition，具体是否执行partition下的任务，由集群来分配

Random U-A
Random Refer

Parsed 和 Download的 BloomFilter 分开
Offset文件放项目文件夹里面