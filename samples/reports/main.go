package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"

	"github.com/neosus/itc-go"
)

var (
	keyID          = flag.String("key_id", "", "KeyID")
	issuerID       = flag.String("issuer_id", "", "IssueID")
	vendorNumber   = flag.String("vendor_number", "", "Vendor Number")
	privateKeyFile = flag.String("private_key", "", "Path to private key associated with key_id")
	salesOut       = flag.String("sales_out", "sales.gzip", "file to save sales and trends report")
	financeOut     = flag.String("finance_out", "finance.gzip", "file to save finance report")
)

func main() {
	flag.Parse()
	if *keyID == "" || *issuerID == "" || *vendorNumber == "" || *privateKeyFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	keyData, err := ioutil.ReadFile(*privateKeyFile)
	if err != nil {
		log.Fatal(err)
	}

	block, _ := pem.Decode([]byte(keyData))
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	privateKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		log.Fatal("failed to parse key")
	}

	client := itc.NewClient(*keyID, *issuerID, privateKey)

	log.Println("request sales and trends reports")
	data := url.Values{}
	data.Set(itc.SalesFrequencyFilter, itc.DailyFrequency)
	data.Set(itc.SalesReportSubTypeFilter, itc.SummaryReportSubType)
	data.Set(itc.SalesReportTypeFilter, itc.SalesReportType)
	data.Set(itc.VendorNumberFilter, *vendorNumber)

	resp, err := client.GetSalesReport(context.Background(), data)
	if err != nil {
		log.Fatal(err)
	}

	if err := saveToFile(resp, *salesOut); err != nil {
		log.Fatalf("failed to save sales and trends reports: %s", err)
	}
	log.Println("sales and trends reports saved to " + *salesOut)

	log.Println("request finance reports")
	data = url.Values{}
	data.Set(itc.FinanceReportRegionCodeFilter, "US")
	data.Set(itc.FinanceReportDateFilter, "2018-11")
	data.Set(itc.VendorNumberFilter, *vendorNumber)

	resp, err = client.GetFinanceReport(context.Background(), data)
	if err != nil {
		log.Fatal(err)
	}

	if err := saveToFile(resp, *financeOut); err != nil {
		log.Fatalf("failed to save finance reports: %s", err)
	}
	log.Println("finance reports saved to " + *financeOut)
}

func saveToFile(in io.Reader, out string) error {
	salesFile, err := os.Create(out)
	if err != nil {
		return err
	}
	defer salesFile.Close()

	_, err = io.Copy(salesFile, in)
	return err
}
