package api

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
	"ynp/env"

	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo/v4"
)

type Number struct {
	Num1 string `json:"num1"`
	Num2 string `json:"num2"`
	Num3 string `json:"num3"`
	Num4 string `json:"num4"`
	Num5 string `json:"num5"`
	Num6 string `json:"num6"`
	NumB string `json:"numb"`
}

// @Summary Index
// @Description Index API
// @Accept json
// @Produce json
// @Param name path string true "name of the user"
// @Success 200
// @Router / [get]
func Index(c echo.Context) error {
	r := &Res{
		Status:   http.StatusOK,
		Response: "hello ynp",
	}
	return c.JSONPretty(http.StatusOK, r, " ")
}

// @Summary Random
// @Description get 6 of random numbers
// @Accept json
// @Produce json
// @Success 200 {object} Res
// @Router /random [get]
func Random(c echo.Context) error {
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	nums := make(map[int]bool)
	var rn int

	for len(nums) < 6 {
		rn = random.Intn(45) + 1
		nums[rn] = true
	}
	var ret []int
	for i := 1; i <= 45; i++ {
		_, exist := nums[i]
		if exist {
			ret = append(ret, i)
		}
	}

	r := &Res{
		Status:   http.StatusOK,
		Response: ret,
	}
	return c.JSONPretty(http.StatusOK, r, " ")
}

func CreateTable(c echo.Context) error {
	err := NewTable()
	if err != nil {
		return c.String(http.StatusMethodNotAllowed, "table create fail")
	}
	r := &Res{
		Status:   http.StatusOK,
		Response: "table created",
	}
	return c.JSONPretty(http.StatusOK, r, " ")
}

func CrawlingLottoNum(c echo.Context) error {
	lottoUrl := "https://dhlottery.co.kr/gameResult.do?method=byWin"

	round := c.Param("round")
	if round != "latest" {
		lottoUrl = lottoUrl + "&drwNo=" + round
	}

	res, err := http.Get(lottoUrl)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusMethodNotAllowed, "http.Get fail")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Println("http.Status not OK")
		return c.String(http.StatusMethodNotAllowed, "http Get fail(no 200)")
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

	winlist := strings.Split(strings.TrimSpace(win), "\n")

	var nums []string
	for i := 0; i < 6; i++ {
		nums = append(nums, strings.TrimSpace(winlist[i]))
	}
	nums = append(nums, strings.TrimSpace(bonus))

	InsertNums(env.MyDB, round, strings.Join(nums, " "))

	r := &Res{
		Status:   http.StatusOK,
		Response: nums,
	}
	return c.JSONPretty(http.StatusOK, r, " ")
}

func CrawlingLottoNumAll(c echo.Context) error {
	go func() {
		lu := "https://dhlottery.co.kr/gameResult.do?method=byWin"

		for round := 1; round < 1098; round++ {

			lottoUrl := lu + "&drwNo=" + strconv.Itoa(round)

			res, err := http.Get(lottoUrl)
			if err != nil {
				log.Println(err)
				continue
			}
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				log.Println("http.Status not OK")
				continue
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

			winlist := strings.Split(strings.TrimSpace(win), "\n")

			var nums []string
			for i := 0; i < 6; i++ {
				nums = append(nums, strings.TrimSpace(winlist[i]))
			}
			nums = append(nums, strings.TrimSpace(bonus))

			InsertNums(env.MyDB, strconv.Itoa(round), strings.Join(nums, " "))

			time.Sleep(time.Millisecond * 500)
		}

	}()
	r := &Res{
		Status:   http.StatusOK,
		Response: "OK",
	}
	return c.JSONPretty(http.StatusOK, r, " ")
}

func GetLottoNum(c echo.Context) error {
	var nums string
	var err error
	round := c.Param("round")

	if round == "" || round == "latest" {
		nums, err = getLatestNums(env.MyDB)
	} else {
		nums, err = getNums(env.MyDB, round)
	}
	if err != nil {
		return c.JSON(http.StatusMethodNotAllowed, "check fail")
	}
	r := &Res{
		Status:   http.StatusOK,
		Response: nums,
	}
	return c.JSONPretty(http.StatusOK, r, " ")
}
