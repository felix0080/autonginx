package main

import (
	"time"
	"fmt"
	"azh/monitor/server"
	"log"
)

/*


实现邮箱发送容器，对接行内接口，和

每个实例，以及对应的用户
1. 按天，GPU资源的平均使用率（30%），最高使用率；
2. 使用率为0的时间（需要时间段，用来建议用户在非使用时间停止应用）。
3. 一直空跑

邮件通知内容：
1. 建议用户降低资源分配
2. 建议用户在非使用时间停止应用
3. 建议用户停止应用，释放资源。
建议用户把中间结果存在云硬盘。

仅仅统计每一天的
没到一天的从开始运行时间计算到当前时间停止

只关注
*/

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	server.HttpServer()
}