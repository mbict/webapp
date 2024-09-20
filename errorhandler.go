package webapp

import (
	"errors"
	"net/http"
)

// ErrorHandler is a centralized error handler.
type ErrorHandler func(Context, error) error

var DefaultErrorHandler = func(c Context, err error) error {

	if IsBindError(err) {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
			"error":   errors.Unwrap(err),
			"type":    "bind",
		})
	}

	if IsValidationError(err) {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
			"error":   err,
			"type":    "validation",
		})
	}

	he, ok := err.(*HTTPError)
	if ok {
		if he.Internal != nil {
			if herr, ok := he.Internal.(*HTTPError); ok {
				he = herr
			}
		}
	} else {
		he = &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		}
	}

	message := he.Message
	if m, ok := he.Message.(string); ok {
		//if e.Debug {
		//	message = Map{"message": m, "error": err.Error()}
		//} else {
		message = Map{"message": m}
		//}
	}

	// Send response
	if c.Request().Method == http.MethodHead {
		c.Response().WriteHeader(he.Code)
	} else {
		err = c.JSON(he.Code, message)
	}
	//if err != nil {
	//e.Logger().Error(err)
	//}

	return err
}
