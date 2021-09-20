package esios

type PriceData struct {
	Dia  string
	Hora string
	PCB  string
}

type PriceResponse struct {
	PVPC []PriceData
}

func GetPrice() PriceResponse {
	resp := PriceResponse{}
	err := GetJson(BaseURL, &resp)
	if err != nil {
		panic(err)
	}

	return resp
}
