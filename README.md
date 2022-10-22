# gh2
go-http2-client


``` go get github.com/lllk140/gh2 ```

```

package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/lllk140/gh2/GH2"
	"golang.org/x/net/http2/hpack"
)

func main() {
	var TlsConfig *tls.Config = &tls.Config{
		NextProtos: []string{"h2"},
	}

	var con, _ = tls.Dial("tcp", "httpbin.org:443", TlsConfig)

	var th2 = new(GH2.H2Connection)

	th2.InitiateConnection()
	th2.SendSettings(0, nil, 0)
	var SettingsData = th2.DataToSend()
	_, _ = con.Write(SettingsData)

	//POST
	var headers = []hpack.HeaderField{
		{Name: ":method", Value: "POST"},
		{Name: ":path", Value: "/post"},
		{Name: ":authority", Value: "httpbin.org"},
		{Name: ":scheme", Value: "https"},
		{Name: "user-agent", Value: "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:105.0) Gecko/20100101 Firefox/105.0"},
	}
	th2.SendHeaders(1, headers, 4)
	th2.SendData(1, []byte("next"), 1)
	var HeadersData = th2.DataToSend()
	_, _ = con.Write(HeadersData)

	var data []byte
	for {
		var buf = make([]byte, 8196)
		var length, _ = con.Read(buf)

		var events = th2.ReceiveData(buf[:length])
		for _, event := range events {
			if value, ok := event.(*GH2.HeadersFrame); ok == true {
				fmt.Printf("%v\n", value.Headers)
			}
			if value, ok := event.(*GH2.DataFrame); ok == true {
				data = bytes.Join([][]byte{data, value.Body}, []byte(""))
			}
			if _, ok := event.(*GH2.EndStream); ok == true {
				goto EXIT
			}
		}
	}
EXIT:
	fmt.Printf("%v\n", string(data))
	th2.CloseConnection(1, 0, 0)
	_, _ = con.Write(th2.DataToSend())
	_ = con.Close()
}

```
