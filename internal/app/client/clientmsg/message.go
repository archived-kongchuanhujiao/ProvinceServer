package clientmsg

type (
	Callback func(*Message) // Callback 回调

	Element interface { // Element 消息链
		TypeName() string // TypeName 类型名
	}

	Text  struct{ Content string } // Text 文本
	At    struct{ Target uint64 }  // At @
	Image struct {                 // Image 图片
		URL  string
		Data []byte
	}

	Message struct { // Message 消息
		Chain  []Element // 消息链
		Target *Target   // 目标
	}
	Target struct { // Target 目标
		ID    uint64 // 唯一识别码
		Name  string // 名称
		Group *Group // 群
	}
	Group struct { // Group 群
		ID   uint64 // 唯一识别码
		Name string // 名称
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

func NewMessage() *Message              { return &Message{Chain: []Element{}, Target: &Target{}} } // NewMessage 新建消息
func NewTextMessage(s string) *Message  { return NewMessage().AddText(s) }                         // NewTextMessage 新建文本消息
func NewAtMessage(t uint64) *Message    { return NewMessage().AddAt(t) }                           // NewAtMessage 新建@消息
func NewImageMessage(d []byte) *Message { return NewMessage().AddImage(d) }                        // NewImageMessage 新建图片消息
