package middleware

type MiddlewareProvider struct {
	JWTAuthMiddleware
}

func ProvideMiddlewareProvider(jwtAuthMiddleware JWTAuthMiddleware) *MiddlewareProvider {
	return &MiddlewareProvider{
		jwtAuthMiddleware,
	}
}
