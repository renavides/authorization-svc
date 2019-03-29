package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io"
	"log"
	"net/http"
	"strings"
)

func Home(w http.ResponseWriter, r *http.Request){
	val := r.URL.Path[len("/lm/home/"):]
	fmt.Println(val)
	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(string(auth[1]), claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretCon.JwtSigningKey), nil
	})
	// ... error handling
	if err != nil{
		log.Println(err)
		respondWithError(w,http.StatusForbidden, "403 Forbidden")
		return
	}
	if !token.Valid {
		respondWithError(w,http.StatusForbidden, "403 Forbidden")
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
		r.Host = "localhost:31062"
		r.URL.Host = "localhost:31062"
		r.URL.Scheme = "http"
		r.RequestURI = ""
		r.URL.Path = "/integrator/v1/page/homepage/"
		resp, err := http.DefaultClient.Do(r)
		if err != nil {
			log.Println(err)
			respondWithError(w,http.StatusInternalServerError, "500 Internal Server Error")
			return
		}
		io.Copy(w,resp.Body)
		w.WriteHeader(http.StatusOK)

	}


	//respondWithJson(w, http.StatusOK, "200 Ok")

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