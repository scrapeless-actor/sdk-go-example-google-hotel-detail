package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func detail(ctx context.Context, params *RequestParam) (string, error) {
	var (
		valueForm  = &url.Values{}
		valueQuery = &url.Values{}
	)
	valueQuery.Set("hl", params.Hl)
	valueQuery.Set("rpcids", "AtySUc")
	valueQuery.Set("source-path", "/travel/search")
	valueQuery.Set("gl", params.Gl)
	valueQuery.Set("soc-device", "1")
	valueQuery.Set("soc-app", "162")
	valueQuery.Set("soc-platform", "1")
	valueQuery.Set("rt", "c")
	parse, _ := url.Parse("https://www.google.com/_/TravelFrontendUi/data/batchexecute")

	parse.RawQuery = valueQuery.Encode()
	param, _ := someParams(params)
	valueForm.Set("f.req", string(param))
	request, _ := http.NewRequestWithContext(ctx, http.MethodPost, "https://www.google.com/_/TravelFrontendUi/data/batchexecute", strings.NewReader(valueForm.Encode()))
	request.Header.Set("accept", "*/*")
	request.Header.Set("Content-Length", "*/*")
	request.Header.Set("Host", "*/*")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("content-type", "application/x-www-form-urlencoded")
	request.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36")
	do, err := client.Do(request)
	if err != nil {
		return "", err
	}
	all, _ := io.ReadAll(do.Body)
	defer do.Body.Close()
	wrbs := GetWrbs(all)
	result := gjson.Parse(string(wrbs[0]))
	data := gjson.Parse(result.Get("0.2").Array()[0].String()) //.Get("1")
	res := data.Get("1").String()
	hotelDetail := getAll(res)
	if strings.Contains(strings.ToLower(string(wrbs[0])), "vacation rental") {
		hotelDetail = getVacationRental(res)
	}
	var resultBytes, _ = json.Marshal(hotelDetail)

	return string(resultBytes), nil
}

func someParams(params *RequestParam) ([]byte, error) {

	t, err := time.Parse("2006-01-02", params.CheckInDate)
	if err != nil {
		return nil, err
	}
	year, month, day := t.Date()
	var checkInDateArr = []int{year, int(month), day}

	t2, err := time.Parse("2006-01-02", params.CheckOutDate)
	if err != nil {
		return nil, err
	}
	year2, month2, day2 := t2.Date()
	var checkOutDateArr = []int{year2, int(month2), day2}

	diff := t2.Sub(t)
	daysDiff := int(diff.Hours() / 24)

	var checkDateArr = []any{checkInDateArr, checkOutDateArr, daysDiff}

	var peopleArr []any
	for i := 0; i < params.Adults; i++ {
		peopleArr = append(peopleArr, []int{3})
	}
	if params.Children > 0 {
		for i := 0; i < params.Children; i++ {
			ages := strings.Split(params.ChildrenAges, ",")
			age, _ := strconv.Atoi(ages[i])
			peopleArr = append(peopleArr, []int{2, age})
		}
	}
	var allPeopleArr = []any{peopleArr, 1}

	priceArr := make([]any, 3)
	priceArr = []any{[]any{nil, params.MinPrice}, []any{nil, params.MaxPrice}, 1}

	// property type
	var propertyTypeArr []any
	if params.PropertyTypes != "" {
		splitArr := strings.Split(params.PropertyTypes, ",")
		for i := range splitArr {
			propertyType, _ := strconv.Atoi(splitArr[i])
			propertyTypeArr = append(propertyTypeArr, propertyType)
		}
	}

	// amenities
	var amenitiesArr []any
	if params.Amenities != "" {
		splitArr := strings.Split(params.Amenities, ",")
		for i := range splitArr {
			amenities, _ := strconv.Atoi(splitArr[i])
			amenitiesArr = append(amenitiesArr, amenities)
		}
	}

	// rating
	var ratingArr []any
	if params.Rating != "" {
		splitArr := strings.Split(params.Rating, ",")
		for i := range splitArr {
			rating, _ := strconv.Atoi(splitArr[i])
			ratingArr = append(ratingArr, rating)
		}
	}

	// hotel_class
	var hotelClassArr []any
	if params.HotelClass != "" {
		splitArr := strings.Split(params.HotelClass, ",")
		for i := range splitArr {
			hotelClass, _ := strconv.Atoi(splitArr[i])
			hotelClassArr = append(hotelClassArr, hotelClass)
		}
	}

	var filterArr2Child = []any{amenitiesArr, hotelClassArr, nil, params.FreeCancellation, nil, nil, params.Currency, nil, nil, nil, propertyTypeArr}
	var checkInInfoArr5 = []any{filterArr2Child, nil, []any{nil, params.Bedrooms, params.Bathrooms}, priceArr, ratingArr, params.SpecialOffers}

	var checkInfoArr3Child1 = []any{nil, []any{[]any{nil, nil, nil, nil, nil, nil, nil}}, []any{}}
	var checkInfoArr3Child2 = []any{nil, checkDateArr, nil, nil, nil, []any{1}}
	var checkInInfoArr3 = []any{checkInfoArr3Child1, checkInfoArr3Child2}
	var checkInInfoArr = []any{params.VacationRentals, allPeopleArr, checkInInfoArr3, nil, checkInInfoArr5}

	var filterArr []any
	if params.PropertyToken != "" {
		filterArr = []any{1, nil, nil, 0, 0, params.PropertyToken, params.SortBy, nil, 0}
	} else {
		filterArr = []any{1, nil, nil, 0, 0, nil, params.SortBy, nil, 0}
	}

	var hotelArrChild = []any{params.Q, checkInInfoArr, filterArr, []any{}}
	hotelJsonChild, err := json.Marshal(hotelArrChild)
	if err != nil {
		return nil, err
	}

	var hotelArr = []any{[]any{[]any{"AtySUc", string(hotelJsonChild), nil, "1"}}}

	// 最终的参数
	jsonData, err := json.Marshal(hotelArr)
	if err != nil {
		return nil, err
	}
	return jsonData, err
}

func GetWrbs(respBytes []byte) [][]byte {
	lines := bytes.Split(respBytes, []byte("\n"))
	var wrbs [][]byte
	for _, line := range lines {
		if bytes.HasPrefix(line, []byte(`[["wrb.fr"`)) {
			wrbs = append(wrbs, line)
		}
	}
	return wrbs
}

func getVacationRental(q string) HotelDetail {
	var (
		hotelDetail HotelDetail
	)

	a := gjson.Parse(q).Get("7").Array()[0].Get("441552390") //gps_coordinates
	var essentialInfo []string
	essentialInfoArray := a.Get("6.5.9.0.1").Get("179305178").Get("10.3.1").Array() //essential_info
	for _, result := range essentialInfoArray {
		essentialInfo = append(essentialInfo, result.Get("0").String())
	}
	var amenities []string
	var excludedAmenities []string
	for _, result := range a.Get("10.1.1").Array() {
		if result.Get("1").Bool() {
			amenities = append(amenities, result.Get("0").String())
			continue
		}
		excludedAmenities = append(excludedAmenities, result.Get("0").String())
	}

	overallRating := a.Get("7.0.0").Float() // 0-->overall_rating 1-->reviews
	reviews := a.Get("7.0.1").Float()       // 0-->overall_rating 1-->reviews
	var images []DetailImage
	for _, v := range a.Get("45").Array() {
		thumbnail := fmt.Sprintf("http:%s", v.Get("0").String())
		originalImage := v.Get("6.1").String()
		images = append(images, DetailImage{
			Thumbnail:     thumbnail,
			OriginalImage: originalImage,
		})
	}
	var price []DetailPrice
	for _, result := range a.Get("6.2.21").Array() {
		source := result.Get("0.0").String()
		link := fmt.Sprintf("https://www.google.com%s", result.Get("0.2").String())
		logo := strings.Replace(result.Get("0.3.0").String(), "//", "", -1)

		ratePerNight := result.Get("12.4.0").String()
		ratePerNightStr := strings.Replace(ratePerNight, "$", "", -1)
		ratePerNightStr = strings.Replace(ratePerNightStr, ",", "", -1)
		ratePerNightInt, _ := strconv.Atoi(ratePerNightStr)

		totalRate := result.Get("12.5.0").String()
		totalRateStr := strings.Replace(totalRate, "$", "", -1)
		totalRateStr = strings.Replace(totalRateStr, ",", "", -1)
		totalRateInt, _ := strconv.Atoi(totalRateStr)

		numGuests := result.Get("12.12.0.0").Int()
		price = append(price, DetailPrice{
			Source:    source,
			Link:      link,
			Logo:      logo,
			NumGuests: int(numGuests),
			RatePerNight: PriceRate{
				Lowest:          ratePerNight,
				ExtractedLowest: ratePerNightInt,
			},
			TotalRate: PriceRate{
				Lowest:          totalRate,
				ExtractedLowest: totalRateInt,
			},
		})
	}
	var nearbyPlaces []DetailNearbyPlace
	for _, result := range a.Get("2.19.0.1.2").Array() {
		var nearbyPlace DetailNearbyPlace
		name := result.Get("0").String()
		for _, v := range result.Get("2").Array() {
			transportationsTypeInt := transportationsType(v.Get("0").Int()) // 0 -->taxi  2 -->Walking  3-->Public transport
			transportationsDuration := v.Get("1").String()                  //
			nearbyPlace.Transportations = append(nearbyPlace.Transportations, struct {
				Type     string `json:"type"`
				Duration string `json:"duration"`
			}{
				Type:     transportationsTypeMapping[transportationsTypeInt],
				Duration: transportationsDuration,
			})
		}
		category := result.Get("13").String()
		nearbyPlace.Name = name
		nearbyPlace.Category = category
		nearbyPlace.GpsCoordinates = GpsCoordinate{
			Latitude:  result.Get("8.0").Float(),
			Longitude: result.Get("8.1").Float(),
		}
		nearbyPlaces = append(nearbyPlaces, nearbyPlace)
	}

	beforeTaxesFees := a.Get("6.2.1.0").String()
	beforeTaxesFeesStr := strings.Replace(beforeTaxesFees, "$", "", -1)
	beforeTaxesFeesStr = strings.Replace(beforeTaxesFeesStr, ",", "", -1)
	extractedBeforeTaxesFees, _ := strconv.Atoi(beforeTaxesFeesStr)

	lowest := a.Get("6.2.1.1").String()
	lowestStr := strings.Replace(lowest, "$", "", -1)
	lowestStr = strings.Replace(lowestStr, ",", "", -1)
	extractedLowest, _ := strconv.Atoi(lowestStr)

	ratePerNight := PriceRate{
		Lowest:                   lowest,
		ExtractedLowest:          extractedLowest,
		BeforeTaxesFees:          beforeTaxesFees,
		ExtractedBeforeTaxesFees: extractedBeforeTaxesFees,
	}

	totalBeforeTaxesFees := a.Get("6.2.9.0").String()
	totalBeforeTaxesFeesStr := strings.Replace(totalBeforeTaxesFees, "$", "", -1)
	totalBeforeTaxesFeesStr = strings.Replace(totalBeforeTaxesFeesStr, ",", "", -1)
	extractedTotalBeforeTaxesFees, _ := strconv.Atoi(totalBeforeTaxesFeesStr)

	totalLowest := a.Get("6.2.9.1").String()
	totalLowestStr := strings.Replace(totalLowest, "$", "", -1)
	totalLowestStr = strings.Replace(totalLowestStr, ",", "", -1)
	extractedTotalLowest, _ := strconv.Atoi(totalLowestStr)
	totalRate := PriceRate{
		Lowest:                   totalLowest,
		ExtractedLowest:          extractedTotalLowest,
		BeforeTaxesFees:          totalBeforeTaxesFees,
		ExtractedBeforeTaxesFees: extractedTotalBeforeTaxesFees,
	}

	gpsCoordinates := GpsCoordinate{
		Latitude:  a.Get("2.0.0").Float(),
		Longitude: a.Get("2.0.1").Float(),
	}

	description := a.Get("11.3.1").String()

	propertyToken := a.Get("20").String()

	link := a.Get("19.8.0.4").String()

	name := a.Get("1").String()

	hotelDetail = HotelDetail{
		Type:              "vacation rental",
		Name:              name,
		Description:       description,
		Link:              link,
		PropertyToken:     propertyToken,
		GpsCoordinates:    gpsCoordinates,
		RatePerNight:      ratePerNight,
		TotalRate:         totalRate,
		Prices:            price,
		NearbyPlaces:      nearbyPlaces,
		Images:            images,
		OverallRating:     overallRating,
		Reviews:           reviews,
		Amenities:         amenities,
		ExcludedAmenities: excludedAmenities,
		EssentialInfo:     essentialInfo,
	}
	//marshal, _ := json.Marshal(hotelDetail)
	//create, _ := os.Create("hotelDetail.json")
	//create.Write(marshal)
	//create.Close()
	return hotelDetail
}
func getAll(q string) HotelDetail {
	var (
		hotelDetail HotelDetail
	)

	a := gjson.Parse(q).Get("7").Array()[0].Get("441552390") //gps_coordinates
	propertyToken := a.Get("20").String()
	link := a.Get("2.29.2").String()
	hotelDetail.PropertyToken = propertyToken
	hotelDetail.Link = link
	ratePerNight := a.Get("6")
	beforeTaxesFees := ratePerNight.Get("2.10.0").String()
	beforeTaxesFeesStr := strings.Replace(beforeTaxesFees, "$", "", -1)
	beforeTaxesFeesStr = strings.Replace(beforeTaxesFeesStr, ",", "", -1)
	extractedBeforeTaxesFees, _ := strconv.Atoi(beforeTaxesFeesStr)

	lowest := ratePerNight.Get("2.10.1").String()
	lowestStr := strings.Replace(lowest, "$", "", -1)
	lowestStr = strings.Replace(lowestStr, ",", "", -1)
	extractedLowest, _ := strconv.Atoi(lowestStr)

	totalBeforeTaxesFees := ratePerNight.Get("2.9.0").String()
	totalBeforeTaxesFeesStr := strings.Replace(totalBeforeTaxesFees, "$", "", -1)
	totalBeforeTaxesFeesStr = strings.Replace(totalBeforeTaxesFeesStr, ",", "", -1)
	extractedTotalBeforeTaxesFees, _ := strconv.Atoi(totalBeforeTaxesFeesStr)

	totalLowest := ratePerNight.Get("2.9.1").String()
	totalLowestStr := strings.Replace(totalLowest, "$", "", -1)
	totalLowestStr = strings.Replace(totalLowestStr, ",", "", -1)
	extractedTotalLowest, _ := strconv.Atoi(totalLowestStr)

	hotelDetail.RatePerNight = PriceRate{
		Lowest:                   lowest,
		ExtractedLowest:          extractedLowest,
		BeforeTaxesFees:          beforeTaxesFees,
		ExtractedBeforeTaxesFees: extractedBeforeTaxesFees,
	}
	hotelDetail.TotalRate = PriceRate{
		Lowest:                   totalLowest,
		ExtractedLowest:          extractedTotalLowest,
		BeforeTaxesFees:          totalBeforeTaxesFees,
		ExtractedBeforeTaxesFees: extractedTotalBeforeTaxesFees,
	}

	overallRating := a.Get("7.0.0").Float()
	reviews := a.Get("7.0.1").Float()
	hotelDetail.OverallRating = overallRating
	hotelDetail.Reviews = reviews

	// ratings
	for _, result := range a.Get("7.1.0").Array() {
		stars := result.Get("0").Int()
		count := result.Get("2").Int()
		hotelDetail.Ratings = append(hotelDetail.Ratings, DetailRating{
			Stars: int(stars),
			Count: int(count),
		})
	}

	//reviews_breakdown
	for _, result := range a.Get("7.9").Array() {
		name := result.Get("6").String()
		description := result.Get("4").String()
		totalMentioned := result.Get("2").Int()
		negative := result.Get("8").Int()
		positive := result.Get("7").Int()
		neutral := totalMentioned - negative - positive
		hotelDetail.ReviewsBreakdown = append(hotelDetail.ReviewsBreakdown, ReviewsBreakdown{
			Name:           name,
			Description:    description,
			TotalMentioned: int(totalMentioned),
			Negative:       int(negative),
			Positive:       int(positive),
			Neutral:        int(neutral),
		})
	}

	// images
	for _, result := range a.Get("5.1").Array() {
		thumbnail := result.Get("1.0").String()
		originalImage := result.Get("1.4.1").String()
		hotelDetail.Images = append(hotelDetail.Images, DetailImage{
			Thumbnail:     thumbnail,
			OriginalImage: originalImage,
		})
	}

	hotelClass := a.Get("3.0").String()
	extractedHotelClass := a.Get("3.1").Int()
	hotelDetail.HotelClass = hotelClass
	hotelDetail.ExtractedHotelClass = int(extractedHotelClass)

	name := a.Get("1").String()
	description := a.Get("11.1").Array()
	hotelDetail.Name = name
	for _, result := range description {
		hotelDetail.Description = hotelDetail.Description + result.String()
	}

	gpsCoordinatesLatitude := a.Get("2.0.0")
	gpsCoordinatesLongitude := a.Get("2.0.1")
	address := a.Get("2.1.0.0.0")
	phone := a.Get("2.2.0")
	phoneLink := a.Get("2.2.1")
	hotelDetail.GpsCoordinates = GpsCoordinate{
		Latitude:  gpsCoordinatesLatitude.Float(),
		Longitude: gpsCoordinatesLongitude.Float(),
	}
	hotelDetail.Address = address.String()
	hotelDetail.Phone = phone.String()
	hotelDetail.PhoneLink = phoneLink.String()

	checkInTime := a.Get("2.17.0").String()
	checkOutTime := a.Get("2.17.1").String()
	hotelDetail.CheckInTime = checkInTime
	hotelDetail.CheckOutTime = checkOutTime

	for _, result := range a.Get("2.19.0.1.2").Array() { //
		hotelDetail.NearbyPlaces = append(hotelDetail.NearbyPlaces, getNearbyPlaces(result.String()))
	}

	for _, result := range a.Get("6.2.21").Array() {
		hotelDetail.Prices = append(hotelDetail.Prices, getPrices(result.String()))
	}

	for _, result := range a.Get("6.2.2").Array() { // 6.2.2
		hotelDetail.FeaturedPrices = append(hotelDetail.FeaturedPrices, getFeaturedPrices(result.String()))
	}
	hotelDetail.Type = "hotel"
	return hotelDetail
}

func getFeaturedPrices(s string) DetailFeaturedPrice {
	var (
		featuredPrice DetailFeaturedPrice
	)
	source := gjson.Parse(s).Get("0.0").String()
	link := fmt.Sprintf("https://www.google.com%s", gjson.Parse(s).Get("0.2").String())
	logo := gjson.Parse(s).Get("0.3.0").String()
	remarks := gjson.Parse(s).Get("0.6").Array()
	featuredPrice.Source = source
	featuredPrice.Link = link
	featuredPrice.Logo = logo
	for _, remark := range remarks {
		featuredPrice.Remark = append(featuredPrice.Remark, remark.Get("0").String())
	}
	for _, result := range gjson.Parse(s).Get("7").Array() {
		featuredPrice.Rooms = append(featuredPrice.Rooms, getFeaturedPricesRoom(result.String()))
	}
	return featuredPrice
}

func getFeaturedPricesRoom(s string) RoomInfo {
	result := gjson.Parse(s)
	name := result.Get("0").String()
	link := fmt.Sprintf("https://www.google.com%s", result.Get("2.0.0").String())
	numGuests := result.Get("2.0.1.0").Int()

	ratePerNightBeforeTaxesFees := result.Get("2.0.4.0").String()
	beforeTaxesFeesStr := strings.Replace(ratePerNightBeforeTaxesFees, "$", "", -1)
	beforeTaxesFeesStr = strings.Replace(beforeTaxesFeesStr, ",", "", -1)
	extractedBeforeTaxesFees, _ := strconv.Atoi(beforeTaxesFeesStr)

	ratePerNightLowest := result.Get("2.0.4.1").String()
	lowestStr := strings.Replace(ratePerNightLowest, "$", "", -1)
	lowestStr = strings.Replace(lowestStr, ",", "", -1)
	extractedLowest, _ := strconv.Atoi(lowestStr)

	totalRateBeforeTaxesFees := result.Get("2.0.5.0").String()
	totalBeforeTaxesFeesStr := strings.Replace(totalRateBeforeTaxesFees, "$", "", -1)
	totalBeforeTaxesFeesStr = strings.Replace(totalBeforeTaxesFeesStr, ",", "", -1)
	extractedTotalBeforeTaxesFees, _ := strconv.Atoi(totalBeforeTaxesFeesStr)

	totalRateBeforeLowest := result.Get("2.0.5.1").String()
	totalLowestStr := strings.Replace(totalRateBeforeLowest, "$", "", -1)
	totalLowestStr = strings.Replace(totalLowestStr, ",", "", -1)
	extractedTotalLowest, _ := strconv.Atoi(totalLowestStr)

	return RoomInfo{
		Name:      name,
		Link:      link,
		NumGuests: int(numGuests),
		RatePerNight: PriceRate{
			Lowest:                   ratePerNightLowest,
			ExtractedLowest:          extractedLowest,
			BeforeTaxesFees:          ratePerNightBeforeTaxesFees,
			ExtractedBeforeTaxesFees: extractedBeforeTaxesFees,
		},
		TotalRate: PriceRate{
			Lowest:                   totalRateBeforeLowest,
			ExtractedLowest:          extractedTotalLowest,
			BeforeTaxesFees:          totalRateBeforeTaxesFees,
			ExtractedBeforeTaxesFees: extractedTotalBeforeTaxesFees,
		},
	}
}

func getPrices(s string) DetailPrice {
	result := gjson.Parse(s)
	source := result.Get("0.0").String()
	link := fmt.Sprintf("https://www.google.com%s", result.Get("0.2").String())
	logo := result.Get("0.3.0").String()
	logo = strings.Replace(logo, "//", "", -1)
	official := result.Get("0.5").Bool()
	numGuests := result.Get("12.12.0.0").Int()

	ratePerNightLowest := result.Get("12.4.0").String()
	ratePerNightLowestStr := strings.Replace(ratePerNightLowest, "$", "", -1)
	ratePerNightLowestStr = strings.Replace(ratePerNightLowestStr, ",", "", -1)
	ratePerNightLowestBeforeTaxesFees, _ := strconv.Atoi(ratePerNightLowestStr)

	totalRateLowest := result.Get("12.5.0").String()
	totalRateLowestStr := strings.Replace(totalRateLowest, "$", "", -1)
	totalRateLowestStr = strings.Replace(totalRateLowestStr, ",", "", -1)
	totalRateLowestBeforeTaxesFees, _ := strconv.Atoi(totalRateLowestStr)

	return DetailPrice{
		Source:    source,
		Link:      link,
		Logo:      logo,
		Official:  official,
		NumGuests: int(numGuests),
		RatePerNight: PriceRate{
			Lowest:          ratePerNightLowest,
			ExtractedLowest: ratePerNightLowestBeforeTaxesFees,
		},
		TotalRate: PriceRate{
			Lowest:          totalRateLowest,
			ExtractedLowest: totalRateLowestBeforeTaxesFees,
		},
	}
}

type transportationsType int64

var (
	transportationsTypeMapping = map[transportationsType]string{
		0: "taxi",
		2: "walking",
		3: "public transport",
	}
)

func (t transportationsType) GetType() string {
	if _, ok := transportationsTypeMapping[t]; ok {
		return transportationsTypeMapping[t]
	}
	return fmt.Sprintf("%d", t)
}

func getNearbyPlaces(s string) DetailNearbyPlace {
	var (
		nearbyPlace DetailNearbyPlace
	)
	result := gjson.Parse(s)
	name := result.Get("0").String()
	thumbnail := result.Get("1").String()
	category := result.Get("13").String()
	link := fmt.Sprintf("https://www.google.com%s", result.Get("18").String())
	description := result.Get("7").String()
	gpsCoordinatesLatitude := result.Get("8.0").Float()
	gpsCoordinatesLongitude := result.Get("8.1").Float()
	nearbyPlace = DetailNearbyPlace{
		Name:        name,
		Thumbnail:   thumbnail,
		Category:    category,
		Link:        link,
		Description: description,
		GpsCoordinates: GpsCoordinate{
			Latitude:  gpsCoordinatesLatitude,
			Longitude: gpsCoordinatesLongitude,
		},
	}
	for _, r := range result.Get("2").Array() {
		transportationsTypeInt := transportationsType(r.Get("0").Int()) // 0 -->taxi  2 -->Walking  3-->Public transport
		transportationsDuration := r.Get("1").String()                  //
		nearbyPlace.Transportations = append(nearbyPlace.Transportations, struct {
			Type     string `json:"type"`
			Duration string `json:"duration"`
		}{
			Type:     transportationsTypeMapping[transportationsTypeInt],
			Duration: transportationsDuration,
		})
	}
	return nearbyPlace
}

type HotelDetail struct {
	Type                string                `json:"type"`
	Name                string                `json:"name"`
	Description         string                `json:"description"`
	Link                string                `json:"link"`
	PropertyToken       string                `json:"property_token"`
	Address             string                `json:"address,omitempty"`
	Phone               string                `json:"phone,omitempty"`
	PhoneLink           string                `json:"phone_link,omitempty"`
	GpsCoordinates      GpsCoordinate         `json:"gps_coordinates"`
	CheckInTime         string                `json:"check_in_time,omitempty"`
	CheckOutTime        string                `json:"check_out_time,omitempty"`
	RatePerNight        PriceRate             `json:"rate_per_night"`
	TotalRate           PriceRate             `json:"total_rate"`
	FeaturedPrices      []DetailFeaturedPrice `json:"featured_prices,omitempty"`
	Prices              []DetailPrice         `json:"prices"`
	NearbyPlaces        []DetailNearbyPlace   `json:"nearby_places"`
	HotelClass          string                `json:"hotel_class,omitempty"`
	ExtractedHotelClass int                   `json:"extracted_hotel_class,omitempty"`
	Images              []DetailImage         `json:"images"`
	OverallRating       float64               `json:"overall_rating"`
	Reviews             float64               `json:"reviews"`
	Ratings             []DetailRating        `json:"ratings,omitempty"`
	ReviewsBreakdown    []ReviewsBreakdown    `json:"reviews_breakdown,omitempty"`
	LocationRating      float64               `json:"location_rating,omitempty"`
	Amenities           []string              `json:"amenities,omitempty"`
	ExcludedAmenities   []string              `json:"excluded_amenities,omitempty"`
	EssentialInfo       []string              `json:"essential_info,omitempty"`
}
type ReviewsBreakdown struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	TotalMentioned int    `json:"total_mentioned"`
	Positive       int    `json:"positive"`
	Negative       int    `json:"negative"`
	Neutral        int    `json:"neutral"`
}
type DetailRating struct {
	Stars int `json:"stars"`
	Count int `json:"count"`
}
type DetailImage struct {
	Thumbnail     string `json:"thumbnail"`
	OriginalImage string `json:"original_image"`
}
type GpsCoordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type PriceRate struct {
	Lowest                   string `json:"lowest,omitempty"`
	ExtractedLowest          int    `json:"extracted_lowest,omitempty"`
	BeforeTaxesFees          string `json:"before_taxes_fees,omitempty"`
	ExtractedBeforeTaxesFees int    `json:"extracted_before_taxes_fees,omitempty"`
}

type DetailFeaturedPrice struct {
	Source string     `json:"source"`
	Link   string     `json:"link"`
	Logo   string     `json:"logo"`
	Remark []string   `json:"remark,omitempty"`
	Rooms  []RoomInfo `json:"rooms,omitempty"`
}

type RoomInfo struct {
	Name         string    `json:"name"`
	Link         string    `json:"link"`
	NumGuests    int       `json:"num_guests"`
	RatePerNight PriceRate `json:"rate_per_night"`
	TotalRate    PriceRate `json:"total_rate"`
}
type DetailPrice struct {
	Source       string    `json:"source"`
	Link         string    `json:"link"`
	Logo         string    `json:"logo"`
	Official     bool      `json:"official"`
	NumGuests    int       `json:"num_guests"`
	RatePerNight PriceRate `json:"rate_per_night"`
	TotalRate    PriceRate `json:"total_rate"`
}

type DetailNearbyPlace struct {
	Category        string `json:"category"`
	Name            string `json:"name"`
	Link            string `json:"link,omitempty"`
	Thumbnail       string `json:"thumbnail,omitempty"`
	Transportations []struct {
		Type     string `json:"type"`
		Duration string `json:"duration"`
	} `json:"transportations"`
	Description    string        `json:"description,omitempty"`
	GpsCoordinates GpsCoordinate `json:"gps_coordinates"`
}
