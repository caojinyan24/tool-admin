package alert

import (
	lark "github.com/larksuite/oapi-sdk-go/v3"
)

var (
	appId     = "cli_a8fcbf41bc489013"
	appSecret = "bQeeQPiZsBc3rr3aWILDodNrA4fNU4rf"
)

// // 创建一个自定义HTTP客户端，临时跳过TLS证书验证
// var httpClient = &http.Client{
// 	Transport: &http.Transport{
// 		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 注意：这仅用于测试环境，生产环境不应使用
// 	},
// }

// // 使用自定义HTTP客户端初始化LarkClient
// var LarkClient = lark.NewClient(
//
//	appId,
//	appSecret,
//	lark.WithHttpClient(httpClient),
//
// )
// 使用自定义HTTP客户端初始化LarkClient
var LarkClient = lark.NewClient(appId, appSecret)
