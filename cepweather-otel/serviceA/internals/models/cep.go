package models

type CepInput struct {
	Cep string `json:"cep" validate:"required,len=8,number=true"`
}

type WeatherData struct {
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
	TempK float64 `json:"temp_k"`
}
