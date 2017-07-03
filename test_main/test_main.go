package main

import (
	"fmt"
        "net/http"
        "io/ioutil"
        "strings"
	"strconv"
)

func httpGet() {
        resp, err := http.Get("http://172.30.0.8:65534/server_info")
        if err != nil {
                fmt.Println("123123")
        }else{
                defer resp.Body.Close()
                body, err := ioutil.ReadAll(resp.Body)
                if err != nil {
                        fmt.Println("sfsdfadf")
                }else{
                        //restart_epoch, _ := strconv.Atoi(strings.Split(string(body), "")[1])
                        restart_epoch := strings.Split(string(body), " ")[5]
			for key, value := range restart_epoch {
				fmt.Printf("key %v, value %s \n", key, value)
			}
			ne, _:= strconv.Atoi(strings.Replace(restart_epoch[5], "\n", "", -1))
			// 此处取出的成员后面跟随了一个空格
			fmt.Println("restart.....")
                        fmt.Println(restart_epoch)
			fmt.Println(string(body))
			fmt.Println(restart_epoch[5])
			fmt.Println(ne)
                }
        }
}

func xhttpGet() {
	resp, err := http.Get("http://172.30.0.8:65534/server_info")
	if err != nil {
	}else{
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
		}else{
			restart_epoch, _ := strconv.Atoi((strings.Replace(strings.Split(string(body), " ")[5], "\n", "", -1)))
			fmt.Println(restart_epoch)
		}
	}
}

func main(){
        httpGet()
	xhttpGet()
}
