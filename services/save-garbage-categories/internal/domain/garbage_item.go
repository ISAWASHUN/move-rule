package domain

type GarbageItem struct {
	Row                      int    `json:"Row"`
	MunicipalityCode         int    `json:"MunicipalityCode"`
	ID                       string `json:"ID"`
	MunicipalityName         string `json:"MunicipalityName"`
	AreaName                 string `json:"AreaName"`
	ItemName                 string `json:"ItemName"`
	ItemNameKana             string `json:"ItemNameKana"`
	ItemNameEnglish          string `json:"ItemNameEnglish"`
	Category                 string `json:"Category"`
	Notes                    string `json:"Notes"`
	Remarks                  string `json:"Remarks"`
	BulkGarbageCollectionFee string `json:"BulkGarbageCollectionFee"`
}

type FileReader interface {
	Read(path string) ([]GarbageItem, error)
}

type MunicipalityRepository interface {
	GetOrCreate(code int, name string) (int, error)
	Truncate() error
}

type GarbageCategoryRepository interface {
	GetOrCreate(name string) (int, error)
	Truncate() error
}

type GarbageItemRepository interface {
	Create(item GarbageItem, municipalityID, garbageCategoryID int, bulkGarbageFee *int) error
	Truncate() error
}
