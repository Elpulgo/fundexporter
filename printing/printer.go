package printing

import (
	models "avanza/models"
	"fmt"
	"math"
	"text/tabwriter"
)

// PrintFundInfos blaha
func PrintFundInfos(fundList []models.FundInfo, category string, tabWriter *tabwriter.Writer) {

	fmt.Fprintf(tabWriter, "\n\n%s\n", category)

	fmt.Fprintf(tabWriter, "\n%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
		"Id",
		"Name",
		"1M (%)",
		"3M (%)",
		"6M (%)",
		"12M (%)",
		"NAV Above MA",
		"NAV",
		"MA30",
		"MA50",
		"MA200",
		"RSI (14 days)",
		"Mgmt Fee",
		"Trans Fee",
		"Total Fee",
		"AMOM 1+3",
		"AMOM 1+3+6",
		"AMOM 3+6+12")

	for _, fund := range fundList {
		fmt.Fprintf(tabWriter, "\n%v\t%s\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v",
			fund.ID,
			fund.Name,
			math.Round(fund.DevOneMonth*10)/10,
			math.Round(fund.DevThreeMonths*10)/10,
			math.Round(fund.DevSixMonths*10)/10,
			math.Round(fund.DevOneyear*10)/10,
			fund.IsNAVAboveMA(),
			math.Round(fund.Nav*10)/10,
			math.Round(fund.MA30*10)/10,
			math.Round(fund.MA50*10)/10,
			math.Round(fund.MA200*10)/10,
			math.Round(fund.RSI*10)/10,
			math.Round(fund.ManagementFee*10)/10,
			math.Round(fund.TransactionFee*10)/10,
			math.Round(fund.TotalFee*10)/10,
			fund.Amom1_3,
			fund.Amom1_3_6,
			fund.Amom3_6_12)
	}
}
