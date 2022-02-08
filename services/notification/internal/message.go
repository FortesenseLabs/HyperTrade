package internal

import (
	"fmt"
	"strings"
	"time"
)

func (t Telegram) FormatConfigsMessage(r ConfigsResponse) string {
	header := "*Configs*"

	var configs = []string{header}

	for i, config := range r.Configs {
		index := i + 1
		b := fmt.Sprintf(
			"`\n#%v\nSymbol: %v\n"+
				"Minimum: %v\n"+
				"Allowed: %v\n"+
				"Enabled: %v`",
			index, config.Symbol,
			config.Minimum,
			config.AllowedAmount,
			config.TradingEnabled,
		)
		configs = append(configs, b)
	}

	return strings.Join(configs, "\n")
}

func (t Telegram) FormatOrderMessage(p OrderEventPayload) string {
	message := fmt.Sprintf(
		"*Created %v Order*\n\n"+
			"`ID: %v\n"+
			"Type: %v\n"+
			"Symbol: %v\n"+
			"Price: %v\n"+
			"Quantity: %v`",
		p.Side,
		p.ID,
		p.Type,
		p.Symbol,
		p.Price,
		p.Quantity,
	)

	return message
}

func (t Telegram) FormatTradeMessage(p TradeEventPayload) string {
	time := p.Time.Format(time.RFC822)

	message := fmt.Sprintf(
		"*Executed Trade*\n\n"+
			"`ID: %v\n"+
			"Symbol: %v\n"+
			"Entry: %v\n"+
			"Exit: %v\n"+
			"Quantity: %v\n"+
			"Time: %v`",
		p.ID,
		p.Symbol,
		p.Entry,
		p.Exit,
		p.Quantity,
		time,
	)

	return message
}

func (t Telegram) FormatBalanceMessage(r BalanceResponse) string {
	header := "*Balance*\n"

	if r.Test {
		header = fmt.Sprintln("*Test*", header)
	}

	var balances = []string{header}
	var separator rune = '•'

	for _, balance := range r.Balance {
		b := fmt.Sprintf("`%c %v %v`", separator, balance.Asset, balance.Amount)
		balances = append(balances, b)
	}

	return strings.Join(balances, "\n")
}

func (t Telegram) FormatStatsMessage(r StatsResponse) string {
	var message string

	if r.Stats == nil {
		message = "*Stats*\n\n`No data available yet`"
	} else {
		message = fmt.Sprintf("*Stats*\n\n`Profit: %.4f\nLoss: %.4f`", r.Stats.Profit, r.Stats.Loss)
	}

	return message
}

func (t Telegram) FormatDumpMessage(symbol string, r DumpResponse) string {
	message := fmt.Sprintf("*Dump*\n\n`ID: %v\nSymbol: %v\nQuantity: %v`", r.ID, symbol, r.Quantity)

	return message
}

func (t Telegram) FormatErrorMessage(p CriticalErrorEventPayload) string {
	message := fmt.Sprintf("*Critical Error*\n\n`%v`", p.Error)

	return message
}

func (t Telegram) FormatUpdateTradingMessage(symbol string, enable bool) string {
	var message string

	var payload interface{}
	req := UpdateTradingRequest{symbol, enable}
	err := t.pubsub.Request(UpdateTradingEvent, req, &payload)

	if err != nil {
		message = err.Error()
	} else {
		var status string

		switch enable {
		case true:
			status = "enabled"
		case false:
			status = "disabled"
		}
		message = fmt.Sprintf("*Message*\n\n`Trading has been %v`", status)
	}

	return message
}
