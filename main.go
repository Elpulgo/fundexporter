package main

import (
	models "avanza/models"
	printer "avanza/printing"
	reports "avanza/reports"
	requests "avanza/requests"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"text/tabwriter"
)

var (
	printCurrent  bool
	reportCurrent bool
	report        bool
)

func main() {
	log.Print("Starting Avanza fund exporter ..")

	if len(os.Args) > 1 {
		if os.Args[1] == "current" {
			printCurrent = true
		} else if os.Args[1] == "report" {
			report = true
		}

		if len(os.Args) > 2 {
			if os.Args[1] == "report" && os.Args[2] == "current" {
				reportCurrent = true
			}
		}
	}

	execute()
}

func execute() {
	tabWriter := new(tabwriter.Writer)
	tabWriter.Init(os.Stdout, 12, 8, 1, '\t', 0)

	defer tabWriter.Flush()

	configContent, err := ioutil.ReadFile("config.json")
	if err != nil {
		print(err)
	}

	var config models.Config

	err = json.Unmarshal(configContent, &config)
	if err != nil {
		print(err)
	}

	if printCurrent == true || reportCurrent == true {
		getCurrentFunds(config, tabWriter)
		return
	}

	if report == true {
		fmt.Println("Printing report ...")
	}

	file := reports.CreateReportFile("report.csv")

	for _, filter := range config.Filter {
		categoryContent := strings.Split(filter, ":")

		payload := getPayload(filter)
		categoryData := requests.GetCategoryData(payload, categoryContent[1], tabWriter)
		fundList := getFundList(categoryData.SliceList(), tabWriter)

		if report == true {
			reports.Generate(fundList, categoryContent[1], file)
		} else {
			printer.PrintFundInfos(fundList, categoryContent[1], tabWriter)
		}
	}
}

func getCurrentFunds(config models.Config, tabWriter *tabwriter.Writer) {
	var myFunds []models.FundInfo

	for _, fundID := range config.FundIds {
		fund := requests.GetSingleFundInfo(fundID)
		fund.ID = fundID

		ma30, rsi := requests.GetMovingAverage(fund.ID, "six_months", 30)
		ma50, _ := requests.GetMovingAverage(fund.ID, "six_months", 50)
		ma200, _ := requests.GetMovingAverage(fund.ID, "six_months", 200)

		amomValues := requests.GetAMOMValues(fund.ID)

		fund.MA30 = ma30
		fund.MA50 = ma50
		fund.MA200 = ma200
		fund.RSI = rsi
		fund.Amom1_3 = amomValues.GetAMOM1_3()
		fund.Amom1_3_6 = amomValues.GetAMOM1_3_6()
		fund.Amom3_6_12 = amomValues.GetAMOM3_6_12()

		myFunds = append(myFunds, fund)
	}

	sort.SliceStable(myFunds, func(i, j int) bool {
		return myFunds[i].DevThreeMonths > myFunds[j].DevThreeMonths
	})

	if reportCurrent == true {
		reportFile := reports.CreateReportFile("report-current.csv")
		reports.GenerateCurrent(myFunds, reportFile)
		println("Generating report for current funds ...")
		return
	}

	printer.PrintFundInfos(myFunds, "Current", tabWriter)
}

func getFundList(fundList []models.FundInfo, tabWriter *tabwriter.Writer) []models.FundInfo {
	size := len(fundList)
	newFundList := make([]models.FundInfo, 0, size)

	ch := make(chan models.FundInfo, size)
	wg := sync.WaitGroup{}

	for i := 0; i < size; i++ {
		wg.Add(1)
		go getFundInfo(fundList[i], ch, &wg)
	}

	wg.Wait()
	close(ch)

	for fund := range ch {
		newFundList = append(newFundList, fund)
	}

	sort.SliceStable(newFundList, func(i, j int) bool {
		return newFundList[i].DevThreeMonths > newFundList[j].DevThreeMonths
	})

	return newFundList
}

func getFundInfo(info models.FundInfo, ch chan models.FundInfo, wg *sync.WaitGroup) {

	ma30, rsi := requests.GetMovingAverage(info.ID, "six_months", 30)
	ma50, _ := requests.GetMovingAverage(info.ID, "six_months", 50)
	ma200, _ := requests.GetMovingAverage(info.ID, "six_months", 200)

	fundInfo := requests.GetSingleFundInfo(info.ID)
	amomValues := requests.GetAMOMValues(info.ID)

	info.Nav = fundInfo.Nav
	info.MA30 = ma30
	info.MA50 = ma50
	info.MA200 = ma200
	info.RSI = rsi
	info.DevSixMonths = fundInfo.DevSixMonths
	info.Amom1_3 = amomValues.GetAMOM1_3()
	info.Amom1_3_6 = amomValues.GetAMOM1_3_6()
	info.Amom3_6_12 = amomValues.GetAMOM3_6_12()

	ch <- info
	wg.Done()
}

func getPayload(filter string) string {

	filterParts := strings.Split(filter, ":")
	category := filterParts[0]
	value := filterParts[1]

	payload := `{			
		"startIndex": 0,
		"indexFund": $$INDEX$$,
		"sustainabilityProfile": false,
		"lowCo2": false,
		"noFossilFuelInvolvement": false,
		"regionFilter": [$$REGION$$],
		"countryFilter": [],
		"alignmentFilter": [],
		"industryFilter": [$$INDUSTRY$$],
		"fundTypeFilter": [$$FUNDTYPE$$],
		"interestTypeFilter": [],
		"sortField": "developmentThreeMonths",
		"sortDirection": "DESCENDING",
		"name": "",
		"recommendedHoldingPeriodFilter": [],
		"companyFilter": [],
		"productInvolvementsFilter": [],
		"ratingFilter": [],
		"sustainabilityRatingFilter": [],
		"environmentalRatingFilter": [],
		"socialRatingFilter": [],
		"governanceRatingFilter": []
	  }`

	if category == "industry" {
		payload = strings.Replace(payload, "$$INDUSTRY$$", `"`+value+`"`, -1)
	} else {
		payload = strings.Replace(payload, "$$INDUSTRY$$", ``, -1)
	}

	if category == "region" {
		payload = strings.Replace(payload, "$$REGION$$", `"`+value+`"`, -1)
	} else {
		payload = strings.Replace(payload, "$$REGION$$", ``, -1)
	}

	if category == "fundType" {
		payload = strings.Replace(payload, "$$FUNDTYPE$$", `"`+value+`"`, -1)
	} else {
		payload = strings.Replace(payload, "$$FUNDTYPE$$", ``, -1)
	}

	if category == "index" {
		payload = strings.Replace(payload, "$$INDEX$$", `true`, -1)
	} else {
		payload = strings.Replace(payload, "$$INDEX$$", `false`, -1)
	}

	return payload
}
