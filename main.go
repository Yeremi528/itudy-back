package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	assignments "github.com/Yeremi528/itudy-back/assigments"
	assigmentsdb "github.com/Yeremi528/itudy-back/assigments/repository/asigments"
	"github.com/Yeremi528/itudy-back/courses"
	"github.com/Yeremi528/itudy-back/courses/repository/coursesdb"
	"github.com/Yeremi528/itudy-back/email"
	"github.com/Yeremi528/itudy-back/exam"
	"github.com/Yeremi528/itudy-back/exam/repository/examdb"
	"github.com/Yeremi528/itudy-back/kit/logger"
	db "github.com/Yeremi528/itudy-back/kit/mongo"
	"github.com/Yeremi528/itudy-back/kit/secretmanager"
	"github.com/Yeremi528/itudy-back/kit/tracer"
	"github.com/Yeremi528/itudy-back/learning"
	"github.com/Yeremi528/itudy-back/learning/repository/learningdb"
	"github.com/Yeremi528/itudy-back/movements/repository/movementsdb"
	"github.com/Yeremi528/itudy-back/oauth"
	"github.com/Yeremi528/itudy-back/payin"
	"github.com/Yeremi528/itudy-back/payin/repository/payindb"
	"github.com/Yeremi528/itudy-back/user"
	"github.com/Yeremi528/itudy-back/user/repository/userdb"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/resend/resend-go/v2"
)

var build = "dev"

func main() {
	var (
		ctx = context.Background()

		writeTo = os.Stdout
		service = ""
		level   = logger.LevelInfo
	)

	log := logger.New(writeTo, level, service, traceFunc)
	log.Info(ctx, "Startup - Service Details", "logLevel", log.GetLevel().ToString(), "build", build, "cores", runtime.GOMAXPROCS(0))

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "service error, shutting down", "errorDetails", err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {

	sm, err := secretmanager.New(ctx, secretmanager.Config{
		ProjectID: "itudy-485221",
	})
	if err != nil {
		return err
	}

	cfg, err := loadConfig(ctx, sm)
	if err != nil {
		return err
	}

	// -----------------------------------------------------------------------
	// init DB
	db, err := db.ConnectMongo(cfg.Mongo.Conexion)
	if err != nil {
		return err
	}

	// -----------------------------------------------------------------------
	// Resend
	clientResend := resend.NewClient(cfg.Resend.APIKEY)

	// -----------------------------------------------------------------------
	// Repositories

	var (
		coursesRepository             = coursesdb.NewRepository(db)
		userRepository                = userdb.NewRepository(db)
		learningRepository            = learningdb.NewRepository(db)
		assignmentsRepository         = assigmentsdb.NewRepository(db)
		assignmentsRepositoryFlexible = assigmentsdb.NewRepositoryFlexible(db)
		examReposotory                = examdb.NewRepository(db)
		payinRepository               = payindb.New(db)
		movementsRepository           = movementsdb.NewRepository(db)
	)

	// -----------------------------------------------------------------------
	// Services
	emailService := email.NewService(clientResend)

	payinService, err := payin.NewService(payinRepository, movementsRepository, assignmentsRepositoryFlexible, payin.Config{
		MercadoPago: payin.MercadoPagoConfig{
			AccessToken: cfg.MercadoPago.AccessToken,
		},
	}, emailService)
	if err != nil {
		return fmt.Errorf("erro al inicializar payin: %w", err)
	}

	var (
		coursesService  = courses.NewService(coursesRepository)
		userService     = user.NewService(userRepository)
		learningService = learning.NewService(learningRepository)
		oauthService    = oauth.NewService(oauth.Config{
			GoogleClientID: "947017986235-hjvh14vf1mnh04drpnpvapav5bh2oqh7.apps.googleusercontent.com",
		}, userService)
		examService       = exam.NewService(examReposotory, userRepository)
		assigmentsService = assignments.NewService(assignmentsRepository, payinService, examService)
	)

	// -----------------------------------------------------------------------
	// Routes

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		// Permite cualquier origen (usar "*" es común en desarrollo con apps móviles/web)
		// Para producción, se recomienda cambiar "*" por tu dominio específico (ej: "https://itydu.app")
		AllowedOrigins: []string{"*"},

		// Métodos HTTP permitidos
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},

		// Headers permitidos
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},

		// Permite enviar credenciales (como cookies)
		AllowCredentials: true,

		// Tiempo que el navegador puede cachear los resultados de la pre-flight (en segundos)
		MaxAge: 300,
	}))
	r.Use(middleware.Logger)
	courses.MakeHandlerWith(coursesService).SetRoutesTo(r)
	user.MakeHandlerWith(userService).SetRoutesTo(r)
	learning.MakeHandlerWith(learningService).SetRoutesTo(r)
	oauth.MakeHandlerWith(oauthService).SetRoutesTo(r)
	assignments.MakeHandlerWith(assigmentsService).SetRoutesTo(r)
	exam.NewHttpHandler(examService).SetRoutesTo(r)
	payin.MakeHandlerWith(payinService).SetRoutesTo(r)

	// -------------------------------------------------------------------------
	// HTTP App Server

	var (
		shutdownListener = make(chan os.Signal, 1)
		errListener      = make(chan error, 1)
	)

	signal.Notify(shutdownListener, syscall.SIGINT, syscall.SIGTERM)

	api := http.Server{
		Addr:         cfg.Web.Host,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		Handler:      r,
	}

	go func() {
		log.Info(ctx, "Startup - API router started", "host", api.Addr)

		errListener <- api.ListenAndServe()
	}()

	// -----------------------------------------------------------------------
	// Shutdown

	select {
	case err := <-errListener:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdownListener:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown completed", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("cannot stop server gracefully: %w", err)
		}
	}

	return nil
}

func traceFunc(ctx context.Context) []any {
	v := tracer.GetValues(ctx)

	fields := make([]any, 2, 4)
	fields[0], fields[1] = "traceID", v.TraceID

	if v.Rut != "" {
		fields = append(fields, "RUT", v.Rut)
	}

	return fields
}
