package PH2

import (
	"github.com/lllk140/gh2/GH2/IntBinary"
	"golang.org/x/net/http2/hpack"
)

// Frame
// FrameType: 帧类型
// StreamId: 标识
// Length: 内容长度
// Flags: 布尔标识
type Frame struct {
	FrameType int64
	StreamId  int64
	Length    int64
	Flags     int64
}

// buildHeader
// 生成帧头
func (receiver *Frame) buildHeader(body []byte) []byte {
	receiver.Length = int64(len(body))
	//if receiver.Length > 16384 {
	//	return nil
	//}
	var H1 = int(receiver.Length & 0xFF)
	var B1 = int(receiver.FrameType)
	var B2 = int(receiver.Flags)
	var L2 = int(receiver.StreamId & 0x7FFFFFFF)
	var values = []int{0, H1, B1, B2, L2}
	return IntBinary.Pack("BHBBL", values)
}

// EndStream
// 已经结束哩
type EndStream struct {
	Frame
}

// DataFrame
// -> Data帧(0)
// Body: 内容
type DataFrame struct {
	Frame
	Body []byte
}

// HeadersFrame
// -> Headers帧(1)
// Headers: 请求/响应 头
type HeadersFrame struct {
	Frame
	Headers []hpack.HeaderField
}

// PriorityFrame
// -> Priority帧(2)
type PriorityFrame struct {
	Frame
}

// RstStreamFrame
// -> RstStream帧(3)
type RstStreamFrame struct {
	Frame
}

// SettingsFrame
// -> Settings帧(4)
// SettingsHeaderTableSize: 远端报头压缩表的最大承载量
// SettingsEnablePush: 服务器推送
// SettingsMaxConcurrentStreams: 发送端允许接收端创建的最大并发流的数量
// SettingsInitialWindowSize: 发送端流量控制的初始窗口大小
// SettingsMaxFrameSize: 接收最大帧大小
// SettingsMaxHeaderListSize: 可接收的header列表长度
// SettingsMaxClosedStreams: 最大关闭流数
type SettingsFrame struct {
	Frame
	SettingsHeaderTableSize      int
	SettingsEnablePush           int
	SettingsMaxConcurrentStreams int
	SettingsInitialWindowSize    int
	SettingsMaxFrameSize         int
	SettingsMaxHeaderListSize    int
	SettingsMaxClosedStreams     int
}

// setSettings
// 设置设置列表
func (receiver *SettingsFrame) setSettings(body []int64) {
	switch body[0] {
	case 1:
		receiver.SettingsHeaderTableSize = int(body[1])
	case 2:
		receiver.SettingsEnablePush = int(body[1])
	case 3:
		receiver.SettingsMaxConcurrentStreams = int(body[1])
	case 4:
		receiver.SettingsInitialWindowSize = int(body[1])
	case 5:
		receiver.SettingsMaxFrameSize = int(body[1])
	case 6:
		receiver.SettingsMaxHeaderListSize = int(body[1])
	case 8:
		receiver.SettingsMaxClosedStreams = int(body[1])
	default:
		break
	}
}

// getSettings
// 返回设置列表
func (receiver *SettingsFrame) getSettings() [][]int {
	if receiver.Flags == 1 {
		return [][]int{}
	}
	var settingList = make([][]int, 7)
	settingList[0] = []int{1, receiver.SettingsHeaderTableSize}
	settingList[1] = []int{2, receiver.SettingsEnablePush}
	settingList[2] = []int{4, receiver.SettingsInitialWindowSize}
	settingList[3] = []int{5, receiver.SettingsMaxFrameSize}
	settingList[4] = []int{8, receiver.SettingsMaxClosedStreams}
	settingList[5] = []int{3, receiver.SettingsMaxConcurrentStreams}
	settingList[6] = []int{6, receiver.SettingsMaxHeaderListSize}

	var settings [][]int
	for _, setting := range settingList {
		if setting[1] < 0 {
			settings = append(settings, setting)
		}
	}
	return settings
}

// PushPromiseFrame
// -> PushPromise帧(5)
type PushPromiseFrame struct {
	Frame
}

// PingFrame
// -> Ping帧(6)
type PingFrame struct {
	Frame
}

// GoawayFrame
// -> Goaway帧(7)
type GoawayFrame struct {
	Frame
	ErrorCode int
}

// WindowUpdateFrame
// -> WindowUpdate帧(8)
type WindowUpdateFrame struct {
	Frame
	Delta int
}

// ContinuationFrame
// -> Continuation帧(9)
type ContinuationFrame struct {
	Frame
}

// NewDataFrame
// 创建一个包含 Flags,StreamId 的空Data帧
func NewDataFrame(StreamId int64, Flags int64) *DataFrame {
	var frame = new(DataFrame)
	frame.StreamId = StreamId
	frame.FrameType = 0
	frame.Flags = Flags
	return frame
}

// NewHeadersFrame
// 创建一个包含 Flags,StreamId 的空Headers帧
func NewHeadersFrame(StreamId int64, Flags int64) *HeadersFrame {
	var frame = new(HeadersFrame)
	frame.StreamId = StreamId
	frame.FrameType = 1
	frame.Flags = Flags
	return frame
}

// NewPriorityFrame
// 创建一个包含 Flags,StreamId 的空Priority帧
func NewPriorityFrame(StreamId int64, Flags int64) *PriorityFrame {
	var frame = new(PriorityFrame)
	frame.StreamId = StreamId
	frame.FrameType = 2
	frame.Flags = Flags
	return frame
}

// NewRstStreamFrame
// 创建一个包含 Flags,StreamId 的空RstStream帧
func NewRstStreamFrame(StreamId int64, Flags int64) *RstStreamFrame {
	var frame = new(RstStreamFrame)
	frame.StreamId = StreamId
	frame.FrameType = 3
	frame.Flags = Flags
	return frame
}

// NewSettingsFrame
// 创建一个包含 Flags,StreamId, Settings.* 的空SettingsFrame帧
func NewSettingsFrame(StreamId int64, Flags int64) *SettingsFrame {
	var frame = new(SettingsFrame)

	frame.SettingsHeaderTableSize = 4096
	frame.SettingsEnablePush = 1
	frame.SettingsMaxConcurrentStreams = 100
	frame.SettingsInitialWindowSize = 65535
	frame.SettingsMaxFrameSize = 16384
	frame.SettingsMaxHeaderListSize = 65536
	frame.SettingsMaxClosedStreams = 0

	frame.StreamId = StreamId
	frame.FrameType = 4
	frame.Flags = Flags
	return frame
}

// NewGoawayFrame
// 创建一个包含 Flags,StreamId 的空Goaway帧
func NewGoawayFrame(StreamId int64, Flags int64) *GoawayFrame {
	var frame = new(GoawayFrame)
	frame.StreamId = StreamId
	frame.FrameType = 7
	frame.Flags = Flags
	return frame
}

// NewWindowUpdateFrame
// 创建一个包含 Flags,StreamId 的空WindowUpdate帧
func NewWindowUpdateFrame(StreamId int64, Flags int64) *WindowUpdateFrame {
	var frame = new(WindowUpdateFrame)
	frame.StreamId = StreamId
	frame.FrameType = 8
	frame.Flags = Flags
	return frame
}

// NewEndStream
// 已经结束哩
func NewEndStream(StreamId int64, Flags int64) *EndStream {
	var frame = new(EndStream)
	frame.StreamId = StreamId
	frame.FrameType = -1
	frame.Flags = Flags
	return frame
}
