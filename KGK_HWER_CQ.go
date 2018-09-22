package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/cihub/seelog"
	"github.com/satori/go.uuid"
	"github.com/valyala/fasthttp"
	//"os"
	"sync"
	"time"
	//"os/exec"
	//"reflect"
	//"os"
	//"os"
	// "os"
)

var g_lock sync.Mutex
var g_lock_m_ch sync.Mutex
var task_queue = NewQueue()

type Job struct {
	Class string        `json:"Class"`
	Args  []interface{} `json:"Args"`
}

type Queue struct {
	data *list.List
}

func NewQueue() *Queue {
	q := new(Queue)
	q.data = list.New()
	return q
}

func (q *Queue) push(v interface{}) {
	defer g_lock.Unlock()
	g_lock.Lock()
	q.data.PushFront(v)
}

func (q *Queue) pop_safe() interface{} {
	g_lock.Lock()
	if q.data.Len() > 0 {
		iter := q.data.Back()
		v := iter.Value
		q.data.Remove(iter)
		g_lock.Unlock()
		return v
	} else {
		g_lock.Unlock()
		return nil
	}
}

func (q *Queue) remove(v interface{}) {
	defer g_lock.Unlock()
	g_lock.Lock()
	for iter := q.data.Back(); iter != nil; iter = iter.Prev() {
		var task Task
		value_1, ok_1 := iter.Value.([]byte)
		if ok_1 {
			json.Unmarshal(value_1, &task)
		}

		var task_v Task
		value_2, ok_2 := v.([]byte)
		if ok_2 {
			json.Unmarshal(value_2, &task_v)
		}

		if task.Task_uuid == task_v.Task_uuid {
			q.data.Remove(iter)
			break
		}
	}
}

func (q *Queue) dump() {
	for iter := q.data.Back(); iter != nil; iter = iter.Prev() {
		fmt.Println("item:", iter.Value)
	}
}

func (q *Queue) len() int {
	defer g_lock.Unlock()
	g_lock.Lock()
	return q.data.Len()
}

type Req_body struct {
	Id      float64 `json:"id"`
	Img_url string  `json:"img_url"`
}

type Task struct {
	Task_uuid string  `json:"task_uuid"`
	Client_id float64 `json:"client_id"`
	Img_url   string  `json:"img_url"`
	//Maxfaceonly int `json:"maxfaceonly"`
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type Result_body_HWER struct {
	Task_uuid string `json:"task_uuid"`
	//Result    HWER   `json:"result"`
	Result string `json:"result"`
}

//type HWER struct {
//	Activation_idx []int      `json:"activation_idx"`
//	Labels         [][]string `json:"labels"`
//	Latex          string     `json:"latex"`
//	Error          int        `json:"error"`
//}

//type Syms struct {
//	Sym_list []string `json:"sym_list"`
//}

var code_map map[int]string
var m_ch map[string]chan Result_body_HWER

//var g_time_out = 10

func handler_upload_task(ctx *fasthttp.RequestCtx) {
	t1 := time.Now()
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")

	var req_body Req_body
	body := ctx.PostBody()
	err := json.Unmarshal(body, &req_body)
	if err != nil {
		fmt.Println("Fail to parse req body")
		fmt.Println(err)
	}

	var task Task
	//fmt.Println(req_body.Pt_seq)
	u, err := uuid.NewV4()
	uuid_str := u.String()
	task.Task_uuid = fmt.Sprintf("%s", uuid_str)
	task.Client_id = req_body.Id
	task.Img_url = req_body.Img_url
	//task.Maxfaceonly = req_body.Maxfaceonly
	task_queue.push(task)

	var res Response

	c := make(chan Result_body_HWER)
	g_lock_m_ch.Lock()
	m_ch[task.Task_uuid] = c
	g_lock_m_ch.Unlock()

	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(10e9)
		timeout <- true
	}()

	select {
	case result := <-c:
		res.Code = 0
		res.Msg = code_map[res.Code]

		in := []byte(result.Result) // result.Result是srv的result字段的值，是python的json.dumps(.)的内容（非普通的string，是raw string）；原样转发srv返回的识别结果（json字段和值），然后以json的格式返回给客户端
		var raw map[string]interface{}
		json.Unmarshal(in, &raw)
		raw["task_uuid"] = result.Task_uuid
		res.Data = raw
	case <-timeout:
		res.Code = 1004
		res.Msg = code_map[res.Code]
		res.Data = nil
		//fmt.Println("Timeout")
	}

	//fmt.Println("Got task result")
	task_queue.remove(task)
	g_lock_m_ch.Lock()
	delete(m_ch, task.Task_uuid)
	g_lock_m_ch.Unlock()

	data, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
	}
	ctx.Write(data)

	elapsed := time.Since(t1)
	logger.Debug("Task UUID: ", task.Task_uuid, " Time: ", elapsed)
	logger.Flush()
}

func handler_get_task(ctx *fasthttp.RequestCtx) {
	t1 := time.Now()
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	var res Response
	var uuid string
	task := task_queue.pop_safe()

	if task == nil {
		res.Code = 2000
		res.Msg = code_map[res.Code]
		res.Data = nil

	} else {
		tmp, ok := task.(Task)
		if ok {
			res.Code = 0
			res.Msg = code_map[res.Code]
			res.Data = tmp
			uuid = tmp.Task_uuid
		}
	}

	data, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
	}
	ctx.Write(data)

	elapsed := time.Since(t1)
	logger.Debug(uuid, "time:", elapsed)
	logger.Flush()
}

func Handler_upload_result(ctx *fasthttp.RequestCtx) {
	t1 := time.Now()
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")

	var v Result_body_HWER
	body := ctx.PostBody()
	err := json.Unmarshal(body, &v)
	if err != nil {
		fmt.Println("Fail to parse req body in result post.")
		fmt.Println(err)
	}

	//fmt.Println("uuid: " + v.Task_uuid)

	g_lock_m_ch.Lock()
	c, ok := m_ch[v.Task_uuid]
	g_lock_m_ch.Unlock()

	var response Response
	if ok {
		c <- v
		response.Code = 0
		response.Msg = code_map[response.Code]
		response.Data = nil
	} else {
		response.Code = 2004
		response.Msg = code_map[response.Code]
		response.Data = nil
	}

	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("receive2:", v)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Write(data)

	elapsed := time.Since(t1)
	logger.Debug(v.Task_uuid, "time:", elapsed)
	logger.Flush()
}

func create_code_map() {
	code_map[0] = "success."
	code_map[1004] = "Request Time Out. Server is busy..."
	code_map[2000] = "No Task."
	code_map[2004] = "Task Time Out."
}

var logger seelog.LoggerInterface

func main() {
	var port string
	// port = os.Args[1]   # 命令行参数的形式输入端口号
	port = "3611"
	fmt.Printf("Launch server at %s. $_$\n", port)

	var err error
	logger, err = seelog.LoggerFromConfigAsFile("seelog_HWER.xml") //  different business using diff xml config.
	if err != nil {
		seelog.Critical("err parsing config log file", err)
		return
	}
	logger.Info("Launch server on " + port + ". $_$")
	logger.Flush()

	code_map = make(map[int]string)
	m_ch = make(map[string]chan Result_body_HWER)
	create_code_map()

	router := fasthttprouter.New()
	router.POST("/upload_task/", handler_upload_task)
	router.GET("/get_task/", handler_get_task)
	router.POST("/upload_result/", Handler_upload_result)

	if err := fasthttp.ListenAndServe("0.0.0.0:"+port, router.Handler); err != nil {
		fmt.Println("start fasthttp fail:", err.Error())
	}
}
