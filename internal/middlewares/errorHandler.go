package middlewares

import (
	"errors"
	"fmt"
	"github.com/eduardolat/goeasyi18n"
	"github.com/go-playground/validator/v10"
	"github.com/gobeam/stringy"
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	HttpErrorHandler struct {
		statusCodes map[error]int
	}
)

func NewHttpErrorHandler(errorStatusCodeMaps map[error]int) *HttpErrorHandler {
	return &HttpErrorHandler{
		statusCodes: errorStatusCodeMaps,
	}
}

func (eh *HttpErrorHandler) getStatusCode(err error) int {
	for key, value := range eh.statusCodes {
		if errors.Is(err, key) {
			return value
		}
	}

	return http.StatusInternalServerError
}

func unwrapRecursive(err error) error {
	var originalErr = err

	for originalErr != nil {
		var internalErr = errors.Unwrap(originalErr)

		if internalErr == nil {
			break
		}

		originalErr = internalErr
	}

	return originalErr
}

func (eh *HttpErrorHandler) Handler(err error, c echo.Context) {
	var he *echo.HTTPError

	switch errs := err.(type) {
	case *echo.HTTPError:
		if errs.Internal != nil {
			var herr *echo.HTTPError
			if errors.As(errs.Internal, &herr) {
				he = herr
			}
		}
		he = &echo.HTTPError{
			Code:    eh.getStatusCode(err),
			Message: unwrapRecursive(err).Error(),
		}
	case validator.ValidationErrors:
		errorsMessages := map[string]interface{}{}
		i18n := c.Get("i18n").(*goeasyi18n.I18n)
		language := c.Get("lang").(string)
		for _, e := range errs {
			errorsMessages[stringy.New(e.Field()).SnakeCase().ToLower()] = i18n.T(language,
				fmt.Sprintf("field_%s", e.Tag()),
				goeasyi18n.Options{
					Data: map[string]string{
						"valData": e.Param(),
					},
				})
		}
		he = &echo.HTTPError{
			Code:    eh.getStatusCode(err),
			Message: map[string]interface{}{"errors": errorsMessages},
		}
	default:
		he = &echo.HTTPError{
			Code:    eh.getStatusCode(err),
			Message: unwrapRecursive(err).Error(),
		}
	}

	if _, ok := he.Message.(string); ok {
		he.Message = map[string]interface{}{"errors": err.Error()}
	}

	// Send response
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(he.Code)
		} else {
			err = c.JSON(he.Code, he.Message)
		}
		if err != nil {
			c.Logger().Error(err)
		}
	}

}
