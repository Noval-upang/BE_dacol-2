package helpers

import (
	"fmt"
	"os"
	"time"

	"github.com/Hy-Iam-Noval/dacol-2/src/DB"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
)

const (
	UserKey = "user"
	Breaker = "Breaker"
)

type (
	JWTClaim struct {
		jwt.RegisteredClaims
		DB.User
	}
	JwtMethod = *jwt.SigningMethodHMAC
)

var JwtToken string = os.Getenv("JWT_TOKEN")

func parse(signToken string) (*jwt.Token, error) {
	spiltToken := signToken[len(Breaker)+1:]
	return jwt.ParseWithClaims(
		spiltToken,
		&JWTClaim{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(JwtMethod); !ok {
				return nil, fmt.Errorf("method must be %s", t.Method.Alg())
			}
			return []byte(JwtToken), nil
		},
	)
}

func invalidToken(signToken string) bool {
	token, _ := parse(signToken)
	if token == nil {
		return false
	}

	_, ok := token.Claims.(*JWTClaim)
	return ok && token.Valid
}

func GenerateToken(user DB.User) string {
	t := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		JWTClaim{
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
			user,
		},
	)

	s, err := t.SignedString([]byte(JwtToken))
	if err != nil {
		panic(err)
	}

	return s
}

func AuthWare(c Ctx) error {
	signToken := c.Get(UserKey)
	if signToken == "" || !invalidToken(signToken) {
		return c.Status(UnAuth).Send(nil)
	}

	user := ParseTokenUser(signToken)
	DB.
		Select("user", "*",
			fmt.Sprintf("id = %d", user.Id),
			fmt.Sprintf(`email = "%s"`, user.Email),
			fmt.Sprintf(`password = "%s"`, user.Password)).
		Single().
		Scan(&user.Id, &user.Email, &user.Password)
	if (user == DB.User{}) {
		return c.Status(Invalid).Send(nil)
	}

	return c.Status(200).Next()
}

func ParseTokenUser(signToken string) DB.User {
	token, _ := parse(signToken)
	claim := token.Claims.(*JWTClaim)
	return claim.User
}
