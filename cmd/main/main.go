package main

import (
	"github.com/lactobasilusprotectus/go-template/docs"
	authDelivery "github.com/lactobasilusprotectus/go-template/pkg/auth/delivery"
	authUsecase "github.com/lactobasilusprotectus/go-template/pkg/auth/usecase"
	"github.com/lactobasilusprotectus/go-template/pkg/common/config"
	commonTime "github.com/lactobasilusprotectus/go-template/pkg/common/time"
	"github.com/lactobasilusprotectus/go-template/pkg/domain"
	rootDelivery "github.com/lactobasilusprotectus/go-template/pkg/root/delivery"
	userRepository "github.com/lactobasilusprotectus/go-template/pkg/user/repository"
	"github.com/lactobasilusprotectus/go-template/pkg/util/db"
	_ "github.com/lactobasilusprotectus/go-template/pkg/util/http"
	httputil "github.com/lactobasilusprotectus/go-template/pkg/util/http"
	"github.com/lactobasilusprotectus/go-template/pkg/util/jwt"
	"github.com/lactobasilusprotectus/go-template/pkg/util/redis"
	"log"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

// @version					1.0
// @termsOfService				http://swagger.io/terms/
// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
// @schemes					http
// @BasePath					/
// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
// @securityDefinitions.apikey	JWT
// @in							header
// @name						Authorization
// @Security					JWT
func main() {
	// Get env
	env := os.Getenv(config.ENV)

	if env == "" {
		env = config.LOC
	}

	// Read env file
	cfg, err := config.Read(config.GetFilePath(env))

	//swagger info
	docs.SwaggerInfo.Title = cfg.Title
	docs.SwaggerInfo.Description = cfg.Description
	docs.SwaggerInfo.Host = cfg.URL

	if err != nil {
		return
	}

	// Init utils: http server, db connection, etc.
	utils := initUtils(cfg)

	// Init repository and use case layer
	_, uc, err := initRepoAndUseCases(utils, cfg)
	if err != nil {
		log.Fatalln("initRepoAndUseCases err:", err)
	}

	// Init delivery layer (HTTP)
	httpHandler := initHttpHandler(utils, uc, env)

	// Start serving
	registerHttpHandler(utils.HttpServer, httpHandler)
	utils.HttpServer.Run(env)

	// =======================================================

	// Catching signals

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-c

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	utils.HttpServer.Stop()

	log.Println("shutting down...")
	os.Exit(0)

}

// initUtils initialises utility for the app
// this includes httputil server, db connection, and any other dependent tools
func initUtils(cfg config.Config) AppUtil {
	// database connection
	dbConn, err := db.NewDatabaseConnection(cfg.Database)

	//migration base on AppModels
	registerModels(dbConn, AppModels{})

	if err != nil {
		log.Fatalln(err)
	}

	// redis as cache
	redisClient := redis.NewRedisClient(cfg.Redis)

	// time module
	timeModule := commonTime.New()

	// JWT implementation
	jwtModule := jwt.New(timeModule)

	return AppUtil{
		HttpServer:   httputil.NewServer(cfg.Http),
		DbConnection: dbConn,
		Redis:        redisClient,
		Jwt:          jwtModule,
		Time:         timeModule,
	}
}

// initHttpHandler initialises http handler for the app
func initHttpHandler(ut AppUtil, uc AppUseCase, env string) AppHttpHandler {
	rootHandler := rootDelivery.NewRootHandler(env)

	return AppHttpHandler{
		RootHttpHandler: rootHandler,
		AuthHttpHandler: authDelivery.NewAuthHttpHandler(uc.AuthUseCase, uc.AuthUseCase),
	}
}

// initRepoAndUseCases initialises repo and use case layer
func initRepoAndUseCases(util AppUtil, cfg config.Config) (repo AppRepo, uc AppUseCase, err error) {
	repo.User = userRepository.NewUserRepository(util.DbConnection, util.Time)

	//usecase
	uc.AuthUseCase = authUsecase.NewAuthUseCase(repo.User, util.Jwt, util.Redis, util.Time, cfg)

	return repo, uc, nil
}

// registerHttpHandler registers our handlers to the http server.
// reflect docs: https://golang.org/pkg/reflect/.
// The purpose of this function is to register HTTP request handlers to an httputil.Server object based on the fields of the handlers object.
func registerHttpHandler(srv *httputil.Server, handlers AppHttpHandler) {
	h := reflect.ValueOf(handlers)

	for i := 0; i < h.NumField(); i++ {
		srv.RegisterHandler(h.Field(i).Interface().(httputil.RouterHandler))
	}
}

// registerModels registers our models to the database.
// The purpose of this function is to automatically create tables or migrate existing ones in the database based on the fields of the models object.
func registerModels(dbConn *db.DatabaseConnection, models AppModels) {
	m := reflect.ValueOf(models)

	for i := 0; i < m.NumField(); i++ {

		err := dbConn.Master.AutoMigrate(m.Field(i).Interface())

		if err != nil {
			log.Fatalf("AutoMigrate err: %v", err)
		}
	}
}

//================ TYPES =================

// AppUtil wraps utility layer with the app, includes delivery and database
type AppUtil struct {
	HttpServer   *httputil.Server
	DbConnection *db.DatabaseConnection
	Redis        redis.Interface
	Jwt          *jwt.JwtModule
	Time         *commonTime.Time
}

// AppHttpHandler wraps HTTP handlers exposed by the app as a delivery layer
type AppHttpHandler struct {
	RootHttpHandler *rootDelivery.RootHandler
	AuthHttpHandler *authDelivery.AuthHttpHandler
}

// AppUseCase wraps use case layer within the app
type AppUseCase struct {
	AuthUseCase *authUsecase.AuthUseCase
}

// AppRepo wraps repository layer within the app
type AppRepo struct {
	User *userRepository.UserRepository
}

// AppModels wraps domain models within the app
type AppModels struct {
	User *domain.User
}
