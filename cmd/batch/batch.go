package main

import (
	"log"
	"net/http"
	"regexp"
	"strings"
	"ynp/api"
	"ynp/env"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	env.MyDB = api.NewDb()
	defer env.MyDB.Close()
	crawlingLottoNum()
}

func crawlingLottoNum() error {
	lottoUrl := "https://dhlottery.co.kr/gameResult.do?method=byWin"

	res, err := http.Get(lottoUrl)
	if err != nil {
		log.Println(err)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Println("http.Status not OK")
		return err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	win := doc.Find("div.num.win").Contents().FilterFunction(func(i int, s *goquery.Selection) bool {
		return !s.Is("strong")
	}).Text()
	bonus := doc.Find("div.num.bonus").Contents().FilterFunction(func(i int, s *goquery.Selection) bool {
		return !s.Is("strong")
	}).Text()
	round := doc.Find("div.win_result > h4 > strong").Contents().FilterFunction(func(i int, s *goquery.Selection) bool {
		return !s.Is("strong")
	}).Text()

	winlist := strings.Split(strings.TrimSpace(win), "\n")

	var nums []string
	for i := 0; i < 6; i++ {
		nums = append(nums, strings.TrimSpace(winlist[i]))
	}
	nums = append(nums, strings.TrimSpace(bonus))

	re := regexp.MustCompile(`\d+`)
	r := re.FindAllString(round, -1)

	api.InsertNums(env.MyDB, strings.Join(r, ""), strings.Join(nums, " "))
	return nil
}
