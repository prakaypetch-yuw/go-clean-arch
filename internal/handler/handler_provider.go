package handler

type HandlerProvider struct {
	UserHandler UserHandler
}

func ProvideHandlerProvider(userHandler UserHandler) *HandlerProvider {
	return &HandlerProvider{
		userHandler,
	}
}
