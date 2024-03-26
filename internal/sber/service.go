package sber

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"tech_task/utils"
)

func SendRequest(filePath string) (err error) {
	var (
		client  = &http.Client{}
		resp    []Response
		cookies []*http.Cookie
		channel = make(chan Response)
	)

	keywords, err := utils.ReadFile(filePath)
	if err != nil {
		return err
	}

	mainReq, err := http.NewRequest("GET", mainUrl, nil)
	if err != nil {
		log.Println("SendRequest func main url request error:", err.Error())
		return
	}

	mainResp, err := client.Do(mainReq)
	if err != nil {
		log.Println("SendRequest func main url client.Do error:", err.Error())
		return
	}

	defer mainResp.Body.Close()

	cookies = mainReq.Cookies()

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	for _, keyword := range keywords {
		wg.Add(1)
		mu.Lock()
		go sendRequest(keyword, client, cookies, &wg, channel)
		mu.Unlock()
		resp = append(resp, <-channel)
	}

	wg.Wait()

	resBytes, _ := json.Marshal(resp)

	log.Println("resp:", string(resBytes))

	return nil
}
