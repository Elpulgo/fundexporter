package models

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

// Config blaha
type Config struct {
	FundIds []int    `json:"fundIds"`
	Filter  []string `json:"filter"`
}

// FundList blaha
type FundList struct {
	List []FundInfo `json:"fundListViews"`
}

// FundInfo ...
type FundInfo struct {
	Name           string  `json:"name"`
	DevOneMonth    float64 `json:"developmentOneMonth"`
	DevThreeMonths float64 `json:"developmentThreeMonths"`
	DevSixMonths   float64 `json:"developmentSixMonths"`
	DevOneyear     float64 `json:"developmentOneYear"`
	ManagementFee  float64 `json:"managementFee"`
	TransactionFee float64 `json:"transactionFee"`
	TotalFee       float64 `json:"totalFee"`
	ID             int     `json:"orderBookId"`
	Nav            float64
	MA30           float64
	MA50           float64
	MA200          float64
	RSI            float64
	Amom1_3        float64
	Amom1_3_6      float64
	Amom3_6_12     float64
}

// SliceList blaha..
func (fundList *FundList) SliceList() []FundInfo {
	if len(fundList.List) > 10 {
		return fundList.List[0:10]
	}

	return fundList.List
}

// IsNAVAboveMA blaha
func (fundInfo *FundInfo) IsNAVAboveMA() string {
	if fundInfo.Nav > fundInfo.MA50 && fundInfo.Nav > fundInfo.MA200 && fundInfo.Nav > fundInfo.MA30 {
		return "Ja"
	}
	return "Nej"
}

// GetAMOM1_3 blaha..
func (amomValues *AMOMValues) GetAMOM1_3() float64 {
	value := (amomValues.ChangeSinceOneMonth + amomValues.ChangeSinceThreeMonths) / 2
	return math.Round(value*10) / 10
}

// GetAMOM1_3_6 blaha..
func (amomValues *AMOMValues) GetAMOM1_3_6() float64 {
	value := (amomValues.ChangeSinceOneMonth + amomValues.ChangeSinceThreeMonths + amomValues.ChangeSinceSixMonths) / 3
	return math.Round(value*10) / 10
}

// GetAMOM3_6_12 blaha..
func (amomValues *AMOMValues) GetAMOM3_6_12() float64 {
	value := (amomValues.ChangeSinceThreeMonths + amomValues.ChangeSinceSixMonths + amomValues.ChangeSinceOneYear) / 3
	return math.Round(value*10) / 10
}

// AMOMValues blaha..
type AMOMValues struct {
	ChangeSinceOneMonth    float64
	ChangeSinceThreeMonths float64
	ChangeSinceSixMonths   float64
	ChangeSinceOneYear     float64
}

// MovingAverageInfo blaha..
type MovingAverageInfo struct {
	TechnicalAnalysis []TechAnalysis `json:"technicalAnalysis"`
	High              float64        `json:"high"`
	Low               float64        `json:"low"`
	NAV               float64        `json:"lastPrice"`
	ChangePercent     float64        `json:"changePercent"`
}

// TechAnalysis blaha
type TechAnalysis struct {
	DataPoints []Point `json:"dataPoints"`
	TimeFrame  int     `json:"timeFrame"`
	Type       string  `json:"type"`
}

// Point blaha..
type Point struct {
	// Timestamp float64
	Value float64
}

// Series blaha..
type Series struct {
	DataSeries []DataSerie `json:"dataSerie"`
}

// DataSerie blaha
type DataSerie struct {
	Timestamp  int     `json:"x"`
	Percentage float64 `json:"y"`
}

// UnmarshalJSON blaha..
func (point *Point) UnmarshalJSON(data []byte) error {
	var v []interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		fmt.Printf("Error while decoding %v\n", err)
		return err
	}

	point.Value = v[1].(float64)

	return nil
}

// UERateResult blaha..
type UERateResult struct {
	Status  string   `json:"status"`
	Results UESeries `json:"Results"`
}

// UESeries blaha
type UESeries struct {
	Series []UEDataContainer `json:"series"`
}

// UEDataContainer blaha
type UEDataContainer struct {
	Data []UEData `json:"data"`
}

// UEData blaha
type UEData struct {
	Value string `json:"value"`
}

func parseAndRoundFloatValue(value string) float64 {
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		fmt.Print(err.Error())
	}
	return math.Round(parsed*10) / 10
}
