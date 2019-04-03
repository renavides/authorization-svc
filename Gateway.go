package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
)

type Gateway struct {
	Path string
    Url string
    Context string
}


func (g *Gateway)createGateway(w http.ResponseWriter, r *http.Request){
	//gate
	u, err := url.Parse(g.Url)
	if err != nil {
		log.Println(err)
	}
	log.Print(u.Scheme)
	log.Println(u.Host)
	log.Println(u.Path)
	r.Host = u.Host
	r.URL.Host = u.Host
	r.URL.Scheme = u.Scheme
	r.RequestURI = ""
	r.URL.Path = u.Path
	s, _, _ := net.SplitHostPort(r.RemoteAddr)
	r.Header.Set("X-Forwarded-For", s)
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Println(err)
		respondWithError(w,http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	for key, values := range resp.Header {
		for _ , value := range values{
			w.Header().Add( string(key), string(value))
		}
	}

	io.Copy(w,resp.Body)
	w.WriteHeader(http.StatusOK)
}