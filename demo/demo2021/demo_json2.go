
package main
 
import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
)
 
func main() {
       //拼凑json   body为map数组
	var rbody []map[string]interface{}
	t := make(map[string]interface{})
	t["device_id"] = "dddddd"
	t["device_hid"] = "ddddddd"
	t["device_cid"] = nil
 
	rbody = append(rbody, t)
	t1 := make(map[string]interface{})
	t1["device_id"] = "aaaaa"
	t1["device_hid"] = "aaaaa"
 
	rbody = append(rbody, t1)
 
	cnnJson := make(map[string]interface{})
	cnnJson["code"] = 0
	cnnJson["request_id"] = 123
	cnnJson["code_msg"] = ""
	cnnJson["body"] = rbody
	cnnJson["page"] = 0
	cnnJson["page_size"] = 0
 
	b, _ := json.Marshal(cnnJson)
	cnnn := string(b)
	fmt.Println("cnnn:%s", cnnn)
	cn_json, _ := simplejson.NewJson([]byte(cnnn))
	cn_body, _ := cn_json.Get("body").Array()
 
	for _, di := range cn_body {
                //就在这里对di进行类型判断
		newdi, _ := di.(map[string]interface{})
		device_id := newdi["device_id"]
		device_hid := newdi["device_hid"]
        device_cid := newdi["device_cid"]
		fmt.Println(device_hid, device_id, device_cid)
	}
 
}
