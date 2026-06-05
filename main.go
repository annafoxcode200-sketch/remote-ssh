package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"remote-ssh/config"
	"strings"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	openapiutil "github.com/alibabacloud-go/openapi-util/service"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	credential "github.com/aliyun/credentials-go/credentials"
)

var Config *config.Config

func loadConfig() {
	// 配置文件，默认在根目录下 config.yaml
	configFile := "config.yaml"
	Config = config.Load(configFile)
}

// Description:
//
// 使用凭据初始化账号Client
//
// @return Client
//
// @throws Exception
func CreateClient() (_result *openapi.Client, _err error) {
	// 工程代码建议使用更安全的无AK方式，凭据配置方式请参见：https://help.aliyun.com/document_detail/378661.html。
	credential, _err := credential.NewCredential(nil)
	if _err != nil {
		return _result, _err
	}

	config := &openapi.Config{
		Credential:      credential,
		AccessKeyId:     tea.String(Config.AliECS.AccessKeyID), // 阿里云 AccessId 和 AccessKey
		AccessKeySecret: tea.String(Config.AliECS.AccessKeySecret),
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Ecs
	config.Endpoint = tea.String(Config.AliECS.Endpoint)
	_result = &openapi.Client{}
	_result, _err = openapi.NewClient(config)
	return _result, _err
}

// Description:
//
// # API 相关
//
// @param opt - int 0 表示关闭实例，对应接口 StopInstance, 1表示启动实例，对应接口 StartInstance
//
// @return OpenApi.Params
func CreateApiInfo(opt int) (_result *openapi.Params) {
	var action string
	// 0 表示关闭实例，对应接口 StopInstance, 1表示启动实例，对应接口 StartInstance
	if opt == 0 {
		action = "StopInstance"
	} else {
		action = "StartInstance"
	}
	params := &openapi.Params{
		// 接口名称
		Action: tea.String(action),
		// 接口版本
		Version: tea.String("2014-05-26"),
		// 接口协议
		Protocol: tea.String("HTTPS"),
		// 接口 HTTP 方法
		Method:   tea.String("POST"),
		AuthType: tea.String("AK"),
		Style:    tea.String("RPC"),
		// 接口 PATH
		Pathname: tea.String("/"),
		// 接口请求体内容格式
		ReqBodyType: tea.String("json"),
		// 接口响应体内容格式
		BodyType: tea.String("json"),
	}
	_result = params
	return _result
}

func close(client *openapi.Client) error {
	// 0 表示关闭实例，对应接口 StopInstance, 1表示启动实例，对应接口 StartInstance
	closeParams := CreateApiInfo(0)

	queries := map[string]interface{}{}
	queries["InstanceId"] = tea.String(Config.AliECS.InstanceId)

	queries["DryRun"] = tea.Bool(false)                 // 是否演练模式，true 进行正常校验，但最终不会执行，返回对应 DryRun 码
	queries["StoppedMode"] = tea.String("StopCharging") // 节省停机模式

	// runtime options
	runtime := &util.RuntimeOptions{}
	request := &openapi.OpenApiRequest{
		Query: openapiutil.Query(queries),
	}

	resp, _err := client.CallApi(closeParams, request, runtime)
	if _err != nil {
		return _err
	}

	// 格式化打印输出
	fomatPrint(resp)
	return nil
}

func start(client *openapi.Client) error {
	// 0 表示关闭实例，对应接口 StopInstance, 1表示启动实例，对应接口 StartInstance
	closeParams := CreateApiInfo(1)

	queries := map[string]interface{}{}
	queries["InstanceId"] = tea.String(Config.AliECS.InstanceId)

	queries["DryRun"] = tea.Bool(false) // 是否演练模式，true 进行正常校验，但最终不会执行，返回对应 DryRun 码

	// runtime options
	runtime := &util.RuntimeOptions{}
	request := &openapi.OpenApiRequest{
		Query: openapiutil.Query(queries),
	}

	resp, _err := client.CallApi(closeParams, request, runtime)
	if _err != nil {
		return _err
	}

	// 格式化打印输出
	fomatPrint(resp)
	return nil
}

func fomatPrint(resp map[string]any) {
	jsonBytes, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		fmt.Printf("格式化 resp 失败: %v\n", err)
	} else {
		fmt.Printf("[LOG] %s\n", string(jsonBytes))
	}
}

func main() {
	loadConfig()

	client, _err := CreateClient()
	if _err != nil {
		log.Fatalf("CreateClient() error: %v", _err)
	}

	// 创建一个读取标准输入的 Reader
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("option start|s:start remote   close|c:close remote  :\n")

	// 读取用户输入直到遇到换行符，并去除首尾空格和换行符
	option, _ := reader.ReadString('\n')
	option = strings.TrimSpace(option)

	// 根据输入执行对应的逻辑
	if option == "start" || option == "s" {
		fmt.Println("start remote ecs")
		err := start(client)
		if err != nil {
			fmt.Printf("start error, %s\n", err)
		}
	} else if option == "close" || option == "c" {
		fmt.Println("close remote ecs")
		err := close(client)
		if err != nil {
			fmt.Printf("close error, %s\n", err)
		}
	} else {
		fmt.Println("选择格式不正确")
	}

	fmt.Println("程序执行完毕,请按回车键关闭窗口")
	reader.ReadString('\n')

	// go build -trimpath -ldflags="-s -w" -o remote-ssh.exe .
}
