package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type KafkaBrokerController struct {
	beego.Controller
}



//返回broker数据
func (c *KafkaBrokerController) List() {
	
	
	c.Data["List"] = "Hello world"	
	logs.Info("dataList :")
	c.TplName = "kafka/BrokerList.html"

}
