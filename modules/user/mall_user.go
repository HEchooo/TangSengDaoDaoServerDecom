package user

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// UserInfo 定义返回的用户信息结构体
type UserInfo struct {
	ID              int    `json:"id"`
	UserID          string `json:"userId"`
	Nickname        string `json:"nickname"`
	Description     string `json:"description"`
	Gender          int    `json:"gender"`
	Password        string `json:"password"`
	ImUsername      string `json:"imUsername"`
	ImPassword      string `json:"imPassword"`
	ImUid           string `json:"imUid"`
	GmtCreate       string `json:"gmtCreate"`
	GmtUpdate       string `json:"gmtUpdate"`
	Birthday        string `json:"birthday"`
	Photo           string `json:"photo"`
	PhoneNumber     string `json:"phoneNumber"`
	Email           string `json:"email"`
	ThirdPartyEmail string `json:"thirdPartyEmail"`
	SignUpType      int    `json:"signUpType"`
	Platform        string `json:"platform"`
	GoogleUsername  string `json:"googleUsername"`
	AppleUsername   string `json:"appleUsername"`
	CreateTime      string `json:"createTime"`
	UserType        int    `json:"userType"`
	Deleted         int    `json:"deleted"`
	Status          int    `json:"status"`
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
	// todo 这里后面改成从nacos获取地址
	baseURL := "http://10.10.10.10:8004/user/inner/im/batchGet"

	uidParam := strings.Join(uids, ",")
	// 构建请求 URL
	reqURL := fmt.Sprintf("%s?uids=%s", baseURL, url.QueryEscape(uidParam))
	fmt.Errorf("请求体打印: %v", reqURL)

	// 创建一个新的 HTTP 请求
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	// 设置头信息
	//req.Header.Set("Connection", "Keep-Alive")
	//req.Header.Set("User-Agent", "Apache-HttpClient/4.5.14 (Java/17.0.7)")
	//req.Header.Set("Accept-Encoding", "br,deflate,gzip,x-gzip")
	req.Header.Set("Accept", "application/json")

	// 创建一个 HTTP 客户端并发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
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

func GetMallUserDetailsTest(uids []string) (map[string]UserInfo, error) {
	//if len(m.serverAddresses) > 0 {
	//	servers := strings.Split(m.serverAddresses, ",")

	//for _, server := range servers {
	//m.Info("echooo inner Push server", zap.String("server", server), zap.String("uid", uid))
	//baseURL := "http://" + server + "/user/inner/im/batchGet"
	baseURL := "http://127.0.0.1:8004/user-service/user/inner/im/batchGet"

	uidParam := strings.Join(uids, ",")
	// 构建请求 URL
	reqURL := fmt.Sprintf("%s?uids=%s", baseURL, url.QueryEscape(uidParam))
	log.Println("Failed to create request:  " + reqURL)

	// 创建一个新的 HTTP 请求
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	// 设置头信息
	//req.Header.Set("Connection", "Keep-Alive")
	//req.Header.Set("User-Agent", "Apache-HttpClient/4.5.14 (Java/17.0.7)")
	//req.Header.Set("Accept-Encoding", "br,deflate,gzip,x-gzip")
	req.Header.Set("Accept", "application/json")

	// 创建一个 HTTP 客户端并发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
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
		//Info("Error reading response body:", zap.Error(err))
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
