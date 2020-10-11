package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type KafkaTopicController struct {
	beego.Controller
}



//返回Topic数据
func (c *KafkaTopicController) TopicList() {
	
	
	c.Data["List"] = "Hello world"	
	logs.Info("dataList :")
	c.TplName = "kafka/TopicList.html"

}

