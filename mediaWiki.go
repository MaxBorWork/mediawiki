package mediawiki

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)

func MediaWikiRequest(url string) string {
	resp := CreateHttp(url)

	doc, err :=	goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Println(err.Error())
	}

	text := doc.Find("pre").Text()
	var jsonBody map[string]interface{}
	if err := json.Unmarshal([]byte(text), &jsonBody); err != nil {
		panic(err)
	}

	return ParseJsonResponse(jsonBody)
}

func ParseJsonResponse(jsonBody map[string]interface{}) string {
	query := jsonBody["query"].(map[string]interface{})
	pages := query["pages"].(map[string]interface{})

	var pagesVal map[string]interface{}
	for key, _ := range pages {
		pagesVal = pages[key].(map[string]interface{})
	}

	fullText := pagesVal["extract"].(string)

	return TrimText(fullText)
}

func TrimText(text string) string {
	var smallText string
	textArr := strings.Fields(text)
	for i, word := range textArr {
		if i < 200 {
			smallText = smallText + " " + word
		} else {
			break
		}
	}

	smallText = smallText + "."

	return smallText
}

func CreateHttp(url string) *http.Response {
	req, e := http.NewRequest("GET", url, nil)
	if e != nil {
		log.Fatal(e)
	}

	res, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		log.Fatal(respErr)
	}

	return res
}

func CreateUrl(lang, title string) string {
	return "https://" + lang + ".wikipedia.org/w/api.php?fomat=json&prop=extracts&action=query&explaintext&titles=" + title
}
