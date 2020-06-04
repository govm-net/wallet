package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	mux := NewRouter()
	mux.Handle("/", http.FileServer(http.Dir(conf.StaticFiles)))
	srv := http.Server{}
	srv.Handler = mux
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", conf.HTTPPort))
	if err != nil {
		return
	}
	addr := "http://" + ln.Addr().String()
	fmt.Println("listen address:", addr)
	go func(addr string) {
		time.Sleep(time.Second)
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "windows":
			cmd = exec.Command("cmd", "/c", "start", addr)
		case "darwin":
			cmd = exec.Command("open", addr)
		case "linux":
			cmd = exec.Command("xdg-open", addr)
		}
		err := cmd.Start()
		if err != nil {
			fmt.Println("fail to open browser. ", err)
			return
		}
	}(addr)
	err = srv.Serve(ln)
	if err != nil {
		fmt.Println("fail to Serve", err)
	}
}
