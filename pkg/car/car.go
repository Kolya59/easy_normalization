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
}

var Data []Car

func FillData() {
	Data = []Car{
		{
			Model:                   "",
			BrandName:               "",
			BrandCreatorCountry:     "",
			EngineModel:             "",
			EnginePower:             0,
			EngineVolume:            0,
			EngineType:              "",
			TransmissionModel:       "",
			TransmissionType:        "",
			TransmissionGearsNumber: 0,
		},
		{
			Model:                   "",
			BrandName:               "",
			BrandCreatorCountry:     "",
			EngineModel:             "",
			EnginePower:             0,
			EngineVolume:            0,
			EngineType:              "",
			TransmissionModel:       "",
			TransmissionType:        "",
			TransmissionGearsNumber: 0,
		},
		{
			Model:                   "",
			BrandName:               "",
			BrandCreatorCountry:     "",
			EngineModel:             "",
			EnginePower:             0,
			EngineVolume:            0,
			EngineType:              "",
			TransmissionModel:       "",
			TransmissionType:        "",
			TransmissionGearsNumber: 0,
		},
		{
			Model:                   "",
			BrandName:               "",
			BrandCreatorCountry:     "",
			EngineModel:             "",
			EnginePower:             0,
			EngineVolume:            0,
			EngineType:              "",
			TransmissionModel:       "",
			TransmissionType:        "",
			TransmissionGearsNumber: 0,
		},
		{
			Model:                   "",
			BrandName:               "",
			BrandCreatorCountry:     "",
			EngineModel:             "",
			EnginePower:             0,
			EngineVolume:            0,
			EngineType:              "",
			TransmissionModel:       "",
			TransmissionType:        "",
			TransmissionGearsNumber: 0,
		},
		{
			Model:                   "",
			BrandName:               "",
			BrandCreatorCountry:     "",
			EngineModel:             "",
			EnginePower:             0,
			EngineVolume:            0,
			EngineType:              "",
			TransmissionModel:       "",
			TransmissionType:        "",
			TransmissionGearsNumber: 0,
		},
		{
			Model:                   "",
			BrandName:               "",
			BrandCreatorCountry:     "",
			EngineModel:             "",
			EnginePower:             0,
			EngineVolume:            0,
			EngineType:              "",
			TransmissionModel:       "",
			TransmissionType:        "",
			TransmissionGearsNumber: 0,
		},
	}
}
