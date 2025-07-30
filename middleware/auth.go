package middleware

import (
	"NotificationManagement/config"
	"NotificationManagement/utils/errutil"

	"github.com/Nerzal/gocloak/v13"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type CustomContext struct {
	echo.Context
	Roles *[]string
}

// KeycloakMiddleware creates a middleware to validate JWT and extract roles
func KeycloakMiddleware() echo.MiddlewareFunc {
	keycloakCfg := config.Keycloak()
	client := gocloak.NewClient(keycloakCfg.ServerURL)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return errutil.NewAppError(errutil.ErrMissingAuthHeader, errutil.ErrNoAuthHeader)
			}

			token := ""
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				token = authHeader[7:]
			} else {
				return errutil.NewAppError(errutil.ErrInvalidAuthFormat, errutil.ErrInvalidHeaderFormat)
			}

			ctx := c.Request().Context()
			result, err := client.RetrospectToken(ctx, token, keycloakCfg.ClientID, keycloakCfg.ClientSecret, keycloakCfg.Realm)
			if err != nil || result == nil || !*result.Active {
				return errutil.NewAppError(errutil.ErrInvalidToken, err)
			}

			// Parse the JWT token to extract roles
			parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				// Verify the token signing method
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, errutil.NewAppError(errutil.ErrInvalidSigningMethod, errutil.ErrInvalidTokenSignature)
				}

				// Get the public key from Keycloak
				cert, err := client.GetCerts(ctx, keycloakCfg.Realm)
				if err != nil {
					return nil, errutil.NewAppError(errutil.ErrCertificateRetrieval, err)
				}

				if cert == nil || len(*cert.Keys) == 0 {
					return nil, errutil.NewAppError(errutil.ErrNoCertificates, errutil.ErrNoCertificateFound)
				}

				// Use the first key in the certificate
				certKey := (*cert.Keys)[0]
				if certKey.X5c == nil || len(*certKey.X5c) == 0 {
					return nil, errutil.NewAppError(errutil.ErrNoCertificates, errutil.ErrNoCertificateFound)
				}

				// The certificate is in base64-encoded DER format, wrapped in PEM
				certPEM := "-----BEGIN CERTIFICATE-----\n" + (*certKey.X5c)[0] + "\n-----END CERTIFICATE-----"
				key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(certPEM))
				if err != nil {
					return nil, errutil.NewAppError(errutil.ErrCertificateRetrieval, err)
				}

				return key, nil
			})
			if err != nil {
				return errutil.NewAppError(errutil.ErrInvalidToken, err)
			}

			if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
				var roles []string
				if realmAccess, ok := claims["realm_access"].(map[string]interface{}); ok {
					if rolesInterface, ok := realmAccess["roles"].([]interface{}); ok {
						for _, role := range rolesInterface {
							if roleStr, ok := role.(string); ok {
								roles = append(roles, roleStr)
							}
						}
					}
				}

				// Create custom context with roles
				cc := &CustomContext{
					Context: c,
					Roles:   &roles,
				}

				return next(cc)
			}

			return next(c)
		}
	}
}

// RequireRoles middleware to check if user has required roles
func RequireRoles(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc, ok := c.(*CustomContext)
			if !ok {
				return errutil.NewAppError(errutil.ErrNoRoleInfo, errutil.ErrNoRoleInformation)
			}

			// Check if user has any of the required roles
			for _, requiredRole := range roles {
				for _, userRole := range *cc.Roles {
					if requiredRole == userRole {
						return next(c)
					}
				}
			}

			return errutil.NewAppError(errutil.ErrInsufficientPrivileges, errutil.ErrInsufficientRoles)
		}
	}
}
