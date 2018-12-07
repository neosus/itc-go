# iTunes Connect Go library

A client for accessing the iTunesConnect API

Look how-to-use [sample](samples/reports/main.go)

Or look small example:
```go
	client := itc.NewClient("key-id", "issuer-id", privateKey)

	data := url.Values{}
	data.Set(itc.SalesFrequencyFilter, itc.DailyFrequency)
	data.Set(itc.SalesReportSubTypeFilter, itc.SummaryReportSubType)
	data.Set(itc.SalesReportTypeFilter, itc.SalesReportType)
	data.Set(itc.VendorNumberFilter, "my-vendor-number")

	resp, err := client.GetSalesReport(context.Background(), data)
	if err != nil {
		log.Fatal(err)
	}

	if err := saveToFile(resp, *salesOut); err != nil {
		log.Fatalf(err)
	}
```

## Supported API's
- Reports
  - Sales and trends
  - Finance