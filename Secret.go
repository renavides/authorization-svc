package main

type Secret struct {
	JwtSigningKey string `json:"security.jwt.signing_key"`
	Password string `json:"spring.datasource.password"`
	Username string `json:"spring.datasource.username"`
}
