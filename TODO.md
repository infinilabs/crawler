# Gopa # Check List
[狗爬],A Spider Written in GO.

核心

存储
    使用weedfs存储
    domain or url hash，sharding,为避免bloom过大，优先处理部分domain，剩下的都入domain命名的队列，处理完一个domain，然后卸载bloom，加载其它的bloom，分别处理
    未保存文件，优先解析url，本地已存储文件，入解析url队列，使用本地路径作为url
    使用一个channel，处理多个事件，实现fetch及parser gorouting的优雅关闭      [done]


内容处理
    文字block抓取
    处理跳转：<meta http-equiv="refresh" content="0;url=http://www.baidu.com/">

检查内存泄露的原因       [done]

taskItem任务未接收到新的，没有进行下载操作

配置文件支持多个参数，通过，分割，转换成fields

任务
    每个gopa能够注册一系列partition，具体是否执行partition下的任务，由集群来分配

Random U-A
Random Refer

Parsed 和 Download的 BloomFilter 分开[done]
Offset文件放项目文件夹里面[done]

满足不了Save规则的，但是满足Fetch规则，需要在内存里面解析并记录url，只是不持久化

GOPA集群化，每个gopa只设置一个cluster参数和node参数[可选]，
通过集群web面板来管理任务，seed参数非必须，通过web来添加
每个gopa可以分别设置角色：fetch、parse、master
gopa也分shard【考虑一致性hash】

目录太大，自动shard，切分，需要统计目录文件大小

分页参数，自动保存到文件，文件名自动重命名,broken_by_parameter   [done]

各bloomfilter关闭时持久化   [done]

页面保存的时候，丢失了当前页面的地址，如果页面的url路径是相对路径，则匹配会失败，需要修复页面的相对路径为绝对路径   [done]

职责单一化，下载的只负责下载，可分别启动

url可以保存到本地文件，一行一个,每个节点预先分配一个shard段，只处理本段的url，其它段的url，集群自动同步

shard下载队列是主动获取，最外面的master分配任务的时候，只有当前workers有空闲的时候，才分配任务

检测本地是否存在，如果存在则不处理，并添加到bloomfilter  [done]

根据url参数模板来批量下载网页 [done]

Cookie [done]

速度控制,阀值控制  [done]
    超时
    返回错误页面
    自动控制速度，暂停，自动调整合适的速度

[1 任务定制(URL获取)]->[2 任务解析]->[3 任务预览(URL查看)]->[4 任务执行]->[5 保存文件]
  
[1.1 自动爬取]
[1.2 精确爬取]
    [1.2.1 变量控制]->[变量控制... 管道式控制]
    
[5.1 保存前处理]
    [文件名] [文件名重命名] [文件名移除] [文件参数处理] [文件名]
[5.2 保存]
[5.3 保存后处理]
    

每个任务都有多个预判，如抓取，预判：主机、端口、后缀等，每个预判多种表达方式：符合、不符合、包含、正则、脚本（返回BOOL）

任务按域名分组，每个域名一个协程，同一域名下面的请求走同一个协程，可以配置每个域名的并发协程和多少个域名并发抓取，所有协程的总算可配

使用Hyperloglog来去重：1.count；2.ADD；3.count；3.比较count；计算换内存（与bloomfilter相比）

https://github.com/dgryski/go-minhash.git

文件是否抽取链接也需要进行多种方式判断

sharding、本地持久化、避免网络拷贝，默认不网络转发、如果其他节点该分片没有任务了则网络迁移过去

保存前：404判断：1.状态码；2.标题：<title>404 Not Found</title>；3.其它自定义条件

进度监控：parse、fetch


COMMAND chan, 任务都进COMMAND 协程，然后统一分发处理