package domain

type Storage interface {
	Save(items []GarbageItem) error
}
