package itc

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const ReportsVersion = "v1"

const VendorNumberFilter = "filter[vendorNumber]"
const VersionFilter = "filter[version]"

const (
	SalesFrequencyFilter = "filter[frequency]"
	DailyFrequency       = "DAILY"
	WeeklyFrequency      = "WEEKLY"
	MonthlyFrequency     = "MONTHLY"
	YearlyFrequency      = "YEARLY"
)

const (
	SalesReportSubTypeFilter = "filter[reportSubType]"
	SummaryReportSubType     = "SUMMARY"
	DetailedReportSubType    = "DETAILED"
	OptInReportSubType       = "OPT_IN"
)

const (
	SalesReportTypeFilter       = "filter[reportType]"
	SalesReportType             = "SALES"
	PreOrderReportType          = "PRE_ORDER"
	NewStandReportType          = "NEWSTAND"
	SubscriptionReportType      = "SUBSCRIPTION"
	SubscriptionEventReportType = "SUBSCRIPTION_EVENT"
	SubscriberReportType        = "SUBSCRIBER"
)

const SalesReportDateFilter = "filter[reportDate]"

func (c *client) GetSalesReport(ctx context.Context, data url.Values) (io.Reader, error) {
	headers := map[string]string{"Accept": "application/a-gzip"}
	return c.makeRequest(ctx, http.MethodGet, headers,
		strings.Join([]string{ReportsVersion, "salesReports"}, "/"), data)
}

const FinanceReportRegionCodeFilter = "filter[regionCode]"
const FinanceReportDateFilter = "filter[reportDate]"

func (c *client) GetFinanceReport(ctx context.Context, data url.Values) (io.Reader, error) {
	headers := map[string]string{"Accept": "application/a-gzip"}

	data.Set("filter[reportType]", "FINANCIAL")

	return c.makeRequest(ctx, http.MethodGet, headers,
		strings.Join([]string{ReportsVersion, "financeReports"}, "/"), data)
}
