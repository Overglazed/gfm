package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	//_ "github.com/browserutils/kooky/browser/all"
)

const pageSize int = 50

func main() {

	gfmSlugPtr := flag.String("campaign", "", "The name of the gfm campaign to target")
	cookieStringPtr := flag.String("cookie-string", "", "The cookie string to add to the cookie request header for authentication.")
	//outputPtr := flag.String("output", "csv", "Specifies the output format. Default: csv. Valid values: [csv, text]")
	//cookieLookupPtr := flag.Bool("auto-cookie-lookup", false, "Attempt to automatically grab auth cookie from open browser session.")

	flag.Parse()

	if *gfmSlugPtr == "" {
		log.Fatalln("no gfm fundraiser campaign specified. Use --campaign to specify")
	}

	// if *cookieLookupPtr {
	// 	kookies := kooky.ReadCookies(kooky.Valid, kooky.DomainHasSuffix(`gofundme.com`))

	// 	for _, k := range kookies {
	// 		fmt.Println(k.String())
	// 	}
	// 	os.Exit(0)
	// }

	donation_url := fmt.Sprintf("https://api.gofundme.com/co/v1/feeds/%s/donations", *gfmSlugPtr)
	cookieString := *cookieStringPtr

	donations := getDonations(donation_url, cookieString)

	donations = removeDuplicateIds(donations)

	transformed := transformTo2D(donations)

	w := csv.NewWriter(os.Stdout)
	w.WriteAll(transformed)
	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv: ", err)
	}
}

func getDonations(url string, cookie string) []Donation {
	c := newClient()
	var donations []Donation
	cToken := ""
	res, err := getPaginatedDonation(c, url, cookie, pageSize, cToken)
	if err != nil {
		fmt.Printf("error %s", err)
		os.Exit(1)
	}

	donations = append(donations, res.References.Donations...)

	for res.ViewModels[0].Next.HasNext {
		cToken = res.ViewModels[0].Next.Params.CToken
		res, err = getPaginatedDonation(c, url, cookie, pageSize, cToken)
		if err != nil {
			fmt.Printf("error %s", err)
			os.Exit(1)
		}

		donations = append(donations, res.References.Donations...)
	}

	return donations
}

func getPaginatedDonation(c *http.Client, api_url string, cookie string, pageSize int, cToken string) (*ApiResponse, error) {
	req, err := http.NewRequest("GET", api_url, nil)
	if err != nil {
		fmt.Printf("error %s", err)
		return nil, err
	}

	if cookie == "" {
		cookie = os.Getenv("GFM_AUTH_COOKIE")
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cookie", cookie)

	if len(cToken) > 0 {
		params := req.URL.Query()
		params.Set("ctoken", cToken)
		req.URL.RawQuery = params.Encode()
	}

	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("error %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	r := ApiResponse{
		References: ApiReference{},
		ViewModels: []ApiViewModel{},
	}

	json.NewDecoder(resp.Body).Decode(&r)

	return &r, nil
}

func newClient() *http.Client {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	return c
}

func removeDuplicateIds(d []Donation) []Donation {
	records := make(map[int]bool)
	result := []Donation{}

	for _, donation := range d {
		if _, ok := records[donation.DonationId]; !ok {
			records[donation.DonationId] = true
			result = append(result, donation)
		}
	}
	return result
}

func transformTo2D(d []Donation) [][]string {
	rows := len(d)
	res := make([][]string, rows+1)

	//add headers
	res[0] = []string{
		"donation_id",
		"amount",
		"currencycode",
		"name",
		"first_name",
		"last_name",
		"anonymous",
		"comment",
		"country",
		"timestamp",
	}

	//data rows
	for i := 0; i < rows; i++ {
		res[i+1] = []string{
			strconv.Itoa(d[i].DonationId),
			strconv.Itoa(d[i].Amount),
			d[i].CurrencyCode,
			d[i].Name,
			d[i].FirstName,
			d[i].LastName,
			strconv.FormatBool(d[i].Anonymous),
			d[i].Comment,
			d[i].Country,
			d[i].Timestamp,
		}
	}
	return res
}
