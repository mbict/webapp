package webapp

type Option func(WebApp)

// WithErrorHandler will overwrite the default handling of error handling with this instance
func WithErrorHandler(errorHandlers ...ErrorHandler) Option {

	if len(errorHandlers) == 0 {
		panic("at least one error handler is required")
	}

	eh := errorHandlers[len(errorHandlers)-1]
	for i := len(errorHandlers) - 2; i >= 0; i-- {
		current := errorHandlers[i]
		next := eh
		eh = func(c Context, err error) error {
			err = current(c, err)
			if err != nil {
				return next(c, err)
			}
			return err
		}
	}

	return func(app WebApp) {
		app.(*webapp).errorHandler = eh
	}
}

// WithErrorHandlerFallback will set the error handler to first try to handle the error with the provided error handler
// If the error handler could nto handle the error and returns a non nil error, the default error handler will act as
// a fallback error handler
func WithErrorHandlerFallback(errorHandlers ...ErrorHandler) Option {
	return WithErrorHandler(append(errorHandlers, DefaultErrorHandler)...)
}

func WithRouter(router Router) Option {
	return func(app WebApp) {
		app.(*webapp).router = router
	}
}

func WithBinder(binder Binder) Option {
	return func(app WebApp) {
		app.(*webapp).binder = binder
	}
}

func WithJsonEncoder(encoder JSONEncoding) Option {
	return func(app WebApp) {
		app.(*webapp).jsonEncoder = encoder
	}
}

func WithRenderer(renderer Renderer) Option {
	return func(app WebApp) {
		app.(*webapp).renderer = renderer
	}
}

func WithValidator(validator Validator) Option {
	return func(app WebApp) {
		app.(*webapp).validator = validator
	}
}
