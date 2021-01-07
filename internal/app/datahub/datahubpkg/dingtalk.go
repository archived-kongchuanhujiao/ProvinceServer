package datahubpkg

import "coding.net/kongchuanhujiao/server/internal/app/datahub/internal/dingtalk"

func PushMessage(accessToken string, secret string, content string, atMobiles []string, isAtAll bool) error {
	return dingtalk.Push(accessToken, secret, content, atMobiles, isAtAll)
}
