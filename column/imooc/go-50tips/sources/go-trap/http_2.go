package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

var sites = []string{
	"https://tip.golang.org",
	"https://www.oracle.com/java",
	"https://python.org",
}

func main() {
	var wg sync.WaitGroup
	wg.Add(len(sites))

	for _, site := range sites {
		site := site
		go func() {
			defer wg.Done()
			req, err := http.NewRequest("GET", site, nil)
			if err != nil {
				fmt.Println(err)
				return
			}
			req.Close = true

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("get response from %s, resp length = %d\n", site, len(body))
		}()
	}
	wg.Wait()
}
