package main

// junta toda a estrutura dos arquivos
import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"

	httpHandler "github.com/GuiRissato/ticket-api/internal/events/infra/http"

	"github.com/GuiRissato/ticket-api/internal/events/infra/repository"
	"github.com/GuiRissato/ticket-api/internal/events/infra/service"
	"github.com/GuiRissato/ticket-api/internal/events/usecase"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	db, err := sql.Open("mysql", "test_user:test_password@tcp(localhost:3306)/test_db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Repositório
	eventRepo, err := repository.NewMysqlEventRepository(db)
	if err != nil {
		log.Fatal(err)
	}

	// URLs base específicas para cada parceiro
	partnerBaseURLs := map[int]string{
		1: "http://localhost:9000/api1",
		2: "http://localhost:9000/api2",
	}

	listEventsUseCase := usecase.NewListEventsUseCase(eventRepo)
	getEventUseCase := usecase.NewGetEventUseCase(eventRepo)
	createEventUseCase := usecase.NewCreateEventUseCase(eventRepo)
	partnerFactory := service.NewPartnerFactory(partnerBaseURLs)
	buyTicketsUseCase := usecase.NewBuyTicketsUseCase(eventRepo, partnerFactory)
	createSpotsUseCase := usecase.NewCreateSpotsUseCase(eventRepo)
	listSpotsUseCase := usecase.NewListSpotsUseCase(eventRepo)

	// Handlers HTTP
	eventsHandler := httpHandler.NewEventsHandler(
		listEventsUseCase,
		getEventUseCase,
		createEventUseCase,
		buyTicketsUseCase,
		createSpotsUseCase,
		listSpotsUseCase,
	)

	r := http.NewServeMux()
	r.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	r.HandleFunc("/events", eventsHandler.ListEvents)
	r.HandleFunc("/events/{eventID}", eventsHandler.GetEvent)
	r.HandleFunc("/events/{eventID}/spots", eventsHandler.ListSpots)
	r.HandleFunc("POST /events", eventsHandler.CreateEvent)
	r.HandleFunc("POST /checkout", eventsHandler.BuyTickets)
	r.HandleFunc("POST /events/{eventID}/spots", eventsHandler.CreateSpots)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Canal para escutar sinais do sistema operacional
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		// Recebido sinal de interrupção, iniciando o graceful shutdown
		log.Println("Recebido sinal de interrupção, iniciando o graceful shutdown...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Erro no graceful shutdown: %v\n", err)
		}
		close(idleConnsClosed)
	}()

	// Iniciando o servidor HTTP
	log.Println("Servidor HTTP rodando na porta 8080")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Erro ao iniciar o servidor HTTP: %v\n", err)
	}

	<-idleConnsClosed
	log.Println("Servidor HTTP finalizado")

}
