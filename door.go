package main

import (
	"fmt"
	"github.com/davecheney/gpio"
	"net/http"
	//	"os"
	//	"os/signal"
	"time"
)

func main() {
	p1, _ := gpio.OpenPin(gpio.GPIO2, gpio.ModeOutput)
	p2, _ := gpio.OpenPin(gpio.GPIO3, gpio.ModeOutput)
	p3, _ := gpio.OpenPin(gpio.GPIO4, gpio.ModeOutput)
	p1.Clear()
	p2.Clear()
	p3.Clear()

	http.HandleFunc("/api/door/close", func(w http.ResponseWriter, r *http.Request) {
		p1.Set()
		fmt.Println("Press CLOSE")
		time.Sleep(1000 * time.Millisecond)
		fmt.Println("Release CLOSE")
		p1.Clear()
		fmt.Fprintf(w, "OK: CLOSE")
	})
	http.HandleFunc("/api/door/open", func(w http.ResponseWriter, r *http.Request) {
		p2.Set()
		fmt.Println("Press OPEN")
		time.Sleep(1000 * time.Millisecond)
		fmt.Println("Release OPEN")
		p2.Clear()
		fmt.Fprintf(w, "OK: OPEN")
	})

	http.ListenAndServe(":8011", nil)
}
