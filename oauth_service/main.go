package main

import (
	"context"
	"log"
	"net/http"
	"oauth2_provider/internal/encoding"
	"oauth2_provider/internal/server"
	"oauth2_provider/internal/service"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	connStr, ok := os.LookupEnv("OAUTH_DB_CONNECTION_STRING")
	if !ok {
		log.Fatal("DB_CONNECTION_STRING is not set")
	}

	conn, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	clientService := service.ClientServiceImpl{Db: conn}
	sessionService := service.SessionServiceImpl{Db: conn}
	userService := service.UserServiceImpl{Db: conn}
	keyStorage, err := encoding.NewKeyStorage(filepath.Join(".", "storage"), "private.pem", "public.pem")
	if err != nil {
		log.Fatal(err)
	}
	jwtProvider := service.NewJwtProvider(*keyStorage)
	authProvider := service.AuthProviderImpl{UserService: &userService,
		ClientService:  &clientService,
		JwtProvider:    jwtProvider,
		SessionService: &sessionService}

	authHandler := server.OauthServerImpl{Ap: &authProvider}
	clientHandler := server.ClientServer{ClientService: &clientService}
	userServer := server.UserServer{UserService: &userService}

	mux := http.NewServeMux()

	mux.HandleFunc("/authorize", authHandler.HandleAuthorize)
	mux.HandleFunc("/token", authHandler.HandleToken)

	mux.HandleFunc("/client", clientHandler.SaveHandler)
	mux.HandleFunc("/user", userServer.SaveHandler)
	mux.HandleFunc("/public-key", service.GetPublicKeyHandler(keyStorage))

	http.ListenAndServe(":5000", mux)

}
