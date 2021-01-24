package reports

import (
	"avanza/models"
	requests "avanza/requests"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

// CreateReportFile blaha..
func CreateReportFile(fileName string) *os.File {
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

// Generate blaha..
func Generate(fundList []models.FundInfo, header string, f *os.File) {
	fileText := `
` + header +
		`;

`

	for _, fund := range fundList {
		fileText += `;` + strconv.Itoa(fund.ID) + `;` + fund.Name + `;` + roundFloatValue(fund.DevOneMonth) + `;` + roundFloatValue(fund.DevThreeMonths) + `;` + roundFloatValue(fund.DevSixMonths) + `;` + roundFloatValue(fund.DevOneyear) + `;` + roundFloatValue(fund.Amom1_3) + `;` + roundFloatValue(fund.Amom1_3_6) + `;` + roundFloatValue(fund.Amom3_6_12) + `;` + fund.IsNAVAboveMA() + `;` + roundFloatValue(fund.Nav) + `;` + roundFloatValue(fund.MA30) + `;` + roundFloatValue(fund.MA50) + `;` + roundFloatValue(fund.MA200) + `;` + roundFloatValue(fund.RSI) + `;` + roundFloatValue(fund.ManagementFee) + `;` + roundFloatValue(fund.TransactionFee) + `;` + roundFloatValue(fund.TotalFee) + `;
`
	}

	f.WriteString(fileText)
}

// GenerateCurrent blaha..
func GenerateCurrent(fundList []models.FundInfo, f *os.File) {
	print("Printing current..")
	omxspi := requests.GetOMXSPIPercentage("three_months")
	MA12Unrate, currentUnrate, reccesionRisk := requests.GetUERate()
	dateNow := time.Now()
	dateFormatted := dateNow.Format("2006-01-02")

	fileText := ``

	for _, fund := range fundList {
		fileText += `;` + strconv.Itoa(fund.ID) + `;` + fund.Name + `;;;;;;` + roundFloatValue(fund.Nav) + `;` + roundFloatValue(omxspi) + `;` + roundFloatValue(fund.DevOneMonth) + `;` + roundFloatValue(fund.DevThreeMonths) + `;` + roundFloatValue(fund.DevSixMonths) + `;` + roundFloatValue(fund.DevOneyear) + `;` + roundFloatValue(fund.MA30) + `;` + roundFloatValue(fund.MA50) + `;` + roundFloatValue(fund.MA200) + `;` + roundFloatValue(fund.RSI) + `;
`
	}

	reccesionRiskText := "NEJ"
	if reccesionRisk == true {
		reccesionRiskText = "JA"
	}

	fileText += `
MA12 UNRATE:;` + roundFloatValue(MA12Unrate) + `
CURR UNRATE:;` + roundFloatValue(currentUnrate) + `
RECESSION RISK:;` + reccesionRiskText + `

Datum:;` + dateFormatted + `
`

	f.WriteString(fileText)
}

func roundFloatValue(value float64) string {
	return strconv.FormatFloat(math.Round(value*10)/10, 'f', 1, 64)
}
