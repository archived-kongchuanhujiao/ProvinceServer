package clientspublic

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

func (m *Message) AddText(s string) *Message { m.Chain = append(m.Chain, &Text{s}); return m } // AddText 新增文本

func (m *Message) AddAt(t uint64) *Message { m.Chain = append(m.Chain, &At{t}); return m } // AddAt 新增@

func (m *Message) AddImage(d []byte) *Message { m.Chain = append(m.Chain, &Image{Data: d}); return m } // AddImage 新增图片
