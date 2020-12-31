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

func (t *Text) TypeName() string  { return "text" }
func (a *At) TypeName() string    { return "at" }
func (i *Image) TypeName() string { return "image" }

func (m *Message) AddText(s string) *Message        { m.Chain = append(m.Chain, &Text{s}); return m }        // AddText 新增文本
func (m *Message) AddAt(t uint64) *Message          { m.Chain = append(m.Chain, &At{t}); return m }          // AddAt 新增@
func (m *Message) AddImage(d []byte) *Message       { m.Chain = append(m.Chain, &Image{Data: d}); return m } // AddImage 新增图片
func (m *Message) SetTarget(t *Target) *Message     { m.Target = t; return m }                               // SetTarget 设置目标
func (m *Message) SetGroupTarget(t *Group) *Message { m.Target.Group = t; return m }                         // SetGroupTarget 设置群目标
func (m *Message) QuickMessage(ms *Message)         { m.Chain = ms.Chain; m.Client.SendMessage(m) }          // QuickMessage 快速消息
