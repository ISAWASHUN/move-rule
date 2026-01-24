package domain

type APIClient interface {
	FetchData(url string) ([]GarbageItem, error)
}
