package aggregate_test

import (
	"bytes"
	"encoding/csv"
	"github.com/kserik/horizon-task/internal/app/aggregate"
	"github.com/kserik/horizon-task/internal/pkg/exchange"
	"github.com/kserik/horizon-task/internal/pkg/repos"
	"github.com/kserik/horizon-task/internal/pkg/repos/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestRunSuccess(t *testing.T) {
	var correcttestData = `"app","ts","event","project_id","source","ident","user_id","session_id","country","device_type","device_os","device_os_ver","device_browser","device_browser_ver","props","nums"
"seq-market","2024-04-15 02:15:07.167","BUY_ITEMS","4974","","1","0896ae95dcaeee38e83fa5c43bef99780d7b2be23bcab36214","5d8afd8fec2fbf3e","DE","desktop","linux","x86_64","chrome","122.0.0.0","{""tokenId"":""215"",""txnHash"":""0xd919290e80df271e77d1cbca61f350d2727531e0334266671ec20d626b2104a2"",""chainId"":""137"",""collectionAddress"":""0x22d5f9b75c524fec1d6619787e582644cd4d7422"",""currencyAddress"":""0xd1f9c58e33933a993a3891f8acfe05a68e1afc05"",""currencySymbol"":""SFL"",""marketplaceType"":""amm"",""requestId"":""""}","{""currencyValueDecimal"":""0.6136203411678249"",""currencyValueRaw"":""613620341167824900""}"`

	exSvc := exchange.NewService()
	repo := &mocks.MarketplaceData{}
	repo.On("InsertTransactions", mock.Anything).Return(nil).Once()

	a := aggregate.NewApp(exSvc, repo)

	csvReader := csv.NewReader(bytes.NewBufferString(correcttestData))
	assert.NoError(t, a.Run(csvReader))

	arg := repo.Calls[0].Arguments.Get(0)
	trxs, ok := arg.([]repos.Transaction)
	assert.True(t, ok)

	ts, err := time.Parse("2006-01-02 15:04:05.999", "2024-04-15 02:15:07.167")
	assert.NoError(t, err)
	assert.Equal(t, trxs[0], repos.Transaction{
		Timestamp:         ts,
		Event:             "BUY_ITEMS",
		ProjectID:         4974,
		CurrencySymbol:    "SFL",
		CurrencyValue:     decimal.RequireFromString("0.6136203411678249"),
		USDConversionRate: decimal.RequireFromString("0.10629575068267864"),
		USDValue:          decimal.RequireFromString("0.065225234798595323598961714490136"),
	})
}

func TestRunNoTransactions(t *testing.T) {
	var invalidData = `"seq-market","2024-04-15 02:15:07.167","BUY_ITEMS","4974","","1","0896ae95dcaeee38e83fa5c43bef99780d7b2be23bcab36214","5d8afd8fec2fbf3e","DE","desktop","linux","x86_64","chrome","122.0.0.0","{""tokenId"":""215"",""txnHash"":""0xd919290e80df271e77d1cbca61f350d2727531e0334266671ec20d626b2104a2"",""chainId"":""137"",""collectionAddress"":""0x22d5f9b75c524fec1d6619787e582644cd4d7422"",""currencyAddress"":""0xd1f9c58e33933a993a3891f8acfe05a68e1afc05"",""currencySymbol"":""SFL"",""marketplaceType"":""amm"",""requestId"":""""}","{""currencyValueDecimal"":""0.6136203411678249"",""currencyValueRaw"":""613620341167824900""}"`

	exSvc := exchange.NewService()
	repo := &mocks.MarketplaceData{}
	repo.On("InsertTransactions", mock.Anything).Return(nil).Once()

	a := aggregate.NewApp(exSvc, repo)

	csvReader := csv.NewReader(bytes.NewBufferString(invalidData))
	assert.Error(t, a.Run(csvReader))
}

func TestRunInvalidDate(t *testing.T) {
	var invlidData = `"app","ts","event","project_id","source","ident","user_id","session_id","country","device_type","device_os","device_os_ver","device_browser","device_browser_ver","props","nums"
"seq-market","202-04-15 02:15:07.167","BUY_ITEMS","4974","","1","0896ae95dcaeee38e83fa5c43bef99780d7b2be23bcab36214","5d8afd8fec2fbf3e","DE","desktop","linux","x86_64","chrome","122.0.0.0","{""tokenId"":""215"",""txnHash"":""0xd919290e80df271e77d1cbca61f350d2727531e0334266671ec20d626b2104a2"",""chainId"":""137"",""collectionAddress"":""0x22d5f9b75c524fec1d6619787e582644cd4d7422"",""currencyAddress"":""0xd1f9c58e33933a993a3891f8acfe05a68e1afc05"",""currencySymbol"":""SFL"",""marketplaceType"":""amm"",""requestId"":""""}","{""currencyValueDecimal"":""0.6136203411678249"",""currencyValueRaw"":""613620341167824900""}"`

	exSvc := exchange.NewService()
	repo := &mocks.MarketplaceData{}
	repo.On("InsertTransactions", mock.Anything).Return(nil).Once()

	a := aggregate.NewApp(exSvc, repo)

	csvReader := csv.NewReader(bytes.NewBufferString(invlidData))
	assert.Error(t, a.Run(csvReader))
}
