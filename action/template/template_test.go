package template

import (
	"testing"
	"github.com/qiniu/log"
	"azh/subnginx/model"
)

func TestTemplate(t *testing.T) {
	guardians := []model.AppConf{{
		AppName:"mServer",
		Url:[]string{"2048.demo1.sz-qh-dev1-gpu01.aipower.cmbchina.cn",
			"2048.demo2.sz-qh-dev1-gpu01.aipower.cmbchina.cn"},
	},{
		AppName:"mServer2",
		Url:[]string{"test-2048.demo1.sz-qh-dev1-gpu01.aipower.cmbchina.cn",
			"test-2048.demo2.sz-qh-dev1-gpu01.aipower.cmbchina.cn"},
	}}
	log.Println(Template(guardians))
}
