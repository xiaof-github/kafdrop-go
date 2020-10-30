package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xiaof-github/kafdrop-go/server/models"
)

type KafkaBrokerController struct {
	beego.Controller
}



//返回broker数据
func (c *KafkaBrokerController) BrokerList() {	
	
	dataList, err := models.GetBrokers()	
	if err == nil {
		c.Data["brokerList"] = dataList
	}	
	logs.Info("brokerList :", dataList)
	c.TplName = "kafka/BrokerList.html"

}
