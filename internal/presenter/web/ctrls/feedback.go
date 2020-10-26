package ctrls

import (
	"net/http"

	"clevergo.tech/jsend"
	"github.com/gin-gonic/gin"
	"github.com/oraksil/azumma/internal/domain/usecases"
	"github.com/oraksil/azumma/internal/presenter/web"
	"github.com/oraksil/azumma/internal/presenter/web/ctrls/dto"
	"github.com/oraksil/azumma/internal/presenter/web/ctrls/helpers"
)

type UserFeedbackController struct {
	UserFeedbackUseCase *usecases.UserFeedbackUseCase
}

func (ctrl *UserFeedbackController) createNewFeedback(c *gin.Context) {
	sessionCtx := helpers.NewSessionCtx(c)
	if sessionCtx.Validate() != nil {
		c.JSON(http.StatusOK, jsend.NewFail(map[string]interface{}{
			"code":    400,
			"message": "invalid session",
		}))
		return
	}

	type JsonParams struct {
		Feedback string `json:"feedback"`
	}

	var jsonParams JsonParams
	c.BindJSON(&jsonParams)

	_, err := ctrl.UserFeedbackUseCase.CreateNewUserFeedback(jsonParams.Feedback)
	if err != nil {
		c.JSON(http.StatusOK, jsend.NewFail(map[string]interface{}{
			"code":    400,
			"message": err.Error(),
		}))
		return
	}

	c.JSON(http.StatusOK, jsend.New(dto.Empty()))
}

func (ctrl *UserFeedbackController) Routes() []web.Route {
	return []web.Route{
		{Spec: "POST /api/v1/feedbacks/new", Handler: ctrl.createNewFeedback},
	}
}
