package usecase

import (
	"context"
	"math/rand"

	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/quiz/internal/domain"
	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/quiz/internal/repository"
)

var categoryMasterList = []string{
	"可燃",
	"不燃",
	"資源",
	"粗大",
	"有害ごみ",
	"容器包装プラスチック",
	"製品プラスチック",
	"古布",
	"新聞・折込チラシ",
	"雑誌・本・雑がみ",
	"段ボール・茶色紙",
	"牛乳等紙パック",
	"缶",
	"びん",
	"ペットボトル",
	"燃やせるごみ",
	"燃やせないごみ",
	"粗大ごみ",
	"拠点",
	"不可",
	"処理困難物【市での収集は不可】",
	"パソコン【市での収集は不可】",
	"家電リサイクル法対象品【市での収集は不可】",
}

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
		return []Question{}, nil
	}

	// ランダムに選択
	selectedItems := selectRandomItems(allItems, count)

	questions := make([]Question, len(selectedItems))
	for i, item := range selectedItems {
		category, err := u.garbageCategoryRepo.GetByID(ctx, int(item.GarbageCategoryID))
		if err != nil {
			return nil, err
		}

		choices := generateChoices(category.Name)

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

	// インデックスをシャッフルしてから選択
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

func generateChoices(correctCategory string) []string {
	wrongCategories := make([]string, 0, len(categoryMasterList)-1)
	for _, cat := range categoryMasterList {
		if cat != correctCategory {
			wrongCategories = append(wrongCategories, cat)
		}
	}

	rand.Shuffle(len(wrongCategories), func(i, j int) {
		wrongCategories[i], wrongCategories[j] = wrongCategories[j], wrongCategories[i]
	})

	choices := make([]string, 4)
	choices[0] = correctCategory
	copy(choices[1:], wrongCategories[:3])

	rand.Shuffle(len(choices), func(i, j int) {
		choices[i], choices[j] = choices[j], choices[i]
	})

	return choices
}
