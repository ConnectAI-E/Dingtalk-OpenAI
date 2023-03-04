package public

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/eryajf/chatgpt-dingtalk/config"
	"github.com/go-resty/resty/v2"
)

func InitAiCli() *resty.Client {
	return resty.New().SetTimeout(30*time.Second).SetHeader("Authorization", fmt.Sprintf("Bearer %s", config.LoadConfig().ApiKey)).SetProxy(config.LoadConfig().HttpProxy)
}

type Billing struct {
	Object         string  `json:"object"`
	TotalGranted   float64 `json:"total_granted"`
	TotalUsed      float64 `json:"total_used"`
	TotalAvailable float64 `json:"total_available"`
	Grants         struct {
		Object string `json:"object"`
		Data   []struct {
			Object      string  `json:"object"`
			ID          string  `json:"id"`
			GrantAmount float64 `json:"grant_amount"`
			UsedAmount  float64 `json:"used_amount"`
			EffectiveAt float64 `json:"effective_at"`
			ExpiresAt   float64 `json:"expires_at"`
		} `json:"data"`
	} `json:"grants"`
}

func GetBalance() (Billing, error) {
	var data Billing
	url := "https://api.openai.com/dashboard/billing/credit_grants"
	resp, err := InitAiCli().R().Get(url)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return data, err
	}
	return data, nil
}
