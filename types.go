package main

type ApiResponse struct {
	References ApiReference   `json:"references"`
	ViewModels []ApiViewModel `json:"view_models"`
}

type ApiViewModel struct {
	Type string          `json:"type"`
	Data ApiDonationData `json:"data"`
	Next ApiHasNext      `json:"next"`
}

type ApiDonationData struct {
	DonationIds    []int `json:"donation_ids"`
	TotalDonations int   `json:"total_donations"`
}

type ApiHasNext struct {
	HasNext bool      `json:"has_next"`
	Params  ApiParams `json:"params"`
}

type ApiParams struct {
	CToken string `json:"ctoken"`
}

type ApiReference struct {
	Donations []Donation `json:"donations"`
}
type Donation struct {
	DonationId   int    `json:"donation_id"`
	Amount       int    `json:"amount"`
	CurrencyCode string `json:"currencycode"`
	Name         string `json:"name"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Anonymous    bool   `json:"is_anonymous"`
	Comment      string `json:"comment"`
	Country      string `json:"country"`
	Timestamp    string `json:"timestamp"`
}
