package main

import (
    "encoding/json"
    "fmt"
    //"reflect"
)

func main() {
    //in := []byte(`{ "votes": { "option_A": "3" } }`)
    //fmt.Println(reflect.TypeOf(`{ "votes": { "option_A": "3" } }`))
    //var raw map[string]interface{}
    //json.Unmarshal(in, &raw)
    //raw["count"] = 1
    //out, _ := json.Marshal(raw)
    //println(string(out))
    str:= "{\"access_token\":\"uAUS6o5g-9rFWjYt39LYa7TKqiMVsIfCGPEN4IZzdAk5-T-ryVhL7xb8kYciuU_m\",\"expires_in\":7200}"
    var dat map[string]interface{}
    if err := json.Unmarshal([]byte(str), &dat); err == nil {
        fmt.Println(dat)
        fmt.Println(dat["expires_in"])
    } else {
        fmt.Println(err)
    }
}