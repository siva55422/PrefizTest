package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
)

type ShoppingBasketReceipt map[string]float64

const (
	TOTAL                   = "Total"
	MUSICCD                 = "music CD"
	SALESTAX                = "salesTax"
	IMPORTED                = "import"
	NONEXEMPTION            = "nonExemption"
	BOTTLEOFPERFUME         = "bottle of perfume"
	IMPORTEDBOXOFCHOCOLATES = "imported box of chocolates"
	IMPORTEDBOTTLEOFPERFUME = "imported bottle of perfume"
)

var taxNonExem []string = []string{MUSICCD, BOTTLEOFPERFUME}
var taxNonExemImported []string = []string{IMPORTEDBOXOFCHOCOLATES, IMPORTEDBOTTLEOFPERFUME}

func handleTestReq(w http.ResponseWriter, r *http.Request) {

	var inputReceipt ShoppingBasketReceipt
	json.NewDecoder(r.Body).Decode(&inputReceipt)

	totalTax := 0.00
	totalAmount := 0.00
	for k, v := range inputReceipt {
		n := v
		for _, importItem := range taxNonExemImported {
			if k == importItem {
				totalTax, v = ProcessTaxNonExemption(IMPORTED, n, v, k, inputReceipt, totalTax)
				for _, item := range taxNonExem {
					if strings.Contains(k, item) {
						totalTax, v = ProcessTaxNonExemption(NONEXEMPTION, n, v, k, inputReceipt, totalTax)
					}
				}
			}
		}

		for _, item := range taxNonExem {
			if k == item {
				totalTax, v = ProcessTaxNonExemption(NONEXEMPTION, n, v, k, inputReceipt, totalTax)
			}
		}

		totalAmount = totalAmount + v
	}

	inputReceipt[SALESTAX] = math.Round(totalTax*100) / 100
	inputReceipt[TOTAL] = math.Round(totalAmount*100) / 100

	json.NewEncoder(w).Encode(&inputReceipt)

}

func ProcessTaxNonExemption(taxType string, n float64, v float64, k string, inputReceipt map[string]float64, totalTax float64) (float64, float64) {
	var imp1 float64
	if taxType == IMPORTED {
		imp1 = (n * 5) / 100
	} else if taxType == NONEXEMPTION {
		imp1 = (n * 10) / 100
	}
	impRoundValue := math.Round(imp1*100) / 100
	v = v + impRoundValue
	inputReceipt[k] = math.Round(v*100) / 100
	totalTax = totalTax + impRoundValue

	return totalTax, v
}

func main() {
	fmt.Println("Server Started with Port as 8080...")
	http.HandleFunc("/salestax", handleTestReq)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
