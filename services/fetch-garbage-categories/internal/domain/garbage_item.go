package domain

type GarbageItem struct {
	Row                      int
	MunicipalityCode         int
	ID                       string
	MunicipalityName         string
	AreaName                 string
	ItemName                 string
	ItemNameKana             string
	ItemNameEnglish          string
	Category                 string
	Notes                    string
	Remarks                  string
	BulkGarbageCollectionFee string
}
