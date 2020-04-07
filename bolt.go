


package main

import (
    "golang.org/x/net/publicsuffix"
    "io/ioutil"
    "log"
    "net/http"
    "net/http/cookiejar"
    "net/url"
    "fmt"
    "golang.org/x/net/html"

)

func getElementById(id string, n *html.Node) (element *html.Node, ok bool) {
	for _, a := range n.Attr {
		if a.Key == "id" && a.Val == id {
			return n, true
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if element, ok = getElementById(id, c); ok {
			return
		}
	}
	return
}

func main() {

    options := cookiejar.Options{
        PublicSuffixList: publicsuffix.List,
    }
    jar, err := cookiejar.New(&options)
    if err != nil {
        log.Fatal(err)
    }
    client := http.Client{Jar: jar}
    resp, err := client.Get("http://test/index.php/bolt/login")
    if err != nil {
        log.Fatal(err)
    }
	defer resp.Body.Close()
	root, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	element, ok := getElementById("user_login__token", root)
	if !ok {
		log.Fatal("element not found")
	}
	for _, a := range element.Attr {
		if a.Key == "value" {
			fmt.Println("user_login__token="+a.Val)
			return
		}
	}
	log.Fatal("element missing value")


    resp, err = client.PostForm("http://127.0.0.1:8081/", url.Values{
     "user_login[username]": {"test"},
     "user_login[password]" : {"password"},
      })
    if err != nil {
        log.Fatal(err)
    }

    data, err := ioutil.ReadAll(resp.Body)
    resp.Body.Close()
    if err != nil {
        log.Fatal(err)
    }
    log.Println(string(data))
}

