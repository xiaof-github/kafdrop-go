package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type KafkaController struct {
	beego.Controller
}



//返回全部数据
func (c *KafkaController) List() {
	
	
	c.Data["List"] = "Hello world"	
	logs.Info("dataList :")
	c.TplName = "kafka/TopicList.html"

}
