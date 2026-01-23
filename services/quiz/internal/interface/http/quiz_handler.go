package http

import (
	"net/http"

	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/quiz/internal/pkg"
	"github.com/ISAWASHUN/garbage-category-rule-quiz/services/quiz/internal/usecase"
	"github.com/gin-gonic/gin"
)

// 板橋区のmunicipality_id（現時点では固定しておく）
const defaultMunicipalityID = 1
const defaultQuestionCount = 5

type QuestionResponse struct {
	ID       int      `json:"id" example:"123"`
	ItemName string   `json:"item_name" example:"スプレー缶"`
	Choices  []string `json:"choices" example:"不燃,可燃,資源,有害ごみ"`
}

type GetQuestionsResponse struct {
	Questions []QuestionResponse `json:"questions"`
}

type PostAnswerRequest struct {
	QuestionID       int    `json:"question_id" binding:"required" example:"123"`
	SelectedCategory string `json:"selected_category" binding:"required" example:"不燃"`
}

type PostAnswerResponse struct {
	IsCorrect       bool   `json:"is_correct" example:"true"`
	CorrectCategory string `json:"correct_category" example:"不燃"`
	Notes           string `json:"notes" example:"中身を使い切ってから"`
	Remarks         string `json:"remarks" example:"スプレー缶は中身を使い切ってから不燃ごみとして出してください"`
	BulkGarbageFee  *int   `json:"bulk_garbage_fee" example:"500"`
}

type QuizHandler struct {
	quizUseCase usecase.QuizUseCase
}

func NewQuizHandler(quizUseCase usecase.QuizUseCase) *QuizHandler {
	return &QuizHandler{
		quizUseCase: quizUseCase,
	}
}

// GetQuestions はクイズの問題を取得するハンドラー
// @Summary クイズ問題取得
// @Description ランダムで5問のクイズ問題を取得する
// @Tags quiz
// @Accept json
// @Produce json
// @Success 200 {object} GetQuestionsResponse "成功時のレスポンス"
// @Failure 500 {object} ErrorResponse "サーバーエラー"
// @Router /quiz/questions [get]
func (h *QuizHandler) GetQuestions(c *gin.Context) {
	ctx := c.Request.Context()

	questions, err := h.quizUseCase.GenerateQuestions(ctx, defaultMunicipalityID, defaultQuestionCount)
	if err != nil {
		pkg.HandleError(c, pkg.NewInternalError("問題の取得に失敗しました", err))
		return
	}

	response := GetQuestionsResponse{
		Questions: make([]QuestionResponse, len(questions)),
	}

	for i, q := range questions {
		response.Questions[i] = QuestionResponse{
			ID:       q.ID,
			ItemName: q.ItemName,
			Choices:  q.Choices,
		}
	}

	c.JSON(http.StatusOK, response)
}

// PostAnswer は回答を送信するハンドラー
// @Summary 回答送信
// @Description クイズの回答を送信し、正誤と正解情報を取得する
// @Tags quiz
// @Accept json
// @Produce json
// @Param request body PostAnswerRequest true "回答リクエスト"
// @Success 200 {object} PostAnswerResponse "成功時のレスポンス"
// @Failure 400 {object} ErrorResponse "リクエストエラー"
// @Failure 404 {object} ErrorResponse "問題が見つからない"
// @Failure 500 {object} ErrorResponse "サーバーエラー"
// @Router /quiz/answer [post]
func (h *QuizHandler) PostAnswer(c *gin.Context) {
	ctx := c.Request.Context()

	var req PostAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.HandleError(c, pkg.NewBadRequestError("リクエストの形式が不正です", err))
		return
	}

	result, err := h.quizUseCase.CheckAnswer(ctx, req.QuestionID, req.SelectedCategory)
	if err != nil {
		pkg.HandleError(c, pkg.NewNotFoundError("指定された問題が見つかりません", err))
		return
	}

	response := PostAnswerResponse{
		IsCorrect:       result.IsCorrect,
		CorrectCategory: result.CorrectCategory,
		Notes:           result.Notes,
		Remarks:         result.Remarks,
		BulkGarbageFee:  result.BulkGarbageFee,
	}

	c.JSON(http.StatusOK, response)
}
