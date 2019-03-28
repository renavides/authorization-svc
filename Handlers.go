package main

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"strings"
)

func Home(w http.ResponseWriter, r *http.Request){
	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(string(auth[1]), claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretCon.JwtSigningKey), nil
	})
	// ... error handling
	if err != nil{
		log.Println(err)
		respondWithError(w,http.StatusUnauthorized, "401 Unauthorized")
		return
	}
	if !token.Valid {
		respondWithError(w,http.StatusUnauthorized, "401 Unauthorized")
		return
	}

	ok, err := validateTid(claims["tid"])
	if err != nil{
		log.Println(err)
		respondWithError(w,http.StatusBadGateway, "400 BadRequest")
		return
	}
	if ok == true {
		//gate
	}


	respondWithJson(w, http.StatusOK, "200 Ok")

}

func validateTid(interface{}) (bool, error){
	return true, nil
}


func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}