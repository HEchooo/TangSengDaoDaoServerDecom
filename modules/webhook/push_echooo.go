package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/config"
	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/pkg/log"
	"net/http"
	"strings"
)

type EchoooPush struct {
	serverAddresses string
	log.Log
	ctx *config.Context
}

func NewEchoooPush(serverAddresses string, ctx *config.Context) *EchoooPush {

	return &EchoooPush{
		ctx:             ctx,
		Log:             log.NewTLog("EchoooPush"),
		serverAddresses: serverAddresses,
	}
}

// Push 推送
func (m *EchoooPush) Push(uid string) error {
	if len(m.serverAddresses) > 0 {
		servers := strings.Split(m.serverAddresses, ",")
		for _, server := range servers {
			fmt.Printf("echooo custom tip Push server=%s,uid=%s", server, uid)
			m.Log.Debug("")
			reqParam := SendSinglePushReq{
				userId:     uid,
				deviceId:   "",
				lang:       "zh-cn",
				pushType:   2,
				templateId: 23,
				params:     make(map[string]interface{}),
			}
			jsonData, _ := json.Marshal(reqParam)
			resp, err := http.Post("http://"+server+"/inner/push/sendNotice", "application/json", bytes.NewBuffer(jsonData))
			defer resp.Body.Close()
			if err != nil {
				continue
			} else {
				break
			}
		}
	}
	return nil
}

type SendSinglePushReq struct {
	userId     string
	deviceId   string
	lang       string
	pushType   int
	templateId int
	params     map[string]interface{}
}
