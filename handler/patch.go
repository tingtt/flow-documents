package handler

import (
	"flow-documents/document"
	"flow-documents/flags"
	"flow-documents/jwt"
	"flow-documents/utils"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func Patch(c echo.Context) error {
	// Check `Content-Type`
	if !strings.Contains(c.Request().Header.Get("Content-Type"), "application/json") {
		// 415: Invalid `Content-Type`
		return c.JSONPretty(http.StatusUnsupportedMediaType, map[string]string{"message": "unsupported media type"}, "	")
	}

	// Check token
	u := c.Get("user").(*jwtGo.Token)
	userId, err := jwt.CheckToken(*flags.Get().JwtIssuer, u)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnauthorized, map[string]string{"message": err.Error()}, "	")
	}

	// id
	idStr := c.Param("id")

	// string -> uint64
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		// 404: Not found
		return echo.ErrNotFound
	}

	// Bind request body
	patch := new(document.PatchBody)
	if err = c.Bind(patch); err != nil {
		// 400: Bad request
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
	}

	// Validate request body
	if err = c.Validate(patch); err != nil {
		// 422: Unprocessable entity
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnprocessableEntity, map[string]string{"message": err.Error()}, "	")
	}

	// Check project id
	if patch.ProjectId != nil {
		status, err := utils.HttpGet(fmt.Sprintf("%s/%d", *flags.Get().ServiceUrlProjects, *patch.ProjectId), &u.Raw)
		if err != nil {
			// 500: Internal server error
			c.Logger().Error(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		if status != http.StatusOK {
			// 400: Bad request
			c.Logger().Debugf("project id: %d does not exist", *patch.ProjectId)
			return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("project id: %d does not exist", *patch.ProjectId)}, "	")
		}
	}

	p, notFound, err := document.Patch(userId, id, *patch)
	if err != nil {
		// 500: Internal server error
		c.Logger().Error(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notFound {
		// 404: Not found
		c.Logger().Debug("project not found")
		return echo.ErrNotFound
	}

	// 200: Success
	return c.JSONPretty(http.StatusOK, p, "	")
}
