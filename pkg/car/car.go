package car

type Car struct {
	Model                   string `json:"model"`
	BrandName               string `json:"brand_name"`
	BrandCreatorCountry     string `json:"brand_creator_country"`
	EngineModel             string `json:"engine_model"`
	EnginePower             int    `json:"engine_power"`
	EngineVolume            int    `json:"engine_volume"`
	EngineType              string `json:"engine_type"`
	TransmissionModel       string `json:"transmission_model"`
	TransmissionType        string `json:"transmission_type"`
	TransmissionGearsNumber int    `json:"transmission_gears_number"`
	WheelModel              string `json:"wheel_model"`
	WheelRadius             int    `json:"wheel_radius"`
	WheelColor              string `json:"wheel_color"`
	Price                   int    `json:"price"`
}

var Data []Car

func FillData() {
	Data = []Car{
		{
			Model:                   "2114",
			BrandName:               "LADA",
			BrandCreatorCountry:     "Russia",
			EngineModel:             "V123",
			EnginePower:             80,
			EngineVolume:            16,
			EngineType:              "L4",
			TransmissionModel:       "M123",
			TransmissionType:        "M",
			TransmissionGearsNumber: 5,
			WheelModel:              "Luchshie kolesa Rossii",
			WheelRadius:             13,
			WheelColor:              "Black",
			Price:                   120000,
		},
		{
			Model:                   "2115",
			BrandName:               "LADA",
			BrandCreatorCountry:     "Russia",
			EngineModel:             "V124",
			EnginePower:             100,
			EngineVolume:            18,
			EngineType:              "L4",
			TransmissionModel:       "M123",
			TransmissionType:        "M",
			TransmissionGearsNumber: 5,
			WheelModel:              "Luchshie kolesa Rossii",
			WheelRadius:             13,
			WheelColor:              "Black",
			Price:                   150000,
		},
		{
			Model:                   "Rio",
			BrandName:               "Kia",
			BrandCreatorCountry:     "Korea",
			EngineModel:             "V14234",
			EnginePower:             100,
			EngineVolume:            90,
			EngineType:              "V4",
			TransmissionModel:       "A123",
			TransmissionType:        "A",
			TransmissionGearsNumber: 4,
			WheelModel:              "Luchie kolesa Kitaya",
			WheelRadius:             15,
			WheelColor:              "Red",
			Price:                   400000,
		},
		{
			Model:                   "Sportage",
			BrandName:               "Kia",
			BrandCreatorCountry:     "Korea",
			EngineModel:             "V14234",
			EnginePower:             100,
			EngineVolume:            90,
			EngineType:              "V4",
			TransmissionModel:       "A1234",
			TransmissionType:        "A",
			TransmissionGearsNumber: 5,
			WheelModel:              "Luchie kolesa Kitaya",
			WheelRadius:             15,
			WheelColor:              "Red",
			Price:                   400000,
		},
		{
			Model:                   "A500",
			BrandName:               "Mercedes",
			BrandCreatorCountry:     "Germany",
			EngineModel:             "E1488",
			EnginePower:             300,
			EngineVolume:            50,
			EngineType:              "V12",
			TransmissionModel:       "R123",
			TransmissionType:        "A",
			TransmissionGearsNumber: 8,
			WheelModel:              "Luchshie kolesa Armenii",
			WheelRadius:             20,
			WheelColor:              "Green",
			Price:                   3000000,
		},
	}
}
