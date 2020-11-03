package routers

import (
	"github.com/astaxie/beego"
	"github.com/xiaof-github/kafdrop-go/controllers"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	userInfoController := &controllers.UserInfoController{}
	beego.Router("/admin/userInfo/edit", userInfoController, "get:Edit")
	beego.Router("/admin/userInfo/delete", userInfoController, "post:Delete")
	beego.Router("/admin/userInfo/save", userInfoController, "post:Save")
	beego.Router("/admin/userInfo/list", userInfoController, "get:List")

	kafkaBrokerController := &controllers.KafkaBrokerController{}
	beego.Router("/admin/kafka/brokerList", kafkaBrokerController, "get:BrokerList")

	kafkaTopicController := &controllers.KafkaTopicController{}
	beego.Router("/admin/kafka/topicList", kafkaTopicController, "get:TopicList")
	beego.Router("/admin/kafka/topicMessage", kafkaTopicController, "get:TopicMessage")
}
