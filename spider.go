package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func getHtml(url string) []byte{
	//建立自定义request
	client:= &http.Client{}
	req,err:=http.NewRequest("GET",url,nil)
	if err!=nil{
		log.Println(err)
	}
	req.Header.Add("User-Agent","Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.183 Safari/537.36")
	req.Header.Add("Referer",url)
	resp,err:=client.Do(req)
	if err!=nil{
		log.Println(err)
	}
	defer resp.Body.Close()
	body,err:=ioutil.ReadAll(resp.Body)
	if err!=nil{
		log.Println(err)
	}
	html:=body
	return html
}
func launchQuery(html io.Reader,fTitles *os.File,fTimes *os.File,fTags *os.File,fBriefs *os.File){
	dom,err:=goquery.NewDocumentFromReader(html)
	if err!=nil{
		log.Fatalln(err)
	}
	titles:=dom.Find(".post-title a").Map(func(i int, selection *goquery.Selection) string {
		return selection.Text()
	})
	times:=dom.Find("time[datetime]").Map(func(i int, selection *goquery.Selection) string {
		return selection.Text()
	})
	tags:=dom.Find("a.post-meta-tag").Map(func(i int, selection *goquery.Selection) string {
		return selection.Text()
	})
	briefs:=dom.Find("p.post-excerpt").Map(func(i int, selection *goquery.Selection) string {
		return  selection.Text()
	})

	for _,v:=range titles{
		fTitles.WriteString(v+"\n\n")
	}
	for _,v:=range times{
		fTimes.WriteString(v+"\n\n")
	}
	for _,v:=range tags{
		fTags.WriteString(v+"\n\n")
	}
	for _,v:=range briefs{
		fBriefs.WriteString(v+"\n\n")
	}
	fmt.Println(titles)
	fmt.Println()
	fmt.Println(times)
	fmt.Println()
	fmt.Println(tags)
	fmt.Println()
	fmt.Println(briefs)
	fmt.Println()
}

func main() {
	fTitles,err:=os.Create("./data/title.txt")
	if err!=nil{
		log.Fatalln(err)
	}
	fTimes,err:=os.Create("./data/times.txt")
	if err!=nil{
		log.Fatalln(err)
	}
	fTags,err:=os.Create("./data/tags.txt")
	if err!=nil{
		log.Fatalln(err)
	}
	fBriefs,err:=os.Create("./data/brief.txt")
	if err!=nil{
		log.Fatalln(err)
	}
	for i:=1;i<8;i++ {
		fmt.Printf("现在是第%v页内容\n",i)
		body := getHtml("https://blog.lenconda.top/page/"+strconv.Itoa(i))
		html := bytes.NewReader(body)
		launchQuery(html,fTitles,fTimes,fTags,fBriefs)
	}
	defer fTitles.Close()
	defer fTimes.Close()
	defer fTags.Close()
	defer fBriefs.Close()
}
