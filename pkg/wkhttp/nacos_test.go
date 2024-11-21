package wkhttp

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"testing"
	"time"
)

func Test1(t *testing.T) {
	// 从控制台命名空间管理的"命名空间详情"中拷贝 End Point、命名空间 ID
	var endpoint = "3.216.154.243"
	var namespaceId = "64ca0ecd-adae-4bde-83c4-d53fde3b45f8"

	// 推荐使用 RAM 用户的 accessKey、secretKey
	var username = "dev"
	var password = "valleysound789"

	clientConfig := constant.ClientConfig{
		//
		Endpoint:       endpoint + ":8848",
		NamespaceId:    namespaceId,
		AppName:        "test",
		Username:       username,
		Password:       password,
		TimeoutMs:      5 * 1000,
		ListenInterval: 30 * 1000,
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"clientConfig": clientConfig,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	var dataId = "com.alibaba.nacos.example.properties"
	var group = "DEFAULT_GROUP"

	// 发布配置
	success, err := configClient.PublishConfig(vo.ConfigParam{
		DataId:  dataId,
		Group:   group,
		Content: "connectTimeoutInMills=3000"})

	if success {
		fmt.Println("Publish config successfully.")
	}

	time.Sleep(3 * time.Second)

	// 获取配置
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group})

	fmt.Println("Get config：" + content)

	// 监听配置
	configClient.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("ListenConfig group:" + group + ", dataId:" + dataId + ", data:" + data)
		},
	})

	// 删除配置
	success, err = configClient.DeleteConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group})

	if success {
		fmt.Println("Delete config successfully.")
	}

}
