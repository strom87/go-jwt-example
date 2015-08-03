package main

/*
 * Important if you don't know how to create the public and private rsa files.
 * In the terminal window write the following two lines to create the files.
 *
 * openssl genrsa -out example.rsa 1024
 * openssl rsa -in example.rsa -pubout > example.rsa.pub
 *
 * Now you should have the two files in the current folder you are in.
 * FYI, in the first line you can change the number 1024 to what ever you prefer, ex 2048.
 */

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"time"
)

// The paths to the rsa key files
const (
	publicKeyFilePath  = "keys/example.rsa.pub"
	privateKeyFilePath = "keys/example.rsa"
)

// Holds the rsa files keys that is used to sign the token
var (
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
)

// Creates the token using RS256 and adds the claims
func TokenHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token := jwt.New(jwt.SigningMethodRS256)

	token.Claims["name"] = "user name example"
	token.Claims["email"] = "useremail@example.com"
	token.Claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	token.Claims["iat"] = time.Now().Unix()

	tokenString, err := token.SignedString(privateKey)
	LogError(err)

	fmt.Fprintln(w, tokenString)
}

// Middleware that is used to check that the token is correct
// before giving access to the protected routes
func AuthMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("Unexpected signing method: " + token.Header["alg"].(string))
		}
		return publicKey, nil
	})

	if err != nil || !token.Valid {
		LogError(err)
		fmt.Fprintln(w, "Not authenticated, route protected")
		return
	}

	next(w, r)
}

// Api path is the standard route but it is protected so
// it needs a correct token to access it
func ApiHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Fprintln(w, "Congrats, you have access to the protected route")
}

// Reads the public and private key file that is used to sign the token
// and adds them to the privateKey and publicKey variable
func Init() {
	publicKeyBytes, err := ioutil.ReadFile(publicKeyFilePath)
	LogError(err)

	privateKeyBytes, err := ioutil.ReadFile(privateKeyFilePath)
	LogError(err)

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	LogError(err)

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	LogError(err)
}

// Only used to make it easy to write errors to the terminal window
func LogError(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

func main() {
	Init()
	router := mux.NewRouter()
	n := negroni.Classic()

	router.Handle("/token", negroni.New(negroni.HandlerFunc(TokenHandler)))

	// Add the AuthMiddleware to the request so it first checks if
	// the token is valid before giving access to the api route
	router.Handle("/api", negroni.New(negroni.HandlerFunc(AuthMiddleware), negroni.HandlerFunc(ApiHandler)))

	n.UseHandler(router)
	http.ListenAndServe(":1337", n)
}
