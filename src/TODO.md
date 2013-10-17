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

满足不了Save规则的，但是满足Fetch规则，需要在内存里面解析并记录url，只是不持久化

GOPA集群化，每个gopa只设置一个cluster参数和node参数[可选]，
通过集群web面板来管理任务
每个gopa可以分别设置角色：fetch、parse、master
gopa也分shard【考虑一致性hash】

目录太大，自动shard，切分，需要统计目录文件大小

分页参数，自动保存到文件，文件名自动重命名,broken_by_parameter

各bloomfilter关闭时持久化


页面保存的时候，丢失了当前页面的地址，如果页面的url路径是相对路径，则匹配会失败，需要修复页面的相对路径为绝对路径


职责单一化，下载的只负责下载，可分启动

url可以保存到本地文件，一行一个