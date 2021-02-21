$module = "github.com/kongchuanhujiao/server"
$main   = "cmd/server/main.go"
$commit = git rev-parse --short HEAD

"开始检查代码错误..."

go vet ${main}

"请确认，按任意键开始编译"
[Console]::ReadKey() | Out-Null

"开始编译..."

$env:GOOS="linux"
$env:GOARCH="amd64"
go build -ldflags "-w -s -X main.Commit=${commit}" ${main}

"编译完成，按任意键退出..."
[Console]::ReadKey() | Out-Null
