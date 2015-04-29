package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"time"
)

func main() {
	url := flag.String("url", "null", "endpoint url")
	flag.Parse()

	fmt.Println("Start3... %s", *url)

	r := NewRadio("/dev/spidev0.0")
	for {
		var buf *bytes.Buffer

		select {
		case rcvPayload := <-r.RcvChan:
			fmt.Printf("Buf: %d \n", rcvPayload)
			b, _ := json.Marshal(rcvPayload)
			fmt.Printf("JSON: %s \n", b)
			buf = bytes.NewBuffer(b)

		case <-time.After(1 * time.Second):
			buf = bytes.NewBufferString("{}")
		}

		resp, err := http.Post(*url, "application/json", buf)

		if err == nil {

			respBuf := new(bytes.Buffer)
			respBuf.ReadFrom(resp.Body)

			sndPayload := new(SndMsg)

			json.Unmarshal(respBuf.Bytes(), sndPayload)
			if sndPayload.Repeat > 0 {
				r.SndChan <- sndPayload
			}
		}

	}

}
