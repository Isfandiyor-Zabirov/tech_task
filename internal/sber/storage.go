package sber

import (
	"bytes"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"sync"
)

const (
	mainUrl        = "https://www.sberbank-ast.ru"
	sendRequestUrl = "https://www.sberbank-ast.ru/UnitedPurchaseList.aspx"
)

func sendRequest(searchKey string, client *http.Client, cookies []*http.Cookie, wg *sync.WaitGroup, channel chan Response) {
	defer wg.Done()
	var (
		body = SendRequestBody{Search: searchKey}
		resp = Response{
			Keyword: searchKey,
		}
	)

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		log.Println("sendRequest func json marshal body error:", err.Error())
		return
	}

	searchReq, err := http.NewRequest("POST", sendRequestUrl, bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Println("sendRequest func create new http request error:", err.Error())
		return
	}

	for _, cookie := range cookies {
		searchReq.AddCookie(cookie)
	}

	response, err := client.Do(searchReq)
	if err != nil {
		log.Println("sendRequest func client.Do error:", err.Error())
		return
	}

	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Println("get doc error:", err.Error())
		return
	}

	// get info about the first auction in the list
	firstDiv := doc.Find("div#resultTable")
	firstTable := firstDiv.Find("table.es-reestr-tbl.its:nth-child(1)")
	tbody := firstTable.Find("tbody")
	nameSpan := tbody.Find("tr").First().Find("td:nth-child(2)").
		Find("span.es-el-name")
	numberSpan := tbody.Find("tr").First().Find("td:nth-child(2)").
		Find("span.es-el-code-term")

	var auction Auction
	auction.AuctionName = nameSpan.Text()
	auction.AuctionNumber = numberSpan.Text()
	resp.Auctions = append(resp.Auctions, auction)

	// get info about the second auction in the list
	secondTable := firstDiv.Find("table.es-reestr-tbl.its:nth-child(2)")
	tbody = secondTable.Find("tbody")
	nameSpan = tbody.Find("tr").First().Find("td:nth-child(2)").
		Find("span.es-el-name")
	numberSpan = tbody.Find("tr").First().Find("td:nth-child(2)").
		Find("span.es-el-code-term")

	auction.AuctionName = nameSpan.Text()
	auction.AuctionNumber = numberSpan.Text()
	resp.Auctions = append(resp.Auctions, auction)

	channel <- resp
}
