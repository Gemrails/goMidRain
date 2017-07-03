package test_main

import (
	"fmt"
	"strconv"
        "midrain.app/midconst"
        "net/http"
        "io/ioutil"
        "strings"
)

func httpGet() {
        resp, err := http.Get(midconst.Infourl)
        if err != nil {
                fmt.Println("123123")
        }else{
                defer resp.Body.Close()
                body, err := ioutil.ReadAll(resp.Body)
                if err != nil {
                        fmt.Println("sfsdfadf")
                }else{
                        restart_epoch, _ := strconv.Atoi(strings.Split(string(body), " ")[5])
                        fmt.Println(restart_epoch)
                }
        }
}

func main(){
        httpGet
}
