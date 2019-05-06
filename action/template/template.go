package template

import (
	"bytes"
	"text/template"
	"log"
	"azh/subnginx/model"
)
//配置模板管理，主要是修改
const nginx= `
{{range .}}
upstream {{ println .AppName}}{
	{{range .Url}}server {{ print .}}
	{{end}}check interval=5000 rise=1 fall=3 timeout=4000 type=http
	check_http_send "HEAD / HTTP/1.1\r\nHost: {{range .Url}}{{ print .}},{{end}}\r\n\r\n"
}
{{end}}
server{
	listen	80;
	server_name localhost;
	set $tmpUrl A;
    location / {
{{range .}}
		if ($host ~* ({{ print .AppName}})\.(.*)\.(.*)\.(.*)\.(.*)){
			proxy_pass http://{{ print .AppName}}
			set $tmpUrl '{{range .Url}}{{ print .}},{{end}}'
			break
		}
{{end}}
	}
}
`
//配置模板生成，根据从数据库查出的appconf构建
func Template(guardians []model.AppConf)  string{
	bf:=bytes.NewBuffer([]byte(""))
	masterTmpl, err := template.New("nginx.conf").Parse(nginx)
	if err != nil {
		log.Println(err)
	}
	if err := masterTmpl.Execute(bf,guardians); err != nil {
		log.Println(err)
	}
	bf.Write([]byte("EOF"))
	return bf.String()
}