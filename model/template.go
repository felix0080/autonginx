package model

import (
	"fmt"
	"sync"
	"crypto/sha256"
	"log"
)
/**
	更新nginx的自定义名 			下发失败机器				成功机器		时间				失败审计
	fff,ddd,xxx,sss   			192.1, 12312.123		sajdhaks    12:00			执行过程以及执行失败原因
 */
//can obtain new template

type TemplateLine struct {
	index int
	template []Template
	sync.Mutex
	len int
}
type Template struct {
	TempText string
	HashCode string
}

func NewTemplate(text string) Template {
	temtemp:=Template{
		TempText:text,
	}
	temtemp.Hash()
	return temtemp
}
func (t *Template)Hash()error{
	h := sha256.New()
	if _,err:=h.Write([]byte(t.TempText)); err != nil {
		log.Println(err)
		return err
	}
	t.HashCode=fmt.Sprintf("%x", h.Sum(nil))
	return nil
}
//---- puts a new template in templateline
func (t *TemplateLine) Update(temp Template)  {
	t.Lock()
	defer t.Unlock()
	index := (t.index + 1) % (t .len - 1 )
	t.template[index] = temp
}

//---- gets  templates history from templateline
func (t *TemplateLine)GetsHistory(temp Template) []Template {
	t.Lock()
	defer t.Unlock()
	tmptemplate:=make([]Template,t.len,t.len)
	tempindex:=0
	for i:=t.index;i >= 0 ; i-- {
		tmptemplate[tempindex]=t.template[i]
		tempindex++
	}
	for i:=t.len ; i > t.index ; i-- {
		tmptemplate[tempindex]=t.template[i]
		tempindex++
	}
	return tmptemplate
}
//0 is new
func (t *TemplateLine) GetRecent(sub int)Template{
	var is int
	if ts := t.index - sub ; ts < 0 {
		is = ts % t.len + t.len
	}else {
		is = ts % t.len
	}
	return t.template[is]
}