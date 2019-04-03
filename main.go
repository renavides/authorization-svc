package main

import (
	"encoding/json"
	"fmt"
	"github.com/dimiro1/health"
	"github.com/dimiro1/health/url"
	"github.com/gorilla/mux"
	"github.com/renavides/vault/client"
	"github.com/renavides/vault/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var SecretCon Secret
var VaultClient *client.Vault

func main() {
	log.Println("Starting server initialization")
    //init configurations
	log.Println("Starting config initialization")
	serverPort, err , VaultClient := intConfig()
	if err != nil {
		log.Fatal(err)
	}

	//Router
	r := mux.NewRouter()
	//API Routes
	r.HandleFunc("/lm/home/", Home).Methods("POST")
	//Health Check Routes
	h := health.NewHandler()
	h.AddChecker("Vault", url.NewChecker(fmt.Sprintf("%s://%s:%s/v1/sys/health?perfstandbyok=true", VaultClient.Scheme, VaultClient.Host, VaultClient.Port)))
	r.Path("/health").Handler(h).Methods("GET")

	//Server config - http
	go func() {
		log.Println(fmt.Sprintf("Server is now accepting http requests on port %v", serverPort))
		if err := http.ListenAndServe(fmt.Sprintf(":%v", serverPort), r); err != nil {
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



func intConfig()  (string , error, *client.Vault){

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
	log.Println("Starting vault initialization")
	err := vault.Initialize()
	if err != nil {
		log.Fatal(err)
	}

	secret, err := vault.GetSecret("secret-v1/lm-auth-svc/dev")

	if err != nil {
		return string(config.Server.Port), err , &vault
	}
	if secret.Data != nil {
		j, _ := json.Marshal(secret.Data)
		json.Unmarshal(j, &SecretCon)
	}
	return string(config.Server.Port), nil , &vault
}

