package sdext

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func getTitle(doc *html.Node) ALink {
	var p ALink
	var linkReady bool
	var title string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			title = title + n.Data
		}

		for _, a := range n.Attr {
			if a.Key == "id" {
				p.Id = strings.TrimPrefix(a.Val, "title-")
			}

			if a.Key == "class" && a.Val == "story-title" {
				linkReady = true
			}

			if a.Key == "href" && linkReady {
				p.Link = strings.TrimPrefix(a.Val, "//")
				linkReady = false
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	p.Title = title
	return p
}

func parseArticle(doc *html.Node, id string, smap *map[string]string) {
	selectmap := *(smap)
	var f func(*html.Node)
	f = func(n *html.Node) {
		for _, a := range n.Attr {
			//fmt.Println("key:", a.Key, a.Val)
			if a.Key == "id" && a.Val == "text-"+id {
				var buf bytes.Buffer
				html.Render(&buf, n)
				selectmap["article"] = buf.String()
			}

		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
}

func parseLink(doc *html.Node) (links []ALink) {
	var f func(*html.Node)
	f = func(n *html.Node) {

		for _, a := range n.Attr {
			if a.Key == "class" && a.Val == "story-title" {
				link := getTitle(n)
				links = append(links, link)
			}

		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return links
}

type ALink struct {
	Link               string
	Id, Title, Content string
}

type UrlDone struct {
	Url string
}

func extractLinks(baseUrl string) []ALink {
	rsp, err := http.Get(baseUrl)
	if err != nil {
		return []ALink{}
	}

	doc, err := html.Parse(rsp.Body)
	return parseLink(doc)
}

func extractContent(url, id string) string {
	var article string
	rsp, err := http.Get(url)
	if err != nil {
		return ""
	}

	doc, err := html.Parse(rsp.Body)
	var selectmap = make(map[string]string)
	parseArticle(doc, id, &selectmap)

	for k, v := range selectmap {
		if k == "article" {
			article = v
		}
	}
	return article
}

var urlDoneMap = make(map[string]bool)

func loadUrlDone(dataPath string) {
	fi, err := os.Open(dataPath + "/.data.json")
	defer fi.Close()
	if err != nil {
		log.Println(err)
	}
	data, err := ioutil.ReadAll(fi)
	slice := strings.Split(string(data), "\n")
	for _, l := range slice {
		var ud []UrlDone
		err := json.Unmarshal([]byte(l), &ud)
		if err != nil {
			log.Println(err)
			continue
		}
		for _, u := range ud {
			log.Println("Load:", u.Url)
			urlDoneMap[u.Url] = true
		}
	}
}

func Extracter(startPage, dataPath string) (res []ALink) {
	loadUrlDone(dataPath)
	links := extractLinks(startPage)
	for _, a := range links {
		if !urlDoneMap[a.Link] {
			a.Content = extractContent("https://"+a.Link, a.Id)
			res = append(res, a)
		}

	}
	return res
}
