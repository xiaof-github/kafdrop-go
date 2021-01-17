package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xiaof-github/kafdrop-go/models"
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

// 返回消息数据
func (c *KafkaTopicController) TopicMessage() {

	//获得topic
	topic := c.GetString("topic")
	logs.Info("topic: ", topic)
	data, err := models.GetTopicMessages(topic)	
	if err == nil {
		c.Data["topicMessage"] = data
	}
	c.TplName = "kafka/TopicMessage.html"
}