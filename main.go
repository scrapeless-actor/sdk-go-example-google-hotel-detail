package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

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
	FreeCancellation int    `json:"free_cancellation"`
	SpecialOffers    int    `json:"special_offers"`
	EcoCertified     int    `json:"eco_certified"`
	VacationRentals  int    `json:"vacation_rentals"`
	Bedrooms         int    `json:"bedrooms"`
	Bathrooms        int    `json:"bathrooms"`
	NextPageToken    string `json:"next_page_token"`
	PropertyToken    string `json:"property_token"`
}
type RequestParamGateway struct {
	Q             string `json:"q"`
	Gl            string `json:"gl"`
	Hl            string `json:"hl"`
	Currency      string `json:"currency"`
	CheckInDate   string `json:"check_in_date"`
	CheckOutDate  string `json:"check_out_date"`
	ChildrenAges  string `json:"children_ages"`
	PropertyTypes string `json:"property_types"`
	Amenities     string `json:"amenities"`
	Rating        string `json:"rating"`
	Brands        string `json:"brands"`
	HotelClass    string `json:"hotel_class"`
	NextPageToken string `json:"next_page_token"`
	PropertyToken string `json:"property_token"`
	Engine        string `json:"engine"`

	Adults           any `json:"adults"`
	Children         any `json:"children"`
	FreeCancellation any `json:"free_cancellation"`
	SpecialOffers    any `json:"special_offers"`
	EcoCertified     any `json:"eco_certified"`
	VacationRentals  any `json:"vacation_rentals"`
	Bedrooms         any `json:"bedrooms"`
	Bathrooms        any `json:"bathrooms"`
	SortBy           any `json:"sort_by"`
	MinPrice         any `json:"min_price"`
	MaxPrice         any `json:"max_Price"`
}

func (gateway RequestParamGateway) RequestParamGateway2RequestParam() *RequestParam {
	adults, _ := strconv.ParseInt(fmt.Sprintf("%v", gateway.Adults), 10, 64)
	sortBy, _ := strconv.ParseInt(fmt.Sprintf("%v", gateway.SortBy), 10, 64)
	children, _ := strconv.ParseInt(fmt.Sprintf("%v", gateway.Children), 10, 64)
	minPrice, _ := strconv.ParseInt(fmt.Sprintf("%v", gateway.MinPrice), 10, 64)
	maxPrice, _ := strconv.ParseInt(fmt.Sprintf("%v", gateway.MaxPrice), 10, 64)
	freeCancellation, _ := strconv.ParseInt(fmt.Sprintf("%v", gateway.FreeCancellation), 10, 64)
	specialOffers, _ := strconv.ParseInt(fmt.Sprintf("%v", gateway.SpecialOffers), 10, 64)
	ecoCertified, _ := strconv.ParseInt(fmt.Sprintf("%v", gateway.EcoCertified), 10, 64)
	vacationRentals, _ := strconv.ParseInt(fmt.Sprintf("%v", gateway.VacationRentals), 10, 64)
	bedrooms, _ := strconv.ParseInt(fmt.Sprintf("%v", gateway.Bedrooms), 10, 64)
	bathrooms, _ := strconv.ParseInt(fmt.Sprintf("%v", gateway.Bathrooms), 10, 64)
	return &RequestParam{
		Q:                gateway.Q,
		Gl:               gateway.Gl,
		Hl:               gateway.Hl,
		Currency:         gateway.Currency,
		CheckInDate:      gateway.CheckInDate,
		CheckOutDate:     gateway.CheckOutDate,
		Adults:           int(adults),
		Children:         int(children),
		ChildrenAges:     gateway.ChildrenAges,
		SortBy:           int(sortBy),
		MinPrice:         int(minPrice),
		MaxPrice:         int(maxPrice),
		PropertyTypes:    gateway.PropertyTypes,
		Amenities:        gateway.Amenities,
		Rating:           gateway.Rating,
		Brands:           gateway.Brands,
		HotelClass:       gateway.HotelClass,
		FreeCancellation: int(freeCancellation),
		SpecialOffers:    int(specialOffers),
		EcoCertified:     int(ecoCertified),
		VacationRentals:  int(vacationRentals),
		Bedrooms:         int(bedrooms),
		Bathrooms:        int(bathrooms),
		NextPageToken:    gateway.NextPageToken,
		PropertyToken:    gateway.PropertyToken,
		//Engine:           gateway.Engine,
	}
}
func main() {
	// new actor
	actor := scrapeless.New(scrapeless.WithProxy(), scrapeless.WithStorage())
	defer actor.Close()
	var param = &RequestParamGateway{}
	if err := actor.Input(param); err != nil {
		log.Fatal(err)
	}
	requestParam := param.RequestParamGateway2RequestParam()
	// get proxy url
	proxy, err := actor.Proxy.Proxy(context.TODO(), proxyModel.ProxyActor{
		Country:         "us",
		SessionDuration: 10,
	})
	if err != nil {
		panic(err)
	}
	//proxy = "http://group_scraper_google_trneds:c8d2279d492a@pm-gw-us.scrapeless.io:24125"
	parse, err := url.Parse(proxy)
	if err != nil {
		panic(err)
	}
	fmt.Println(parse.String())
	// init client with proxy
	client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(parse)}}
	data, err := detail(context.Background(), requestParam)
	if err != nil {
		panic(err)
	}
	log.Println(data)
	ok, err := actor.Storage.GetKv().SetValue(context.Background(), "hotel-detail", data, 0)
	if err != nil {
		panic(err)
	}
	log.Println(ok)
}
