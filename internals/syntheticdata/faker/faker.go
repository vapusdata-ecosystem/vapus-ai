package faker

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/go-faker/faker/v4"
)

type Provider string

const (
	ID                  Provider = "uuid_digit"
	HyphenatedID        Provider = "uuid_hyphenated"
	EmailTag            Provider = "email"
	MacAddressTag       Provider = "mac_address"
	OrganizationNameTag Provider = "Organization_name"
	UserNameTag         Provider = "username"
	URLTag              Provider = "url"
	IPV4Tag             Provider = "ipv4"
	IPV6Tag             Provider = "ipv6"
	PASSWORD            Provider = "password"
	JWT                 Provider = "jwt"
	LATITUDE            Provider = "lat"
	LONGITUDE           Provider = "long"
	RealAddress         Provider = "real_address"
	CreditCardNumber    Provider = "cc_number"
	CreditCardType      Provider = "cc_type"
	PhoneNumber         Provider = "phone_number"
	TollFreeNumber      Provider = "toll_free_number"
	E164PhoneNumberTag  Provider = "e_164_phone_number"
	TitleMaleTag        Provider = "title_male"
	TitleFemaleTag      Provider = "title_female"
	FirstNameTag        Provider = "first_name"
	FirstNameMaleTag    Provider = "first_name_male"
	FirstNameFemaleTag  Provider = "first_name_female"
	LastNameTag         Provider = "last_name"
	NAME                Provider = "name"
	ChineseFirstNameTag Provider = "chinese_first_name"
	ChineseLastNameTag  Provider = "chinese_last_name"
	ChineseNameTag      Provider = "chinese_name"
	GENDER              Provider = "gender"
	UnixTimeTag         Provider = "unix_time"
	DATE                Provider = "date"
	TIME                Provider = "time"
	MonthNameTag        Provider = "month_name"
	YEAR                Provider = "year"
	DayOfWeekTag        Provider = "day_of_week"
	DayOfMonthTag       Provider = "day_of_month"
	TIMESTAMP           Provider = "timestamp"
	CENTURY             Provider = "century"
	TIMEZONE            Provider = "timezone"
	TimePeriodTag       Provider = "time_period"
	WORD                Provider = "word"
	SENTENCE            Provider = "sentence"
	PARAGRAPH           Provider = "paragraph"
	CurrencyTag         Provider = "currency"
	//AmountTag                 Provider = "amount"
	AmountWithCurrencyTag Provider = "amount_with_currency"
	// BloodTypeTag              Provider = "blood_type"
	// CountryInfoTag            Provider = "country_info"
	// UserAgentTag              Provider = "user_agent"
)

var FakerMap = map[Provider]func() any{
	ID:                    func() any { return faker.UUIDDigit() },
	HyphenatedID:          func() any { return faker.UUIDHyphenated() },
	EmailTag:              func() any { return faker.Email() },
	MacAddressTag:         func() any { return faker.MacAddress() },
	OrganizationNameTag:   func() any { return faker.DomainName() },
	UserNameTag:           func() any { return faker.Username() },
	URLTag:                func() any { return faker.URL() },
	IPV4Tag:               func() any { return faker.IPv4() },
	IPV6Tag:               func() any { return faker.IPv6() },
	PASSWORD:              func() any { return faker.Password() },
	JWT:                   func() any { return faker.Jwt() },
	LATITUDE:              func() any { return faker.Latitude() },
	LONGITUDE:             func() any { return faker.Longitude() },
	RealAddress:           func() any { return faker.GetRealAddress() },
	CreditCardNumber:      func() any { return faker.CCNumber() },
	CreditCardType:        func() any { return faker.CCType() },
	PhoneNumber:           func() any { return faker.Phonenumber() },
	TollFreeNumber:        func() any { return faker.TollFreePhoneNumber() },
	E164PhoneNumberTag:    func() any { return faker.E164PhoneNumber() },
	TitleMaleTag:          func() any { return faker.TitleMale() },
	TitleFemaleTag:        func() any { return faker.TitleFemale() },
	FirstNameTag:          func() any { return faker.FirstName() },
	FirstNameMaleTag:      func() any { return faker.FirstNameMale() },
	FirstNameFemaleTag:    func() any { return faker.FirstNameFemale() },
	LastNameTag:           func() any { return faker.LastName() },
	NAME:                  func() any { return faker.Name() },
	GENDER:                func() any { return faker.Gender() },
	UnixTimeTag:           func() any { return faker.UnixTime() },
	DATE:                  func() any { return faker.Date() },
	TIME:                  func() any { return faker.TimeString() },
	MonthNameTag:          func() any { return faker.MonthName() },
	YEAR:                  func() any { return faker.YearString() },
	DayOfWeekTag:          func() any { return faker.DayOfWeek() },
	DayOfMonthTag:         func() any { return faker.DayOfMonth() },
	TIMESTAMP:             func() any { return faker.Timestamp() },
	CENTURY:               func() any { return faker.Century() },
	TIMEZONE:              func() any { return faker.Timezone() },
	TimePeriodTag:         func() any { return faker.Timeperiod() },
	WORD:                  func() any { return faker.Word() },
	SENTENCE:              func() any { return faker.Sentence() },
	PARAGRAPH:             func() any { return faker.Paragraph() },
	CurrencyTag:           func() any { return faker.Currency() },
	AmountWithCurrencyTag: func() any { return faker.AmountWithCurrency() },
	ChineseFirstNameTag:   func() any { return faker.ChineseFirstName() },
	ChineseLastNameTag:    func() any { return faker.ChineseLastName() },
	ChineseNameTag:        func() any { return faker.ChineseName() },
	// UserAgentTag:              func() any { return faker.GetUserAgent() }, returns interface
	// BloodTypeTag:              func() any { return faker.GetBlood() },returns interface
}

type Data struct {
	Provider Provider
	Type     string
	Count    int64
	Result   []any
}

func Generate(request []*Data) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(request))

	for _, req := range request {
		wg.Add(1)
		go func(req *Data) {
			defer wg.Done()

			fakerFunc, exists := FakerMap[req.Provider]
			if !exists {
				errCh <- fmt.Errorf("invalid provider: %s", req.Provider)
				return
			}

			req.Result = make([]any, req.Count)
			for i := int64(0); i < req.Count; i++ {
				req.Result[i] = fakerFunc()
			}

			if req.Count > 0 {
				req.Type = reflect.TypeOf(req.Result[0]).String()
			}
		}(req)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		return err
	}
	return nil
}
