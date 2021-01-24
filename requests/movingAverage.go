package requests

import (
	models "avanza/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"text/tabwriter"
)

// GetMovingAverage returns MA, RSI
func GetMovingAverage(fundID int, period string, timeFrameSma int) (float64, float64) {
	payload := `{
		"orderbookId":"` + strconv.Itoa(fundID) + `",
		"chartResolution": "DAY",
		"navigator": false,
		"percentage": false,
		"timePeriod": "` + period + `",
		"chartType": "AREA",
		"owners": false,
		"volume": false,
		"ta": [
			{
				"type": "sma",
				"timeFrame": ` + strconv.Itoa(timeFrameSma) + `
			},		
			{
				"type": "rsi",
				"timeFrame": 14
			}	
		]
	  }`

	resp, err := http.Post(
		"https://www.avanza.se/ab/component/highstockchart/getchart/orderbook",
		"application/json",
		bytes.NewBuffer([]byte(payload)),
	)
	if err != nil {
		print(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}

	var result models.MovingAverageInfo

	err = json.Unmarshal(body, &result)

	if err != nil {
		print(err.Error())
	}

	var maPoint, rsi float64

	maPointLength := len(result.TechnicalAnalysis[0].DataPoints)
	rsiLength := len(result.TechnicalAnalysis[1].DataPoints)

	if maPointLength > 0 {
		maPoint = result.TechnicalAnalysis[0].DataPoints[len(result.TechnicalAnalysis[0].DataPoints)-1].Value
	}

	if rsiLength > 0 {
		rsi = result.TechnicalAnalysis[1].DataPoints[len(result.TechnicalAnalysis[1].DataPoints)-1].Value
	}

	return maPoint, rsi
}

// GetAMOMValues blaha
func GetAMOMValues(fundID int) models.AMOMValues {
	resp, err := http.Get("https://www.avanza.se/_mobile/market/fund/" + strconv.Itoa(fundID))
	if err != nil {
		fmt.Print(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}

	var amomValues models.AMOMValues

	err = json.Unmarshal(body, &amomValues)

	if err != nil {
		print(err.Error())
	}

	return amomValues
}

// GetSingleFundInfo blaha
func GetSingleFundInfo(fundID int) models.FundInfo {
	resp, err := http.Get("https://www.avanza.se/_api/fund-guide/guide/" + strconv.Itoa(fundID))
	if err != nil {
		fmt.Print(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}

	var fundInfo models.FundInfo

	err = json.Unmarshal(body, &fundInfo)

	if err != nil {
		print(err)
	}

	return fundInfo
}

// GetOMXSPIPercentage blaha
func GetOMXSPIPercentage(period string) float64 {
	payload := `{
		"orderbookId":18988,
		"chartType":"AREA",
		"chartResolution":"DAY",
		"percentage":true,
		"timePeriod": "` + period + `"
	}
	`
	resp, err := http.Post(
		"https://www.avanza.se/ab/component/highstockchart/getchart/orderbook",
		"application/json",
		bytes.NewBuffer([]byte(payload)),
	)
	if err != nil {
		print(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}

	var result models.MovingAverageInfo

	err = json.Unmarshal(body, &result)

	if err != nil {
		print(err)
	}

	return result.ChangePercent
}

// GetCategoryData blaha..
func GetCategoryData(payload string, header string, tabWriter *tabwriter.Writer) models.FundList {
	resp, err := http.Post(
		"https://www.avanza.se/_api/fund-guide/list",
		"application/json",
		bytes.NewBuffer([]byte(payload)),
	)
	if err != nil {
		print(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}

	var result models.FundList

	err = json.Unmarshal(body, &result)

	if err != nil {
		print(err)
	}

	return result
}

// GetUERate return SMA12 and current value, and risk for recession
func GetUERate() (float64, float64, bool) {
	resp, err := http.Get("https://api.bls.gov/publicAPI/v2/timeseries/data/LNS14000000")
	if err != nil {
		fmt.Print(err.Error())
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err.Error())
	}

	var ueRateResult models.UERateResult
	err = json.Unmarshal(body, &ueRateResult)

	if err != nil {
		print(err.Error())
	}

	if ueRateResult.Status == "REQUEST_NOT_PROCESSED" {
		fmt.Println("Threshold is reached for UNRATE")
		return 0, 0, false
	}

	var ueDataPoints = ueRateResult.Results.Series[0].Data
	var buffer float64
	var currentValue float64

	for index, uePoint := range ueDataPoints {
		if index > 11 {
			break
		}

		value, err := strconv.ParseFloat(uePoint.Value, 64)
		if err != nil {
			print(err.Error())
		}

		if index <= 11 {
			buffer += value
		}

		if index == 0 {
			currentValue = value
		}
	}

	var sma12 = buffer / 12
	sma12 = math.Round(sma12*10) / 10
	currentValue = math.Round(currentValue*10) / 10

	return sma12, currentValue, currentValue > sma12
}
