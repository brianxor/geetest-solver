package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/tylerdevx/geetest-solver/solver"
	"net/http"
	"time"
)

func V4PuzzleSolveHandler(ctx echo.Context) error {
	type body struct {
		WebsiteUrl string `json:"websiteUrl"`
		CaptchaId  string `json:"captchaId"`
		UserAgent  string `json:"userAgent"`
		Proxy      string `json:"proxy"`
	}

	reqBody := new(body)

	if err := ctx.Bind(reqBody); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "failed parsing request body",
		})
	}

	websiteUrl := reqBody.WebsiteUrl
	captchaId := reqBody.CaptchaId
	userAgent := reqBody.UserAgent
	proxy := reqBody.Proxy

	if websiteUrl == "" || captchaId == "" || userAgent == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "missing required fields",
		})
	}

	captchaSolver, err := solver.NewGeetestSolver(websiteUrl, captchaId, userAgent, proxy)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "failed creating solver",
		})
	}

	start := time.Now()

	solution, err := captchaSolver.SolveV4Puzzle()

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "failed solving captcha",
		})
	}

	solveTime := time.Since(start)

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"success":   true,
		"solveTime": fmt.Sprintf("%s", solveTime),
		"solution":  solution.Solution,
	})

}
