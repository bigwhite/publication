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

	tr := &http.Transport{
		DisableKeepAlives: true,
	}
	cli := &http.Client{
		Transport: tr,
	}

	for _, site := range sites {
		site := site
		go func() {
			defer wg.Done()

			resp, err := cli.Get(site)
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
