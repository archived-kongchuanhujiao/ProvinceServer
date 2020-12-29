package clients

type (
	Text struct {
		Content string // 内容
	}

	At struct {
		Target uint64 // 目标
	}

	Image struct {
		URL string // 统一资源定位器
	}
)

func (t *Text) TypeName() string { return "text" }

func (a *At) TypeName() string { return "at" }

func (i *Image) TypeName() string { return "image" }
