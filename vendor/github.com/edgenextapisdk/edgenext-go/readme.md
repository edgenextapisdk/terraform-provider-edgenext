# edgenext.com api sdk for go

### 说明

* edgenext scdn产品，官网地址：https://www.edgenext.com/
* 接口遵循RESTful，默认请求体json，接口默认返回json

### 签名算法

* 每次请求都签名，保证传输过程数据不被篡改
* 客户端：sha256签名算法，将参数base64编码+app_secret用sha256签名，每次请求带上签名
* 服务端：拿到参数用相同的算法签名，对比签名是否正确

### sdk 使用说明

* 环境：go >=1.14
* 支持get/post/patch/put/delete方法
* 参数说明
    * AppId 分配的app_id
    * appSecert 分配的appSecert, 用于签名数据
    * apiUrlPre api地址前缀, 请具体咨询运营人员
    * userId 当前使用者的用户ID
* 每次调用会返回JSONObject, 如果执行过程中有异常，会直接抛出异常；
* 如果需要调试，可以调用debug方法
* 注意事项
    针对所有请求，uri与get参数是分离的，如 https://apiv4.local.com/V4/version?v=1, 调用时v=1参数，须通过ReqParams.Query传递

## 安装

```
go get github.com/edgenextapisdk/edgenext-go
```

### 使用

#### 示例类
```
package main

import (
        "os"
        sdk "github.com/edgenextapisdk/edgenext-go"
        "fmt"
)

func main() {
	app_id := os.Getenv("SDK_APP_ID")                     //这里是用的环境变量，业务中使用时，请替换成具体的参数
	app_secret := os.Getenv("SDK_APP_SECERT")             //这里是用的环境变量，业务中使用时，请替换成具体的参数
	api_pre := os.Getenv("SDK_API_PRE")                   //这里是用的环境变量，业务中使用时，请替换成具体的参数

        sdkObj := sdk.Sdk{
                AppId: app_id,
                AppSecret: app_secret,
                ApiPre: api_pre,
                Timeout: 30,
        }
        var api string
        var err error
        //ReqParams 有三个属性，如果用不到，不设置即可：
        //ReqParams.Query 对应的是GET请求的参数，map[string]interface{}
        //ReqParams.Data  对应的是非GET请求的参数，map[string]interface{}
        //ReqParams.Header  对应的是发起请求的header头，map[string]string
        var reqParams sdk.ReqParams

        // Response包括响应的全部信息，其中：
        // Response.HttpCode Http请求响应状态码，成功是200
        // Response.RespBody 请求返回的body字符串
        // Response.BizCode 是业务的状态码，1代码请求成功，非1代码请求失败
        // Response.BizMsg  是业务的状态码对应的信息
        // Response.BizData  返回的业务数据，只有BizCode为1时，才会有数据
        var resp *sdk.Response

        api = "test.sdk.get"
        reqParams = sdk.ReqParams{
                Query: map[string]interface{}{
                        "page": 1,
                        "pagesize": 10,
                        "data": map[string]interface{}{
                                "name": "name名称",
                                "domain": "baidu.com",
                        },
                },
        }
        resp, err = sdkObj.Get(api, reqParams)
        if err == nil {
                if resp.BizCode == 1 {
                        fmt.Println(api, " 业务处理成功")
                } else {
                        fmt.Println(api, " 业务处理错误: ")
                }
                fmt.Println(api, " http_code: ", resp.HttpCode)
                fmt.Println(api, " body: ", resp.RespBody)
                fmt.Println(api, " biz_code: ", resp.BizCode)
                fmt.Println(api, " biz_msg: ", resp.BizMsg)
                fmt.Println(api, " biz_data: ", resp.BizData)
                fmt.Println(api, " err: ", err)
        } else {
                fmt.Println(api, " request error: ", err)
        }
        fmt.Println("")


        api = "test.sdk.post"
        reqParams = sdk.ReqParams{
                Data: map[string]interface{}{
                        "name": 1,
                        "age": 10,
                        "data": map[string]interface{}{
                                "name": "name名称",
                                "domain": "baidu.com",
                        },
                },
        }
        resp, err = sdkObj.Post(api, reqParams)
        if err == nil {
                if resp.BizCode == 1 {
                        fmt.Println(api, " 业务处理成功")
                } else {
                        fmt.Println(api, " 业务处理错误: ")
                }
                fmt.Println(api, " http_code: ", resp.HttpCode)
                fmt.Println(api, " body: ", resp.RespBody)
                fmt.Println(api, " biz_code: ", resp.BizCode)
                fmt.Println(api, " biz_msg: ", resp.BizMsg)
                fmt.Println(api, " biz_data: ", resp.BizData)
                fmt.Println(api, " err: ", err)
        } else {
                fmt.Println(api, " request error: ", err)
        }
        fmt.Println("")


        api = "test.sdk.delete"
        reqParams = sdk.ReqParams{
                Data: map[string]interface{}{
                        "id": 10,
                },
        }
        resp, err = sdkObj.Delete(api, reqParams)
        if err == nil {
                if resp.BizCode == 1 {
                        fmt.Println(api, " 业务处理成功")
                } else {
                        fmt.Println(api, " 业务处理错误: ")
                }
                fmt.Println(api, " http_code: ", resp.HttpCode)
                fmt.Println(api, " body: ", resp.RespBody)
                fmt.Println(api, " biz_code: ", resp.BizCode)
                fmt.Println(api, " biz_msg: ", resp.BizMsg)
                fmt.Println(api, " biz_data: ", resp.BizData)
                fmt.Println(api, " err: ", err)
        } else {
                fmt.Println(api, " request error: ", err)
        }
        fmt.Println("")


        api = "test.sdk.put"
        reqParams = sdk.ReqParams{
                Data: map[string]interface{}{
                        "name": 1,
                        "age": 10,
                        "data": map[string]interface{}{
                                "name": "name名称",
                                "domain": "baidu.com",
                        },
                },
        }
        resp, err = sdkObj.Put(api, reqParams)
        if err == nil {
                if resp.BizCode == 1 {
                        fmt.Println(api, " 业务处理成功")
                } else {
                        fmt.Println(api, " 业务处理错误: ")
                }
                fmt.Println(api, " http_code: ", resp.HttpCode)
                fmt.Println(api, " body: ", resp.RespBody)
                fmt.Println(api, " biz_code: ", resp.BizCode)
                fmt.Println(api, " biz_msg: ", resp.BizMsg)
                fmt.Println(api, " biz_data: ", resp.BizData)
                fmt.Println(api, " err: ", err)
        } else {
                fmt.Println(api, " request error: ", err)
        }
        fmt.Println("")
}
```

#### 示例类执行输出
···
test.sdk.get  业务处理成功
test.sdk.get  http_code:  200
test.sdk.get  body:  {"status":{"code":1,"message":"操作成功"},"data":{"data":{"domain":"baidu.com","name":"name"},"page":"1","pagesize":"10"}}
test.sdk.get  biz_code:  1
test.sdk.get  biz_msg:  操作成功
test.sdk.get  biz_data:  map[data:map[domain:baidu.com name:name] page:1 pagesize:10]
test.sdk.get  err:  <nil>

test.sdk.post  业务处理成功
test.sdk.post  http_code:  200
test.sdk.post  body:  {"status":{"code":1,"message":"操作成功"},"data":{"age":10,"data":{"domain":"baidu.com","name":"name"},"name":1}}
test.sdk.post  biz_code:  1
test.sdk.post  biz_msg:  操作成功
test.sdk.post  biz_data:  map[age:10 data:map[domain:baidu.com name:name] name:1]
test.sdk.post  err:  <nil>

test.sdk.delete  业务处理成功
test.sdk.delete  http_code:  200
test.sdk.delete  body:  {"status":{"code":1,"message":"操作成功"},"data":{"id":10}}
test.sdk.delete  biz_code:  1
test.sdk.delete  biz_msg:  操作成功
test.sdk.delete  biz_data:  map[id:10]
test.sdk.delete  err:  <nil>

test.sdk.put  业务处理成功
test.sdk.put  http_code:  200
test.sdk.put  body:  {"status":{"code":1,"message":"操作成功"},"data":{"age":10,"data":{"domain":"baidu.com","name":"name"},"name":1}}
test.sdk.put  biz_code:  1
test.sdk.put  biz_msg:  操作成功
test.sdk.put  biz_data:  map[age:10 data:map[domain:baidu.com name:name] name:1]
test.sdk.put  err:  <nil>
···

### 更新日志

* 2022.11.08 

第一版SDK

* 2022.11.11

1. 单元测试增加特殊字符
2. GET请求修正参数乱序导致的签名失败问题

* 2024.06.27

1. 调整目录结构，以适配标准化的gomod使用方式，未来发布将使用版本号
