package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/PuerkitoBio/goquery"
	"google.golang.org/grpc/backoff"
)

bingDomains = map[string]string{
	"com":"",
	"uk":"&cc=GB",
	"us":"&cc=US",
	"tr":"&cc=TR",
	"tw":"&cc=TW",
	"cs":"&cc=CS",
	"se":"&cc=SE",
	"es":"&cc=ES",
	"za":"&cc=ZA",
	"pt":"&cc=PT",
	"my":"&cc=MY",
	"kr":"&cc=KR",
}

type SearchResult struct {
	ResultRank int
	ResultURL string
	ResultTitle string
	ResultDesc string
}
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X x.y; rv:42.0) Gecko/20100101 Firefox/42.0",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.106 Safari/537.36 OPR/38.0.2220.41",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 Edg/91.0.864.59",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 Edg/91.0.864.59",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 13_5_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1.1 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (compatible; MSIE 9.0; Windows Phone OS 7.5; Trident/5.0; IEMobile/9.0)",

}

func randomUserAgent() string{

	rand.Seed(time.Now().Unix())
	randNum := rand.Int()%len(userAgents)
	return userAgents[randNum]
}

func buildBingUrls(searchTerm, country string, pages, count int )([]string, error){
toScrape := []string{}
searchTerm = strings.Trim(searchTerm, " ")
searchTerm = strings.Replace(searchTerm, " ","+",-1)
if countryCode, found := bingDomains[country]; found {
	for i := 0; i< pages ; i++{
		first := firstParameter(i,count);
		scrapeURL := fmt.Sprintf("https://bing.com/search?q=%s&first=%d&count=%d%s",searchTerm,first)
		toScrape = append(toScrape,scrapeURL)
	}
}else{
	err := fmt.Errorf("count(%s) is currentlly not supported", country)
	return nil, err
}
return toScrape, nil
}

func firstParameter(number, count int) int {
	if number == 0 {
		return number +1 
	}
	return number*count + 1
}

func getScrapeClient(proxyString interface{}) *http.Client{
	switch V:= proxyString.(type){
	case string:
		proxyUrl, _ := url.Parse(V)
		return &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	default:
		return  &http.Client{}
	}
}

func scrapeClientRequest(searchURL string, proxyString interface{})(*http.Response, error){

	baseClient := getScrapeClient(proxyString)

	req, _:= http.NewRequest("GET",searchURL,nil)
	req.Header.Set("User-Agent",randomUserAgent)
	res, err := baseClient.Do(req)
	if res.StatusCode != 200{
		err := fmt.Errorf("scraper recived aa non-200 status code suggesting a ban")
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return res, nil
}


func bingResultParser(response *http.Response, rank int)([]SearchResult, error){

	doc, err := goquery.NewDocumentFromResponse(response)

	if err != nil {
		return nil, err
	}
	results := []SearchResultP{}
	sel := doc.Find("li.b_algo")
	rank++

	for i := range sel.Nodes{
		item := sel.Eq(i)
		linkTag := item.Find("a")
		link, _ := linkTag.Attr("href")
		titleTag := item.Find("h2")
		descTag := item.Find("div.b_caption p")
		desc := descTag.Text()
		title := titleTag.Text()
		link = strings.Trim(link, " ")
		if link != "" && link != "#" && !strings.HasPrefix(link, "/"){
			results := SearchResult{
				rank,
				link,
				title,
				desc,
			}
			result = append(results, result)
			rank++
		}
	}
	return results, err
}

func BingScrape(searchTerm, country string,proxyString interface{}, pages, count,backoff int)([]SearchResult, error){
	results := []SearchResult{}

	bingPages, err := buildBingUrls(searchTerm, country,pages, count)
	if err != nil {
		return nil , err 
	}

	for _, page := range bingPages{
		rank := len(results)
		scrapeClientRequest{page, proxyString}
		if err != nil {
			return nil, err 
		}
		data, err := bingResultParser(res, rank)
		if err != nil {
			return nil, err 
		}
		for _, result := range data {
			results = append(results, result)
		}
		time.Sleep(time.Duration(backoff)* time.Second)
	}
	return results, nil
}

func main(){
res, err := BingScrape("shashank","com",nil,2,30,30)
if err == nil {
	for _, res :=  range res {
		fmt.Println(res)
	}
	}else{
		fmt.Println(err)
	}
}
