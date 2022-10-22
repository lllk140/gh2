package PH2

import (
	"bytes"
	"gh2/GH2/IntBinary"
	"golang.org/x/net/http2/hpack"
)

type HEADERS []hpack.HeaderField

type H2Connection struct {
	_dataToSend []byte
}

func (receiver *H2Connection) addDataToSend(body []byte) {
	var data = [][]byte{receiver._dataToSend, body}
	receiver._dataToSend = bytes.Join(data, []byte(""))
}

func (receiver *H2Connection) DataToSend() []byte {
	var body = receiver._dataToSend
	receiver._dataToSend = []byte("")
	return body
}

func (receiver *H2Connection) InitiateConnection() {
	var preamble = []byte("PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n")
	var data = [][]byte{receiver._dataToSend, preamble}
	receiver._dataToSend = bytes.Join(data, []byte(""))
}

func (receiver *H2Connection) SendData(StreamId int64, data []byte, Flags int64) {
	var frame = NewDataFrame(StreamId, Flags)
	var s = [][]byte{frame.buildHeader(data), data}
	receiver.addDataToSend(bytes.Join(s, []byte("")))
}

func (receiver *H2Connection) SendHeaders(StreamId int64, hs HEADERS, Flags int64) {
	var frame = NewHeadersFrame(StreamId, Flags)
	var buf bytes.Buffer
	var ehp = hpack.NewEncoder(&buf)
	for _, header := range hs {
		_ = ehp.WriteField(header)
	}
	var data = buf.Bytes()
	var s = [][]byte{frame.buildHeader(data), data}
	receiver.addDataToSend(bytes.Join(s, []byte("")))
}

func (receiver *H2Connection) SendPriority(StreamId int64, Flags int64) {
	var frame = NewPriorityFrame(StreamId, Flags)
	var data = []byte("")
	var s = [][]byte{frame.buildHeader(data), data}
	receiver.addDataToSend(bytes.Join(s, []byte("")))
}

func (receiver *H2Connection) SendRstStream(StreamId int64, Flags int64) {
	var frame = NewRstStreamFrame(StreamId, Flags)
	var data = []byte("")
	var s = [][]byte{frame.buildHeader(data), data}
	receiver.addDataToSend(bytes.Join(s, []byte("")))
}

func (receiver *H2Connection) SendSettings(StreamId int64, setting *SettingsFrame, Flags int64) {
	var frame = setting
	if setting == nil {
		frame = NewSettingsFrame(StreamId, Flags)
	}
	var bodyList [][]byte
	for _, value := range frame.getSettings() {
		var elems = IntBinary.Pack("HL", value)
		bodyList = append(bodyList, elems)
	}
	var data = bytes.Join(bodyList, []byte(""))
	var s = [][]byte{frame.buildHeader(data), data}
	receiver.addDataToSend(bytes.Join(s, []byte("")))
}

func (receiver *H2Connection) CloseConnection(StreamId int64, code int, Flags int64) {
	var frame = NewGoawayFrame(StreamId, Flags)
	frame.ErrorCode = code
	var data = IntBinary.Pack("L", []int{code})
	var s = [][]byte{frame.buildHeader(data), data}
	receiver.addDataToSend(bytes.Join(s, []byte("")))
}

func (receiver *H2Connection) ReceiveData(ReceiveBody []byte) []interface{} {
	var events []interface{}
	for len(ReceiveBody) > 0 {
		var body = ReceiveBody[:9]
		var header = IntBinary.UnPack("BHBBL", body)
		var bodyLength = header[1]
		var FrameType = header[2]
		var Flags = header[3]
		var StreamId = header[4]

		switch FrameType {
		case 0:
			// data帧
			var frame = NewDataFrame(StreamId, Flags)
			frame.Body = ReceiveBody[9 : bodyLength+9]
			frame.Length = bodyLength
			events = append(events, frame)
			ReceiveBody = ReceiveBody[bodyLength+9:]
			if frame.Flags == 1 {
				events = append(events, NewEndStream(StreamId, Flags))
			}
		case 1:
			// headers帧
			var frame = NewHeadersFrame(StreamId, Flags)
			var Dhp = hpack.NewDecoder(uint32(bodyLength), nil)
			var Headers, _ = Dhp.DecodeFull(ReceiveBody[9 : bodyLength+9])
			frame.Headers = Headers
			frame.Length = bodyLength
			events = append(events, frame)
			ReceiveBody = ReceiveBody[bodyLength+9:]
		case 2:
			// Priority帧
			var frame = NewPriorityFrame(StreamId, Flags)
			frame.Length = bodyLength
			events = append(events, frame)
			ReceiveBody = ReceiveBody[bodyLength+9:]
		case 3:
			// RstStream帧
			var frame = NewRstStreamFrame(StreamId, Flags)
			frame.Length = bodyLength
			events = append(events, frame)
			ReceiveBody = ReceiveBody[bodyLength+9:]
		case 4:
			var frame = NewSettingsFrame(StreamId, Flags)
			var settings = ReceiveBody[9 : bodyLength+9]
			for len(settings) > 0 {
				var ss = settings[:6]
				frame.setSettings(IntBinary.UnPack("HL", ss))
				settings = settings[6:]
			}
			frame.Length = bodyLength
			events = append(events, frame)
			ReceiveBody = ReceiveBody[bodyLength+9:]
		case 8:
			// WINDOW_UPDATE帧
			var frame = NewWindowUpdateFrame(StreamId, Flags)
			var Data = ReceiveBody[9 : bodyLength+9]
			var Deltas = IntBinary.UnPack("L", Data)
			frame.Delta = int(Deltas[0])
			frame.Length = bodyLength
			events = append(events, frame)
			ReceiveBody = ReceiveBody[bodyLength+9:]
		default:
			ReceiveBody = ReceiveBody[bodyLength+9:]
		}
	}
	return events
}
