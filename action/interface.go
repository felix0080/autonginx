package action

import (
	"azh/subnginx/db"
	"azh/subnginx/action/template"
	"azh/subnginx/action/command"
	"azh/subnginx/config"
	"fmt"
	"log"
)

func Run()  {
	//保存上一次的template数据，因为可能会回滚
	//拿到配置中的数据
	righthost,_:=config.GetConfig().IpAddress.ScanPort()

	appconf:=db.SelectAll()

	//先发送探针包 -》config.GetConfig().IpAddress
	//需要配置端口探测的具体端口

	//未探测成功的机器将随日志一起打出去，发送给web hook
	//生成模板
	nginxConf:=template.Template(appconf)
	nginxConf=fmt.Sprintf(`/usr/bin/cat << EOF > %s
%s`,config.GetConfig().NginxConfigPath,nginxConf)
	//改为探测成功的机器
	for _,ip:=range righthost{
		callback:=command.SshTo(ip,nginxConf,config.GetConfig().RetryTimesWhenFail)
		if callback.Err!=nil {
			log.Println(callback.Err)
			//调用web hook推送消息

			//根据配置决策是否要回滚还是继续部署下一台，回滚只对前面以前的ip进行回滚（已操作的）

		}
	}


	db.SelectAll()
	//对探针成功的机器请求建立连接，下发文件或者命令

	//记录下发成功的机器，下发成功的记录
}