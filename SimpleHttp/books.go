package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"golang.org/x/net/html"
)

var (
	rwl sync.RWMutex
	wg  sync.WaitGroup
)

func main() {
	client := &http.Client{}
	request, error := http.NewRequest("GET",
		"https://novelfire.net/genre-all/sort-new/status-all/all-novel", nil)
	if error != nil {
		fmt.Println(error.Error())
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")
	request.Header.Add("Referer", "https://novelfire.net")
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
	}

	if response.StatusCode != 200 {
		fmt.Println("Oops, you have got " + strconv.Itoa(response.StatusCode) + " as a response")
	}
	defer response.Body.Close()
	resultHTML, parsingError := html.Parse(response.Body)
	if parsingError != nil {
		fmt.Println(parsingError.Error())
	}

	titles := &[]string{}
	wg.Add(1)
	searchInDoc(resultHTML, titles)
	wg.Wait()
	for _, title := range *titles {
		fmt.Println(title)
	}
}

func searchInDoc(doc *html.Node, result *[]string) {
	for _, attribute := range doc.Attr {
		if attribute.Key == "class" && attribute.Val == "novel-item" {
			for _, attr := range doc.FirstChild.Attr {
				if attr.Key == "title" {
					rwl.Lock()
					*result = append(*result, attr.Val)
					rwl.Unlock()
					wg.Done()
					return
				}
			}
		}
	}

	for nextDoc := range doc.ChildNodes() {
		wg.Add(1)
		go searchInDoc(nextDoc, result)
	}

	wg.Done()
}
