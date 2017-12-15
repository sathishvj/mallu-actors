package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Cast struct {
	id     string
	title  string
	actors []string
}

var req = []string{
	"Mohanlal",
	"Thilakan",
	"Nedumudi Venu",
	"Jagathy Sreekumar",
	"Innocent",
}

func ScrapeMovie(id string) {

	cast := Cast{}
	link := "http://malayalasangeetham.info/m.php?" + id
	doc, err := goquery.NewDocument(link)
	if err != nil {
		log.Fatal(err)
	}

	heading := "div.pheading"
	doc.Find(heading).Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		fmt.Println(id, ":", title)
		cast.id = id
		cast.title = title
	})

	// Find the actor items
	pattern := "td.prowsshort"
	//doc.Find(".sidebar-reviews article .content-block").Each(func(i int, s *goquery.Selection) {
	doc.Find(pattern).Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		//band := s.Find("a").Text()
		//title := s.Find("i").Text()
		//fmt.Printf("%d: %s - %s\n", i, band, title)

		title := s.Text()
		if strings.Index(title, "അഭിനേതാക്ക") >= 0 {
			s2 := s.Siblings()
			// fmt.Printf("s2 length: %d\n", s2.Length())

			subpattern := "a"

			s2.Find(subpattern).Each(func(j int, s3 *goquery.Selection) {
				//fmt.Printf("%+v\n", s3)
				data, _ := s3.Attr("href")
				data = strings.Replace(data, "displayProfile.php?category=actors&artist=", "", -1)
				fmt.Printf("\t%d: %s\n", j+1, data)
				cast.actors = append(cast.actors, data)
			})
		}
	})

	if hasActors(cast, req) {
		fmt.Printf("Found: %+v\n", cast)
	}
}

func hasActors(c Cast, req []string) bool {

	for _, r := range req {
		found := false
		for _, a := range c.actors {
			if r == a {
				found = true
			}
		}
		if !found {
			return false
		}
	}

	return true
}

func ScrapeMovieList() {

	// There are 14 pages of movies on the site
	for id := 0; id < 14; id++ {
		link := "http://malayalasangeetham.info/movies.php?tag=Search&actor=Mohanlal&limit=346&alimit=87&page_num=" + strconv.Itoa(id+1)
		doc, err := goquery.NewDocument(link)
		if err != nil {
			log.Fatal(err)
		}

		pattern := "a"
		//doc.Find(".sidebar-reviews article .content-block").Each(func(i int, s *goquery.Selection) {
		doc.Find(pattern).Each(func(i int, s *goquery.Selection) {
			data, _ := s.Attr("href")
			str := "m.php?"
			if strings.Index(data, str) == 0 {
				id := strings.Replace(data, str, "", -1)
				fmt.Println(id)
			}
		})
	}

}

func main() {
	// first scrape movie list ids and manually add it to movielist.go
	//ScrapeMovieList()

	// then run scraper on individual movie list
	for i := 0; i < len(ids); i++ {
		ScrapeMovie(ids[i])
	}
}
