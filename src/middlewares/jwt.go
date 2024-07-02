package middlewares

import (
	"strings"
	"webapp-template/src/errors"
	"webapp-template/src/utils"

	"github.com/gin-gonic/gin"
)

func VerifyJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		splits := strings.Split(authorization, " ")
		if len(splits) != 2 || splits[0] != "Bearer" {
			utils.HandleError(c, errors.New(errors.ErrUnauthorized))
			c.Abort()
			return
		}

		// token, err := jwt.Parse(splits[1], func(t *jwt.Token) (interface{}, error) {
		// 	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		// 		return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		// 	}
		// 	return nil, nil
		// })

		// if err != nil {
		// 	logger.Context(c.Request.Context()).Errorf("authorize got error: %v", err)
		// 	utils.HandleError(c, errors.New(errors.ErrUnauthorized, "token had expired"))
		// 	c.Abort()
		// 	return
		// }

		// claims, ok := token.Claims.(jwt.MapClaims)
		// if !ok && !token.Valid {
		// 	utils.HandleError(c, errors.New(errors.ErrUnauthorized))
		// 	c.Abort()
		// 	return
		// }

		// data, _ := json.Marshal(claims)
		// rid := dtos.RemembranceID{}
		// if err := json.Unmarshal(data, &rid); err != nil {
		// 	utils.HandleError(c, errors.New(errors.ErrInternalServer))
		// 	c.Abort()
		// 	return
		// }

		// c.Request.WithContext(utils.SetUserID(c.Request.Context(), rid.UserID))
		c.Next()
	}
}
