# Gopa # Check List
[狗爬],A Spider Written in GO.

核心
    持久化抓取队列，恢复抓取队列

存储
    使用weedfs存储
    domain or url hash，sharding

内容处理
    文字block抓取
    处理跳转：<meta http-equiv="refresh" content="0;url=http://www.baidu.com/">