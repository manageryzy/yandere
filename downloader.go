package main

import (
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
	"os"
	"flag"
	"strconv"
	"github.com/PuerkitoBio/goquery"
)

func HTTPDownload(uri string) ([]byte, error) {
	fmt.Printf("HTTPDownload From: %s.\n", uri)
	res, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ReadFile: Size of download: %d\n", len(d))
	return d, err
}

func WriteFile(dst string, d []byte) error {
	fmt.Printf("WriteFile: Size of download: %d\n", len(d))
	err := ioutil.WriteFile(dst, d, 0444)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func DownloadToFile(uri string, dst string) {
	fmt.Printf("DownloadToFile From: %s.\n", uri)
	if d, err := HTTPDownload(uri); err == nil {
		fmt.Printf("downloaded %s.\n", uri)
		if WriteFile(dst, d) == nil {
			fmt.Printf("saved %s as %s\n", uri, dst)
		}
	}
}

var startID,endID int;

func init() {
	flag.IntVar(&startID,"start",1,"start number")
	flag.IntVar(&endID,"end",355850,"end number")
}

func main() {
	os.MkdirAll("image",666);
	for i:=startID;i<=endID;i++{
		downloadImage(i)
	}
}


func downloadImage(id int){
	doc, err := goquery.NewDocument("https://yande.re/post/show/"+strconv.Itoa(id))
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("#highres").Each(func(i int, s *goquery.Selection) {
		href,exist := s.Attr("href")
		if exist{
			println("downloading:"+strconv.Itoa(id)+"@"+href)
			DownloadToFile(href,"./image/"+strconv.Itoa(id)+".jpg")
		}else {
			println("fail to download "+strconv.Itoa(id))
		}
	})


}