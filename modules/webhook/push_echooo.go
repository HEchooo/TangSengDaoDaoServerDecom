package webhook

import (
	"bytes"
	"encoding/json"
	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/config"
	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/pkg/log"
	"go.uber.org/zap"
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
	log.Info("EchoooPush serverAddresses", zap.String("serverAddresses", m.serverAddresses))
	if len(m.serverAddresses) > 0 {
		servers := strings.Split(m.serverAddresses, ",")
		for _, server := range servers {
			m.Info("echooo inner Push server", zap.String("server", server), zap.String("uid", uid))
			reqParam := SendSinglePushReq{
				UserId:     uid,
				PushType:   3,
				TemplateId: 27,
			}
			jsonData, _ := json.Marshal(&reqParam)
			resp, err := http.Post("http://"+server+"/inner/push/sendNotice", "application/json", bytes.NewBuffer(jsonData))
			defer resp.Body.Close()
			if err != nil {
				m.Info("Error reading response body:", zap.Error(err))
				continue
			} else {
				break
			}
		}
	}
	return nil
}

type SendSinglePushReq struct {
	UserId     string                 `json:"userId"`
	DeviceId   string                 `json:"deviceId"`
	Lang       string                 `json:"lang"`
	PushType   int                    `json:"pushType"`
	TemplateId int                    `json:"templateId"`
	Params     map[string]interface{} `json:"params"`
}
