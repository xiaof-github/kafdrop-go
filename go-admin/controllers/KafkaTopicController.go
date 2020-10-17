package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/golangpkg/go-admin/models"
)

type KafkaTopicController struct {
	beego.Controller
}



//返回Topic数据
func (c *KafkaTopicController) TopicList() {
	
	
	dataList, err := models.GetTopics()	
	if err == nil {
		c.Data["topicList"] = dataList
	}	
	logs.Info("topicList :", dataList)	
	c.TplName = "kafka/TopicList.html"

}

