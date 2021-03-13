$module = "github.com/kongchuanhujiao/server"
$main   = "cmd/server/main.go"
$commit = git rev-parse --short HEAD

$env:GOOS   = "linux"
$env:GOARCH = "amd64"

"构建信息：CommitShortID ${commit}"

"开始构建..."

go build -ldflags "-w -s -X ${module}/internal/pkg/config.Commit=${commit}" -o .kongchuanhujiao/serve ${main}

"构建完成"
