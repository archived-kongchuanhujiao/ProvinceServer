package account

type (
	// Tab 账号
	Tab struct {
		ID   string    `json:"id" db:"id"`     // 标识号
		QQ   uint64    `json:"qq" db:"qq"`     // QQ
		Push PushField `json:"push" db:"push"` // 推送
	}

	// PushField 推送字段
	PushField []struct {
		Platform string `json:"platform"` // 平台
		Key      string `json:"key"`      // 键[ QQ群号, 钉钉AccessToken ]
		Secret   string `json:"secret"`   // 机密[ QQ平台留空，钉钉平台为密钥 ]
	}
)
