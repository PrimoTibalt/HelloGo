package main

import (
	"errors"
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
	request, createRequestError := http.NewRequest("GET",
		"https://novelfire.net/genre-all/sort-new/status-all/all-novel", nil)
	if createRequestError != nil {
		panic(createRequestError)
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")
	request.Header.Add("Referer", "https://novelfire.net")
	response, getAllNovelsErr := client.Do(request)
	if getAllNovelsErr != nil {
		panic(getAllNovelsErr)
	}

	if response.StatusCode != 200 {
		panic(errors.New("Oops, you have got " + strconv.Itoa(response.StatusCode) + " as a response"))
	}
	defer response.Body.Close()

	resultHTML, parsingError := html.Parse(response.Body)
	if parsingError != nil {
		panic(parsingError)
	}

	titles := map[string]string{}
	wg.Add(1)
	searchInDoc(resultHTML, titles)
	wg.Wait()
	for title, href := range titles {
		fmt.Println(title, href)
	}
}

func searchInDoc(doc *html.Node, result map[string]string) {
	for _, attribute := range doc.Attr {
		if attribute.Key == "class" && attribute.Val == "novel-item" {
			var href string
			for _, attr := range doc.FirstChild.Attr {
				if attr.Key == "href" {
					href = attr.Val
				}
			}

			if href == "" {
				break
			}

			for _, attr := range doc.FirstChild.Attr {
				if attr.Key == "title" {
					rwl.Lock()
					result[attr.Val] = href
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
