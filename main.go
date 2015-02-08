package main

import (
	"fmt"
    "encoding/json"
)

func main() {
	fmt.Println("Start3...")
  
    

	r:=NewRadio("/dev/spidev0.0")
	for {
		payload := <-r.RcvChan 
        fmt.Printf("Buf: %d \n",payload)
        b, _ := json.Marshal(payload)

        fmt.Printf("JSON: %s \n",b)
	}

}
