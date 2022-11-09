package main

//******** jwt基本使用：默认claim********
//******** jwt自定义claim********
//******** gin框架中使用：claim********

//******** jwt基本使用：默认claim***********************************************
//import (
//	"fmt"
//	"github.com/golang-jwt/jwt/v4"
//	"time"
//)
//
//// 用于签名的字符串
//var mySigningKey = []byte("Accept for wangkun")
//
//// GenRegisteredClaims 使用默认声明创建jwt
//func GenRegisteredClaims() (string, error) {
//	// 创建 Claims
//	claims := &jwt.RegisteredClaims{
//		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // 过期时间
//		Issuer:    "wangkun",                                          // 签发人
//	}
//	// 生成token对象
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	// 生成签名字符串
//	return token.SignedString(mySigningKey)
//}
//
//// ParseRegisteredClaims 解析jwt
//func ValidateRegisteredClaims(tokenString string) bool {
//	// 解析token
//	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//		return mySigningKey, nil
//	})
//	if err != nil { // 解析token失败
//		return false
//	}
//	return token.Valid
//}
//func main() {
//	str, _ := GenRegisteredClaims()
//	fmt.Println(ValidateRegisteredClaims(str))
//}

//******** jwt自定义claim************************************************

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 这里额外记录username、phone两字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
//type MyClaims struct {
//	Username string `json:"username"`
//	Phone    string `json:"phone"`
//	jwt.RegisteredClaims
//}
//
////密钥
//var MySecret = []byte("密钥")
//
//func main() {
//	tokenStr, _ := GenToken("张三")
//	fmt.Println("token:", tokenStr)
//	//tokenStr := ""
//	claim, _ := ParseToken(tokenStr)
//	fmt.Printf("解析后：%#v\n", claim)
//	tokenStr2, _ := RefreshToken(tokenStr)
//	fmt.Println("refToken:", tokenStr2)
//}
//
////生成 Token
//func GenToken(username string) (string, error) {
//	// 创建一个我们自己的声明
//	c := MyClaims{
//		username,
//		"10000000086",
//		jwt.RegisteredClaims{
//			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 1)), // 过期时间
//			Issuer:    "test",                                              // 签发人
//		},
//	}
//	// 使用指定的签名方法创建签名对象
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
//	// 使用指定的secret签名并获得完整的编码后的字符串token
//	return token.SignedString(MySecret)
//}
//
//// 解析 Token
//func ParseToken(tokenStr string) (*MyClaims, error) {
//	token, err := jwt.ParseWithClaims(tokenStr, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
//		return MySecret, nil
//	})
//	if err != nil {
//		fmt.Println(" token parse err:", err)
//		return nil, err
//	}
//	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
//		return claims, nil
//	}
//	return nil, errors.New("invalid token")
//}
//
//// 刷新 Token
//func RefreshToken(tokenStr string) (string, error) {
//	jwt.TimeFunc = func() time.Time {
//		return time.Unix(0, 0)
//	}
//
//	token, err := jwt.ParseWithClaims(tokenStr, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
//		return MySecret, nil
//	})
//	if err != nil {
//		return "", err
//	}
//	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
//		jwt.TimeFunc = time.Now
//		claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * 10))
//		return GenToken(claims.Username)
//	}
//	return "", errors.New("Couldn't handle this token")
//}

//******** gin中使用jwt***********************************************
