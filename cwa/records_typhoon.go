package cwa

import (
	"encoding/json"
	"fmt"
)

type TyphoonRecords struct {
	TropicalCyclones TropicalCyclones `json:"TropicalCyclones"`
}

type TropicalCyclones struct {
	TropicalCyclone []TropicalCyclone `json:"TropicalCyclone"`
}

type TropicalCyclone struct {
	Year           string       `json:"Year"`
	TyphoonName    string       `json:"TyphoonName"`
	CwaTyphoonName string       `json:"CwaTyphoonName"`
	CwaTdNo        string       `json:"CwaTdNo"`
	CwaTyNo        string       `json:"CwaTyNo"`
	AnalysisData   *CycloneData `json:"AnalysisData,omitempty"`
	ForecastData   *CycloneData `json:"ForecastData,omitempty"`
}

type CycloneData struct {
	Fix []CycloneFix `json:"Fix"`
}

type CycloneFix struct {
	DateTime            string `json:"DateTime"`
	CoordinateLongitude string `json:"CoordinateLongitude"`
	CoordinateLatitude  string `json:"CoordinateLatitude"`
	MaxWindSpeed        string `json:"MaxWindSpeed"`
	MaxGustSpeed        string `json:"MaxGustSpeed"`
	Pressure            string `json:"Pressure"`
	MovingSpeed         string `json:"MovingSpeed"`
	MovingDirection     string `json:"MovingDirection"`
}

func (r *Response) ParseTyphoonRecords() (*TyphoonRecords, error) {
	var rec TyphoonRecords
	if err := json.Unmarshal(r.Records, &rec); err != nil {
		return nil, fmt.Errorf("failed to parse typhoon records: %w", err)
	}
	return &rec, nil
}
