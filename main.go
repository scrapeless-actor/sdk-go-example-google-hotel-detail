package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/scrapeless-ai/scrapeless-actor-sdk-go/scrapeless"
	proxyModel "github.com/scrapeless-ai/scrapeless-actor-sdk-go/scrapeless/proxy"
)

var (
	client *http.Client
)

type RequestParam struct {
	Q                string `json:"q"`
	Gl               string `json:"gl"`
	Hl               string `json:"hl"`
	Currency         string `json:"currency"`
	CheckInDate      string `json:"check_in_date"`
	CheckOutDate     string `json:"check_out_date"`
	Adults           int    `json:"adults"`
	Children         int    `json:"children"`
	ChildrenAges     string `json:"children_ages"`
	SortBy           int    `json:"sort_by"`
	MinPrice         int    `json:"min_price"`
	MaxPrice         int    `json:"max_Price"`
	PropertyTypes    string `json:"property_types"`
	Amenities        string `json:"amenities"`
	Rating           string `json:"rating"`
	Brands           string `json:"brands"`
	HotelClass       string `json:"hotel_class"`
	FreeCancellation int    `json:"free_cancellation"` // 传入1生效
	SpecialOffers    int    `json:"special_offers"`    // 传入1生效
	EcoCertified     int    `json:"eco_certified"`     // 传入1生效
	VacationRentals  int    `json:"vacation_rentals"`
	Bedrooms         int    `json:"bedrooms"`
	Bathrooms        int    `json:"bathrooms"`
	NextPageToken    string `json:"next_page_token"`
	PropertyToken    string `json:"property_token"`
}

func main() {
	// new actor
	actor := scrapeless.New(scrapeless.WithProxy(), scrapeless.WithStorage())
	defer actor.Close()
	var param = &RequestParam{}
	if err := actor.Input(param); err != nil {
		log.Fatal(err)
	}
	// get proxy url
	proxy, err := actor.Proxy.Proxy(context.TODO(), proxyModel.ProxyActor{
		Country:         "us",
		SessionDuration: 10,
	})
	fmt.Println(proxy)
	if err != nil {
		panic(err)
	}
	parse, err := url.Parse(proxy)
	if err != nil {
		panic(err)
	}
	// init client with proxy
	client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(parse)}}
	data, err := detail(context.Background(), param)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(data)
}
