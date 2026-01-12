package domain

import "time"

type MotorcycleStatus string

const (
	MotorcycleStatusDraft     MotorcycleStatus = "draft"
	MotorcycleStatusAvailable MotorcycleStatus = "available"
	MotorcycleStatusReserved  MotorcycleStatus = "reserved"
	MotorcycleStatusSold      MotorcycleStatus = "sold"
)

type MotorcycleData struct {
	Mileage      *int    `json:"mileage,omitempty"`
	MileageUnit  string  `json:"mileage_unit,omitempty"`
	Volume       *int    `json:"volume,omitempty"`
	VolumeUnit   string  `json:"volume_unit,omitempty"`
	FrameNumber  string  `json:"frame_number,omitempty"`
	ArrivalDate  string  `json:"arrival_date,omitempty"`
}

type Motorcycle struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Price       float64           `json:"price"`
	Currency    string            `json:"currency"`
	Data        *MotorcycleData   `json:"data,omitempty"`
	Status      MotorcycleStatus  `json:"status"`
	SourceURL   string            `json:"sourceUrl"`
	Photos      []*MotorcyclePhoto `json:"photos,omitempty"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

type MotorcyclePhoto struct {
	ID          string    `json:"id"`
	MotorcycleID string   `json:"motorcycleId"`
	S3URL       string    `json:"s3Url"`
	Order       int       `json:"order"`
	CreatedAt   time.Time `json:"createdAt"`
}

type CreateMotorcycle struct {
	Title       string           `json:"title"`
	Price       float64          `json:"price"`
	Currency    string           `json:"currency"`
	Data        *MotorcycleData  `json:"data,omitempty"`
	Status      MotorcycleStatus `json:"status"`
	SourceURL   string           `json:"sourceUrl"`
	PhotoURLs   []string         `json:"photoUrls"`
}

type PatchMotorcycle struct {
	Title       *string           `json:"title,omitempty"`
	Price       *float64          `json:"price,omitempty"`
	Currency    *string           `json:"currency,omitempty"`
	Data        *MotorcycleData   `json:"data,omitempty"`
	Status      *MotorcycleStatus `json:"status,omitempty"`
}

type FilterMotorcycle struct {
	ID     *string           `json:"id,omitempty"`
	Status *MotorcycleStatus `json:"status,omitempty"`
	Title  *string           `json:"title,omitempty"`
	MinPrice *float64        `json:"minPrice,omitempty"`
	MaxPrice *float64        `json:"maxPrice,omitempty"`
	
	IncludePhotos bool `json:"includePhotos"`
}

type CreateMotorcycleFromURL struct {
	URL string `json:"url"`
}

type ParsedMotorcycleData struct {
	Name     string
	Year     int
	Mileage  int
	Volume   int
	FrameNum string
	Images   []string
}

