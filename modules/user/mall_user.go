package user

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

// UserInfo 定义返回的用户信息结构体
type UserInfo struct {
	ID              int    `json:"id"`
	UserID          string `json:"userId"`
	Nickname        string `json:"nickname"`
	Description     string `json:"description"`
	Gender          string `json:"gender"`
	Birthday        string `json:"birthday"`
	Photo           string `json:"photo"`
	PhoneNumber     string `json:"phoneNumber"`
	Email           string `json:"email"`
	ThirdPartyEmail string `json:"thirdPartyEmail"`
	SignUpType      int    `json:"signUpType"`
	Platform        string `json:"platform"`
	GoogleUsername  string `json:"googleUsername"`
	AppleUsername   string `json:"appleUsername"`
	Password        string `json:"password"`
	ImUsername      string `json:"imUsername"`
	ImPassword      string `json:"imPassword"`
	ImUid           string `json:"imUid"`
	CreateTime      string `json:"createTime"`
	UserType        int    `json:"userType"`
	Deleted         int    `json:"deleted"`
	Status          int    `json:"status"`
	GmtCreate       string `json:"gmtCreate"`
	GmtUpdate       string `json:"gmtUpdate"`
}

// MallAPIResponse 定义API响应结构体
type APIResponse struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Data    []UserInfo `json:"data"`
	TraceID string     `json:"traceId"`
	Success bool       `json:"success"`
}

//type MallUserCtx struct {
//	serverAddresses string
//	log.Log
//	ctx *config.Context
//}

// GetMallUserDetails 发送HTTP请求并获取用户详情
func (m *Service) GetMallUserDetails(uids []string) (map[string]UserInfo, error) {
	//if len(m.serverAddresses) > 0 {
	//	servers := strings.Split(m.serverAddresses, ",")

	//for _, server := range servers {
	//m.Info("echooo inner Push server", zap.String("server", server), zap.String("uid", uid))
	//baseURL := "http://" + server + "/user/inner/im/batchGet"
	baseURL := "http://3.216.154.243:8004/user/inner/im/batchGet"

	params := url.Values{}
	for _, uid := range uids {
		params.Add("uids", uid)
	}
	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// 发送HTTP请求
	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查响应码
	if apiResp.Code != 0 {
		return nil, fmt.Errorf("API请求失败: %s", apiResp.Message)
	}

	// 构造结果map
	resultMap := make(map[string]UserInfo)
	for _, userInfo := range apiResp.Data {
		resultMap[userInfo.ImUid] = userInfo
	}
	if err != nil {
		m.Info("Error reading response body:", zap.Error(err))
		//continue
	} else {
		return resultMap, nil
		//break
	}
	//}
	//}
	//baseURL := "http://3.216.154.243:8004/user/inner/im/batchGet"
	//baseURL := "http://127.0.0.1:8004/user-service/user/inner/im/batchGet"
	return nil, nil

}
