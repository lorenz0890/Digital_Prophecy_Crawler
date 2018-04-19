package bitInfoCharts

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/zoidberg77/DigitalProphecy-WebCrawler/utility"
	"strings"
	"fmt"
	"strconv"
)

func (bw BigWallet) String() string {
	return fmt.Sprintf("%v %v %v", bw.Address, bw.AmountBTC, bw.ChangeBTC)
}

//scrape biggest wallets
func GetBiggestWallets() []BigWallet {
	doc, err := goquery.NewDocument(TOP100_BIGGEST_WALLETS)
	utility.CheckErr(err)
	var bw []BigWallet
	if err == nil {

		level := 1
		doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
			level += 1

			var flagAmount = false
			var flagAddress = false
			var address string
			var amount float64
			var change float64

			s.Find("td").Each(func(i int, s *goquery.Selection) {

				//finde amount and store in var
				s.Find("a").Each(func(i int, s *goquery.Selection) {
					if level > 3 && !strings.Contains(s.Text(), "wallet") && !flagAddress {
						address = s.Text()
						flagAddress = true
					}
				})

				//find amount and store in var
				if level > 3 && strings.Contains(s.Text(),"BTC") && strings.Contains(s.Text(),"$") && !flagAmount {
					amount, err = strconv.ParseFloat(strings.Split(strings.Replace(strings.Replace(s.Text(), " ", "", -1), ",", ".", -1), "BTC")[0], 64)
					utility.CheckErr(err)
					flagAmount = true
				}

				// if there is change, store in var. else set var to 0
				s.Find(".hidden-phone").Each(func(i int, s *goquery.Selection) {
					if level > 3 && strings.Contains(s.Text(), "BTC") {
						change, err = strconv.ParseFloat(strings.Split(strings.Replace(strings.Replace(s.Text(), " ", "", -1), ",", ".", -1), "BTC")[0], 64)
						utility.CheckErr(err)
					} else { change = 0}
				})

				if flagAddress && flagAmount {
					bw =  append(bw, BigWallet{address, amount, change})
					flagAddress = false
					flagAmount = false
					address = ""
					amount = 0
				}
			})
		})
	/*
		for _, element := range bw {
			fmt.Println(element)
	}
*/	}
	return bw
}
