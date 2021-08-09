

package main
 
import (
	"encoding/json"
	"fmt"
    "io/ioutil"

    log "github.com/sirupsen/logrus"
	// "github.com/bitly/go-simplejson"
)

type SoftwareConf struct {
    Id           string `json:"id"`
    ProcessName  string `json:"processName"`
    SoftwareName string `json:"softwareName"`
    OrgCode      string `json:"orgCode"`
    OrgName      string `json:"orgName"`
    OrgType      string `json:"orgType"`
    CityCode     string `json:"cityCode"`
    CityName     string `json:"cityName"`
    DistrictCode string `json:"districtCode"`
    DistrictName string `json:"districtName"`
    VendorName   string `json:"vendorName"`
    Status       int    `json:"status"`
}
 
func main() {

    var conf []SoftwareConf
    var processConf []string
    var ConfName string = "software.json"

    err := ReadConfFromJson(ConfName, &conf)
    if err != nil {
        log.Error(err)    
    } else {
        log.Info("req c")
        log.Info(conf)
        processConf = make([]string, 0)
        for _, item := range conf {
            if item.Status == 1 {
                processConf = append(processConf, item.ProcessName)
            }
        }
        log.Info(processConf)
    }
}

func ReadConfFromJson(fileName string, v interface{}) error {
    data, err := ioutil.ReadFile(fileName)
    if err != nil {
        log.Error("read error:", err)
        return err
    }
    err = json.Unmarshal(data, &v)
    if err != nil {
        log.Error("decode json error: ", err)
        return err
    }
    fmt.Println("v: ", v)
    return err
}


