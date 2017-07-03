package main

import (
	"fmt"
        "net/http"
        "io/ioutil"
        "strings"
	"strconv"
)

func xhttpGet() {
        resp, err := http.Get("http://172.30.0.8:65534/server_info")
        if err != nil {
        }else{
                defer resp.Body.Close()
                body, err := ioutil.ReadAll(resp.Body)
                if err != nil {
                }else{
                        restart_epoch, _ := strconv.Atoi(strings.Replace(strings.Split(string(body), " ")[5], "\n", "", -1))
                        fmt.Println(restart_epoch)
                }
        }
}

func main(){
	xhttpGet()
}
