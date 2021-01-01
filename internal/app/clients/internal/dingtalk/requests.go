package dingtalk

import (
	"bytes"
	"coding.net/kongchuanhujiao/server/internal/app/clients/clientspublic"
	json2 "encoding/json"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strconv"
)

// 储存所有 API 链接
var (
	mainAPI = "https://oapi.dingtalk.com/"
	sendMsg = "robot/send?access_token="
)

func sendMessage(c *DingTalk, m *clientspublic.Message) bool {
	json, err := json2.Marshal(c.transformToDTMessage(m))

	if err != nil {
		loggerr.Warn("转换钉钉消息出现异常", zap.Error(err))
		return false
	}

	_, ok := request(c, sendMsg, json)

	return ok
}

func request(c *DingTalk, subUrl string, json []byte) (errRes ErrResponse, ok bool) {
	ts, sign := generateSign(c.Secret)

	req, err := http.NewRequest(http.MethodPost, mainAPI+subUrl+c.AccessToken+"&timestamp="+ts+"&sign="+sign, bytes.NewReader(json))

	if err != nil {
		loggerr.Warn("调用钉钉接口出现异常", zap.Error(err))
		return
	}

	req.Header.Add("Content-Type", "application/json")

	resp := req.Response

	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		loggerr.Warn("调用钉钉接口出现异常, 响应码: " + strconv.Itoa(resp.StatusCode))
		return
	}

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		loggerr.Warn("解析钉钉接口传回状态异常", zap.Error(err))
		return
	}

	err = json2.Unmarshal(b, &errRes)

	if err != nil {
		loggerr.Warn("解析钉钉接口传回状态异常", zap.Error(err))
		return
	}

	return errRes, true
}
