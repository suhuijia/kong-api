package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	//"strings"
	"bytes"
	//"encoding/json"
)

var count = 1

func httpGet() {
	//resp, err := http.Get("http://face-attr.ai.nd.com.cn:8050/v1/fa_get_task/")
	// resp, err := http.Get("http://192.168.46.39:8000/v1/get_task")
	resp, err := http.Get("http://120.25.253.2:8050/v1/hwer_get_task")

	if err != nil {
		//fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(count, " ", string(body))
	count += 1
}

func httpPost() {
	//resp, err := http.Post("http://192.168.46.123:8004/upload_result/",
	//    "application/json",
	//    strings.NewReader("name=cjb"))
	req := `{'Task_uuid': ,'ativation_idx': [14, 24, 42, 50, 66], 'labels': [['x', 'X', '\\lambda', '\\blank', 'A'], ['+', 't', '\\div', 'H', '\\pm'], ['y', 'j', 'g', 'Y', 'U'], ['=', 'F', '\\neq', 'E', '-'], ['z', 'Z', 'E', 't', '2']]}`
	req_new := bytes.NewBuffer([]byte(req))
	body_type := "application/json;charset=utf-8"
	// resp, err := http.Post("http://192.168.46.123:8000/upload_result/", body_type, req_new)
	resp, err := http.Post("http://120.25.253.2:8050/v1/hwer_post_result", body_type, req_new)
	if err != nil {
		//fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(count, " ", string(body))
	count += 1
}

func main() {
	for {
		time.Sleep(1 * time.Millisecond)
		httpGet()
		// time.Sleep(1 * time.Millisecond)
		// httpPost()
	}
}
