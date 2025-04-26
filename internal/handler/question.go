package handler

import (
	"strconv"

	"github.com/CRobinDev/BCCGembira_Nusastra/internal/service"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/response"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type IQuestionHandler interface {
	GetAllQuestionsByQuizID() fiber.Handler
}

type questionHandler struct {
	questionService service.IQuestionService
	val             validator.Validator
}

func NewQuestionHandler(questionService service.IQuestionService, val validator.Validator) IQuestionHandler {
	return &questionHandler{
		questionService: questionService,
		val:             val,
	}
}

func (qh *questionHandler) GetAllQuestionsByQuizID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		quizID := c.Params("quizID")
		if err := qh.val.Validate(quizID); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		quizIDInt, err := strconv.ParseUint(quizID, 10, 64)
		if err != nil {
			return err 
		}

		questions, err := qh.questionService.GetAllQuestionsByQuizID(c.Context(), quizIDInt)
		if err != nil {
			return err
		}

		return response.Success(c, "questions", questions)
	}
}
