package clients

type (
	// Text 文本
	Text struct {
		Content string // 内容
	}

	// At @
	At struct {
		Target uint64 // 目标
	}

	// Image 图片
	Image struct {
		URL  string // 统一资源定位器
		Data []byte // 数据
	}
)

func (t *Text) TypeName() string { return "text" }

func (a *At) TypeName() string { return "at" }

func (i *Image) TypeName() string { return "image" }
