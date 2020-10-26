package usecases

import (
	"github.com/oraksil/azumma/internal/domain/models"
)

type UserFeedbackUseCase struct {
	FeedbackRepo models.UserFeedbackRepository
}

func (uc *UserFeedbackUseCase) CreateNewUserFeedback(feedback string) (*models.UserFeedback, error) {
	created, err := uc.FeedbackRepo.Save(&models.UserFeedback{Content: feedback})
	if err != nil {
		return nil, err
	}

	return created, nil
}
