package configs

// GetQQConf 获取 QQ 配置
func GetQQConf() (uint64, string) {
	return conf.Number, conf.Password
}

// GetDatabaseConf 获取数据库配置
func GetDatabaseConf() string {
	return conf.Database
}

// GetJWTConf 获取 JWT 配置
func GetJWTConf() JWT {
	return conf.JWT
}
