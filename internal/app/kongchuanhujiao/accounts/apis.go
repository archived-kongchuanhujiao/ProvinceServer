package accounts

// TODO https://qianjunakasumi.coding.net/p/kongchuanhujiao/requirements/issues/15/detail
/*
TODO
   MVC API 参考 wenda 包接入
   注册账号 API 实现步骤：
   接收POST注册账号 POST /apis/accounts/register/verifier 要求QQ字段 QQ：uint64
   向QQ号发送好友消息（消息内容：验证码,4位，有效期5min）
   客户端 POST /apis/accounts/register  要求 QQ 、验证码和ID 字段  验证码：uint8 id:string
   验证 验证码：通过map[uint64(QQ)]uint8(验证码) // 对比，如果一致则注册成功
   向数据库写入，数据API我后面会写
*/
