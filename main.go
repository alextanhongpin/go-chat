package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/alextanhongpin/go-chat/chat"
	"github.com/alextanhongpin/go-chat/controller"
	"github.com/alextanhongpin/go-chat/database"
	"github.com/alextanhongpin/go-chat/entity"
	"github.com/alextanhongpin/go-chat/pkg/token"
	"github.com/alextanhongpin/go-chat/service"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

func main() {
	var (
		dbUser    = os.Getenv("DB_USER")
		dbPass    = os.Getenv("DB_PASS")
		dbName    = os.Getenv("DB_NAME")
		jwtSecret = os.Getenv("JWT_SECRET")
		jwtIssuer = os.Getenv("JWT_ISSUER")
		port      = ":4000"
	)

	db, err := database.New(dbUser, dbPass, dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	signer := token.New(token.SignerOptions{
		Now: func() time.Time {
			return time.Now().UTC()
		},
		Issuer: jwtIssuer,
		TTL:    5 * time.Minute,
		Secret: []byte(jwtSecret),
	})

	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any

	c := chat.New(db, NewRedis(), logger)
	defer c.Close()

	ctl := controller.New()
	authorized := authMiddleware(signer)

	getRoomsService := service.NewGetRoomsService(db)
	getConversationsService := service.NewGetConversationsService(db)
	postAuthorizeService := service.NewAuthorizeService(db)
	postLoginService := service.NewLoginService(db, signer)
	postRegisterService := service.NewRegisterService(db, signer)

	mux := http.NewServeMux()
	// Serve public files.
	mux.Handle("/", http.FileServer(http.Dir("./public")))

	mux.HandleFunc("/ws", c.ServeWS(signer, db))
	mux.HandleFunc("/auth", authorized(ctl.PostAuthorize(postAuthorizeService)))
	mux.HandleFunc("/rooms", authorized(ctl.GetRooms(getRoomsService)))
	mux.HandleFunc("/conversations/", authorized(ctl.GetConversations(getConversationsService)))
	mux.HandleFunc("/register", ctl.PostRegister(postRegisterService))
	mux.HandleFunc("/login", ctl.PostLogin(postLoginService))

	srv := &http.Server{
		Addr:         port,
		Handler:      mux,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		log.Printf("listening to port *%s. press ctrl + c to cancel.\n", port)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("gracefully shut down application")
}

// NewRedis returns a new redis.Client
func NewRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return client
}

type middleware func(http.HandlerFunc) http.HandlerFunc

func authMiddleware(signer token.Signer) middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if values := strings.Split(auth, " "); len(values) != 2 {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			} else {
				bearer, token := values[0], values[1]
				if bearer != "Bearer" {
					http.Error(w, "invalid bearer type", http.StatusUnauthorized)
					return
				}
				userID, err := signer.Verify(token)
				if err != nil {
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}
				ctx := r.Context()
				ctx = context.WithValue(ctx, entity.ContextKeyUserID, userID)
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		}
	}
}
