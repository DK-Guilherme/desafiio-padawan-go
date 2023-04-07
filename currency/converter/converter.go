package converter

import (
	"fmt"
	"strconv"
)

type ConvertedCurrency struct {
	ValorConvertido float64 `json:"valorConvertido"`
	SimboloDaMoeda  string  `json:"simboloDaMoeda"`
}

func Convert(amount float64, from string, to string, rate float64) ConvertedCurrency {
	conv := fmt.Sprintf("%.2f", amount*rate)
	var symbol string

	switch to {
	case "USD":
		symbol = "$"
	case "BRL":
		symbol = "R$"
	case "EUR":
		symbol = "€"
	case "BTC":
		symbol = "₿"
	default:
		symbol = "$"
	}

	res, _ := strconv.ParseFloat(conv, 64)
	var converted = ConvertedCurrency{
		ValorConvertido: res,
		SimboloDaMoeda:  symbol,
	}
	fmt.Println()

	return converted
}
