package sms

import (
    "encoding/json"
    "fmt"
    "GoHub-Service/pkg/logger"

    "github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

// Aliyun 实现 sms.Driver interface
type Aliyun struct{}

// Send 实现 sms.Driver interface 的 Send 方法
func (s *Aliyun) Send(phone string, message Message, config map[string]string) bool {

    // 检查配置是否完整
    if config["access_key_id"] == "" || config["access_key_secret"] == "" {
        logger.ErrorString("短信[阿里云]", "配置错误", "access_key_id 或 access_key_secret 未配置")
        return false
    }

    if config["sign_name"] == "" || message.Template == "" {
        logger.ErrorString("短信[阿里云]", "配置错误", "sign_name 或 template_code 未配置")
        return false
    }

    // 创建短信客户端
    client, err := dysmsapi.NewClientWithAccessKey(
        "cn-hangzhou",
        config["access_key_id"],
        config["access_key_secret"],
    )

    if err != nil {
        logger.ErrorString("短信[阿里云]", "初始化客户端失败", err.Error())
        return false
    }

    // 构建请求
    request := dysmsapi.CreateSendSmsRequest()
    request.Scheme = "https"
    request.PhoneNumbers = phone
    request.SignName = config["sign_name"]
    request.TemplateCode = message.Template

    // 转换模板参数为 JSON
    templateParam, err := json.Marshal(message.Data)
    if err != nil {
        logger.ErrorString("短信[阿里云]", "解析模板参数错误", err.Error())
        return false
    }
    request.TemplateParam = string(templateParam)

    logger.DebugJSON("短信[阿里云]", "发送请求", map[string]interface{}{
        "phone":         phone,
        "sign_name":     config["sign_name"],
        "template_code": message.Template,
        "template_param": string(templateParam),
    })

    // 发送短信
    response, err := client.SendSms(request)
    if err != nil {
        logger.ErrorString("短信[阿里云]", "发送失败", err.Error())
        return false
    }

    logger.DebugJSON("短信[阿里云]", "接口响应", response)

    // 判断是否发送成功
    if response.Code == "OK" {
        logger.DebugString("短信[阿里云]", "发送成功", fmt.Sprintf("BizId: %s", response.BizId))
        return true
    } else {
        logger.ErrorString("短信[阿里云]", "发送失败", fmt.Sprintf("Code: %s, Message: %s", response.Code, response.Message))
        return false
    }
}
