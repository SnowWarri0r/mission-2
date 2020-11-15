package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	html2 "html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)
var data string
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
func launchQuery(html io.Reader){
	dom,err:=goquery.NewDocumentFromReader(html)
	if err!=nil{
		log.Fatalln(err)
	}
        dom.Find("div.box.post-box").Each(func(i int, selection *goquery.Selection) {
        	title:=selection.Find("h2.post-title a").Text()//将数据分类为4种
        	var tags string
        	selection.Find("a.post-meta-tag").Each(func(i int, selection *goquery.Selection) {
				tags=tags+" "+selection.Text()
			})
        	time:=selection.Find("div.box.post-box time[datetime]").Text()
        	brief:=selection.Find("div.box.post-box p.post-excerpt").Text()
        	brief=html2.EscapeString(brief)
        	data=data+"<h2>"+title+"</h2>\n<p>"+time+"&nbsp;&nbsp;"+tags+"</p>\n<p>"+brief+"</p>\n\n"
		})
}
func main() {
	s:=&http.Server{
		Addr: ":8080",//设置服务器映射端口为8080
	}
	for i:=1;i<9;i++ {
		body := getHtml("https://blog.lenconda.top/page/"+strconv.Itoa(i))
		html := bytes.NewReader(body)
		launchQuery(html)
	}
	f,err:=os.Create("./data/list.txt")//将抓取数据写入文件
	if err!=nil{
		log.Fatalln(err)
	}
	f.WriteString(data)
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html;charset=UTF-8")//设置数据显示方式为html，编码保存为utf-8
		fmt.Fprintln(writer,data)
	})
	s.ListenAndServe()
	defer f.Close()
}
