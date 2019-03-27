package main

import (
	"encoding/json"
	"fmt"
	"github.com/dimiro1/health"
	"github.com/dimiro1/health/url"
	"github.com/gorilla/mux"
	"github.com/renavides/vault/client"
	"github.com/renavides/vault/config"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)
func main() {
	log.Println("Starting server initialization")
    //init configurations
	log.Println("Starting config initialization")
	v := intConfig()
	log.Println("Starting vault initialization")
	/*err := v.Initialize()
	if err != nil {
		log.Fatal(err)
	}
*/
	//Router
	r := mux.NewRouter()
	//API Routes
	r.HandleFunc("/lifemiles/home/", home).Methods("GET")
	//Health Check Routes
	h := health.NewHandler()
	h.AddChecker("Vault", url.NewChecker(fmt.Sprintf("%s://%s:%s/v1/sys/health?perfstandbyok=true", v.Scheme, v.Host, v.Port)))
	r.Path("/health").Handler(h).Methods("GET")
	//Server config - http
	go func() {
		log.Println(fmt.Sprintf("Server is now accepting http requests on port %v", 8080))
		if err := http.ListenAndServe(fmt.Sprintf(":%v", 8080), r); err != nil {
			log.Fatal(err)
		}
	}()
	//Catch SIGINT AND SIGTERM to gracefully
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	sig := <-gracefulStop
	fmt.Printf("caught sig: %+v", sig)
}
/*
func Login(w http.ResponseWriter, r *http.Request) {
	var user user.Users
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	respondWithJson(w,http.StatusOK, "ok")
	fmt.Println(user)
}
*/
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func intConfig() *client.Vault {

	var config = config.Config{}
	config.Read()
	//Server params
	var credential = client.Credential{
		Token:          config.Vault.Credential.Token,
		RoleID:         config.Vault.Credential.RoleID,
		SecretID:       config.Vault.Credential.SecretID,
		ServiceAccount: config.Vault.Credential.ServiceAccount,
	}
	var vault = client.Vault{
		Host:           config.Vault.Host,
		Port:           config.Vault.Port,
		Scheme:         config.Vault.Scheme,
		Authentication: config.Vault.Authentication,
		Role:           config.Vault.Role,
		Mount:          config.Vault.Mount,
		Credential:     credential,
	}
	return &vault
}

func home(w http.ResponseWriter, r *http.Request){
	auth := r.Header.Get("Authorization")

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(auth, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})
	// ... error handling
	if err != nil{
		log.Println(err)
	}
	fmt.Println(token)
	// do something with decoded claims
	for key, val := range claims {
		fmt.Printf("Key: %v, value: %v\n", key, val)
	}
	respondWithJson(w,http.StatusOK, "ok")
}