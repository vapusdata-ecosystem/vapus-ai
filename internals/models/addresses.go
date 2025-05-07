package models

import mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"

type Address struct {
	StreetAddress1 string `json:"streetAddress1" yaml:"streetAddress1"`
	StreetAddress2 string `json:"streetAddress2" yaml:"streetAddress2"`
	City           string `json:"city" yaml:"city"`
	State          string `json:"state" yaml:"state"`
	ZipCode        string `json:"zipCode" yaml:"zipCode"`
	Country        string `json:"country" yaml:"country"`
	Others         string `json:"others" yaml:"others"`
}

func (a *Address) GetStreetAddress1() string {
	return a.StreetAddress1
}

func (a *Address) GetStreetAddress2() string {
	return a.StreetAddress2
}

func (a *Address) GetCity() string {
	return a.City
}

func (a *Address) GetState() string {
	return a.State
}

func (a *Address) GetZipCode() string {
	return a.ZipCode
}

func (a *Address) GetCountry() string {
	return a.Country
}

func (a *Address) GetOthers() string {
	return a.Others
}

func (a *Address) ConvertFromPb(address *mpb.Address) *Address {
	a.StreetAddress1 = address.StreetAddress1
	a.StreetAddress2 = address.StreetAddress2
	a.City = address.City
	a.State = address.State
	a.ZipCode = address.ZipCode
	a.Country = address.Country
	a.Others = address.Others
	return a
}

func (a *Address) ConvertToPb() *mpb.Address {
	return &mpb.Address{
		StreetAddress1: a.StreetAddress1,
		StreetAddress2: a.StreetAddress2,
		City:           a.City,
		State:          a.State,
		ZipCode:        a.ZipCode,
		Country:        a.Country,
		Others:         a.Others,
	}
}
