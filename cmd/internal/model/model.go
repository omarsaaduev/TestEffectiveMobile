package model

type Person struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic,omitempty"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}

type NationalizeResponse struct {
	Name    string              `json:"name"`
	Country []CountryPrediction `json:"country"`
}

type CountryPrediction struct {
	CountryID   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}
