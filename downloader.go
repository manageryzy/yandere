package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"flag"
	"strconv"
	"github.com/PuerkitoBio/goquery"
)

func HTTPDownload(uri string) ([]byte, error) {
	//fmt.Printf("HTTPDownload From: %s.\n", uri)
	res, err := http.Get(uri)
	if err != nil {
		println(err)
	}
	defer res.Body.Close()
	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		println(err)
	}
	//fmt.Printf("ReadFile: Size of download: %d\n", len(d))
	return d, err
}

func WriteFile(dst string, d []byte) error {
	fmt.Printf("WriteFile: Size of download: %d\n", len(d))
	err := ioutil.WriteFile(dst, d, 0444)
	if err != nil {
		println(err)
	}
	return err
}

func DownloadToFile(uri string, dst string) {
	//fmt.Printf("DownloadToFile From: %s.\n", uri)
	if d, err := HTTPDownload(uri); err == nil {
		//fmt.Printf("downloaded %s.\n", uri)
		if WriteFile(dst, d) == nil {
			//fmt.Printf("saved %s as %s\n", uri, dst)
		}
	}
}

var startID,endID,workers int;
var jobs = make([]int,0)
var c chan int;

func init() {
	flag.IntVar(&startID,"start",1,"start number")
	flag.IntVar(&endID,"end",355850,"end number")
	flag.IntVar(&workers,"workers",4,"download thread number")

	flag.Parse()
}

func main() {
	os.MkdirAll("image",666);


	for i:=startID;i<=endID;i++{
		jobs = append(jobs,i)
	}

	for i:=0;i<workers;i++{
		go downloadWorker()
	}

	for i := 0; i <workers; i++ {
		<-c
	}
}

func downloadWorker(){
	defer func() {
		c<-0;
		println("worker exit")
	}()

	x := 0
	x, jobs = jobs[len(jobs)-1], jobs[:len(jobs)-1]

	downloadImage(x)
}


func downloadImage(id int){
	println("querying " + strconv.Itoa(id))
	doc, err := goquery.NewDocument("https://yande.re/post/show/"+strconv.Itoa(id))
	if err != nil {
		println(err)
	}

	doc.Find("#highres").Each(func(i int, s *goquery.Selection) {
		href,exist := s.Attr("href")
		if exist{
			println("downloading:"+strconv.Itoa(id)+"\t@"+href)
			DownloadToFile(href,"./image/"+strconv.Itoa(id)+".jpg")
		}else {
			println("fail to download "+strconv.Itoa(id))
		}
	})
}