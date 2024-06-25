package user

import (
	"fmt"
	"testing"
)

// GetUserDetailsHandler 处理获取用户详情的请求
func (s *Service) TestGetUserDetailsHandler(t *testing.T) {
	//url := "http://3.216.154.243:8004/user/inner/im/batchGet"
	// 测试数据
	uids := []string{"a49f8c9b4a4b4465aaec417e0dc69ac1", "0c8cde3d959b41ecad0ba62429d6b42e"}

	// 获取用户详情
	userDetails, err := s.GetMallUserDetails(uids)
	if err != nil {
		fmt.Printf("获取用户详情失败: %v\n", err)
		return
	}

	// 输出结果
	for uid, userInfo := range userDetails {
		fmt.Printf("UID: %s, UserInfo: %+v\n", uid, userInfo)
	}
}
