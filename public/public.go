package public

import (
	"fmt"
	"strings"

	"github.com/eryajf/chatgpt-dingtalk/config"
	"github.com/eryajf/chatgpt-dingtalk/pkg/cache"
	"github.com/eryajf/chatgpt-dingtalk/pkg/dingbot"
	"github.com/eryajf/chatgpt-dingtalk/pkg/logger"
)

var UserService cache.UserServiceInterface
var Config *config.Configuration
var Prompt *[]config.Prompt

func InitSvc() {
	Config = config.LoadConfig()
	Prompt = config.LoadPrompt()
	UserService = cache.NewUserService()
	// 暂时不在初始化时获取余额
	// if Config.Model == openai.GPT3Dot5Turbo0301 || Config.Model == openai.GPT3Dot5Turbo {
	// _, _ = GetBalance()
	// }
}

func FirstCheck(rmsg *dingbot.ReceiveMsg) bool {
	lc := UserService.GetUserMode(rmsg.SenderStaffId)
	if lc == "" {
		if Config.DefaultMode == "串聊" {
			return true
		} else {
			return false
		}
	}
	if lc != "" && strings.Contains(lc, "串聊") {
		return true
	}
	return false
}

// ProcessRequest 分析处理请求逻辑
// 主要提供单日请求限额的功能
func CheckRequest(rmsg *dingbot.ReceiveMsg) bool {
	if Config.MaxRequest == 0 {
		return true
	}
	count := UserService.GetUseRequestCount(rmsg.SenderStaffId)
	// 判断访问次数是否超过限制
	if count >= Config.MaxRequest {
		logger.Info(fmt.Sprintf("亲爱的: %s，您今日请求次数已达上限，请明天再来，交互发问资源有限，请务必斟酌您的问题，给您带来不便，敬请谅解!", rmsg.SenderNick))
		_, err := rmsg.ReplyToDingtalk(string(dingbot.TEXT), fmt.Sprintf("一个好的问题，胜过十个好的答案！\n亲爱的: %s，您今日请求次数已达上限，请明天再来，交互发问资源有限，请务必斟酌您的问题，给您带来不便，敬请谅解!", rmsg.SenderNick))
		if err != nil {
			logger.Warning(fmt.Errorf("send message error: %v", err))
		}
		return false
	}
	// 访问次数未超过限制，将计数加1
	UserService.SetUseRequestCount(rmsg.SenderStaffId, count+1)
	return true
}

var Welcome string = `# 发送信息

若您想给机器人发送信息，有如下两种方式：

1. 在本机器人所在群里@机器人；
2. 点击机器人的头像后，再点击"发消息"。

机器人收到您的信息后，默认会交给chatgpt进行处理。除非，您发送的内容是如下**系统指令**之一。

-----

# 系统指令

系统指令是一些特殊的词语，当您向机器人发送这些词语时，会触发对应的功能：

**单聊**：每条消息都是单独的对话，不包含上下文

**串聊**：对话会携带上下文，除非您主动重置对话或对话长度超过限制

**重置**：重置上下文

**余额**：查询机器人所用OpenAI账号的余额

**模板**：查询机器人内置的快捷模板

**图片**：查看如何根据提示词生成图片

**帮助**：重新获取帮助信息

-----

# 友情提示

使用"串聊模式"会显著加快机器人所用账号的余额消耗速度。

因此，若无保留上下文的需求，建议使用"单聊模式"。

即使有保留上下文的需求，也应适时使用"重置"指令来重置上下文。

-----

# 项目地址

本项目已在GitHub开源，[查看源代码](https://github.com/eryajf/chatgpt-dingtalk)。
`
