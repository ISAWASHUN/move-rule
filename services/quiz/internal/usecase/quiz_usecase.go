package usecase

import (
	"context"
	"math/rand"

	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/quiz/internal/domain"
	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/quiz/internal/infrastructure/repository"
)

type Question struct {
	ID       int
	ItemName string
	Choices  []string
}

type AnswerResult struct {
	IsCorrect       bool
	CorrectCategory string
	Notes           string
	Remarks         string
	BulkGarbageFee  *int
}

type QuizUseCase interface {
	GenerateQuestions(ctx context.Context, municipalityID int, count int) ([]Question, error)
	CheckAnswer(ctx context.Context, questionID int, selectedCategory string) (*AnswerResult, error)
}

type quizUseCase struct {
	garbageItemRepo     repository.GarbageItemRepository
	garbageCategoryRepo repository.GarbageCategoryRepository
	municipalityRepo    repository.MunicipalityRepository
}

func NewQuizUseCase(
	garbageItemRepo repository.GarbageItemRepository,
	garbageCategoryRepo repository.GarbageCategoryRepository,
	municipalityRepo repository.MunicipalityRepository,
) QuizUseCase {
	return &quizUseCase{
		garbageItemRepo:     garbageItemRepo,
		garbageCategoryRepo: garbageCategoryRepo,
		municipalityRepo:    municipalityRepo,
	}
}

func (u *quizUseCase) GenerateQuestions(ctx context.Context, municipalityID int, count int) ([]Question, error) {
	allItems, err := u.garbageItemRepo.GetByMunicipalityID(ctx, municipalityID)
	if err != nil {
		return nil, err
	}

	if len(allItems) == 0 {
		// データが存在しない場合は空配列を返す
		// データベースにデータを投入する必要があります
		return []Question{}, nil
	}

	allCategories, err := u.garbageCategoryRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	categoryNames := make([]string, len(allCategories))
	for i, cat := range allCategories {
		categoryNames[i] = cat.Name
	}

	selectedItems := selectRandomItems(allItems, count)

	questions := make([]Question, len(selectedItems))
	for i, item := range selectedItems {
		category, err := u.garbageCategoryRepo.GetByID(ctx, int(item.GarbageCategoryID))
		if err != nil {
			return nil, err
		}

		choices := generateChoices(category.Name, categoryNames)

		questions[i] = Question{
			ID:       int(item.ID),
			ItemName: item.ItemName,
			Choices:  choices,
		}
	}

	return questions, nil
}

func selectRandomItems(items []domain.GarbageItem, count int) []domain.GarbageItem {
	if len(items) <= count {
		return items
	}

	indices := make([]int, len(items))
	for i := range indices {
		indices[i] = i
	}

	rand.Shuffle(len(indices), func(i, j int) {
		indices[i], indices[j] = indices[j], indices[i]
	})

	result := make([]domain.GarbageItem, count)
	for i := 0; i < count; i++ {
		result[i] = items[indices[i]]
	}

	return result
}

func (u *quizUseCase) CheckAnswer(ctx context.Context, questionID int, selectedCategory string) (*AnswerResult, error) {
	item, category, err := u.garbageItemRepo.GetByIDWithCategory(ctx, questionID)
	if err != nil {
		return nil, err
	}

	isCorrect := category.Name == selectedCategory

	var bulkGarbageFee *int
	if item.BulkGarbageFee > 0 {
		fee := item.BulkGarbageFee
		bulkGarbageFee = &fee
	}

	return &AnswerResult{
		IsCorrect:       isCorrect,
		CorrectCategory: category.Name,
		Notes:           item.Notes,
		Remarks:         item.Remarks,
		BulkGarbageFee:  bulkGarbageFee,
	}, nil
}

func generateChoices(correctCategory string, categoryList []string) []string {
	wrongCategories := make([]string, 0, len(categoryList)-1)
	for _, cat := range categoryList {
		if cat != correctCategory {
			wrongCategories = append(wrongCategories, cat)
		}
	}

	rand.Shuffle(len(wrongCategories), func(i, j int) {
		wrongCategories[i], wrongCategories[j] = wrongCategories[j], wrongCategories[i]
	})

	choices := make([]string, 4)
	choices[0] = correctCategory
	copyCount := 3
	if len(wrongCategories) < copyCount {
		copyCount = len(wrongCategories)
	}
	copy(choices[1:1+copyCount], wrongCategories[:copyCount])

	rand.Shuffle(len(choices), func(i, j int) {
		choices[i], choices[j] = choices[j], choices[i]
	})

	return choices
}
