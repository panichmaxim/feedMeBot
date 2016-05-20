package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"net/url"
	"bytes"
	"encoding/json"
)

func post(urlRes string) {
      var b []byte  
       data := url.Values{}
       data.Set("city_id", "1")
       data.Add("place_id", "61")
       data.Add("customer_phone", "+789989123")
       data.Add("customer_guests", "5")
       data.Add("date", "2015-04-03")
  
    request, _ := http.NewRequest("POST", urlRes, bytes.NewBufferString(data.Encode()))
        request.Header.Add("Content-Type", "application/json")
      client := &http.Client{}
      resp, _ := client.Do(request)
        b, _ = ioutil.ReadAll(resp.Body)
        fmt.Println("Response body:", string(b))
}

func get (url string, id string , model interface{}){
  response, err := http.Get(url+id)
    if err != nil {
        fmt.Printf("%s", err)
    } else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
          //  fmt.Printf("%s", err)
        }
       json.Unmarshal(contents, &model)
      // fmt.Println(model)
      // fmt.Printf("%s\n", string(contents))
    }
}