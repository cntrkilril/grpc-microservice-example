package app

import (
	v1 "appointment-service/internal/controller/grpc/v1"
	natsController "appointment-service/internal/controller/nats"
	"appointment-service/internal/gateway"
	"appointment-service/internal/service"
	"github.com/nats-io/nats.go"
	gen "gitlab.com/d1zero-online-booking/common/pb/gen/appointment_service"
	"gitlab.com/d1zero-online-booking/common/pkg/govalidator"
	"gitlab.com/d1zero-online-booking/common/pkg/logger"
	postgresCommon "gitlab.com/d1zero-online-booking/common/pkg/postgres"
	gogrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"
	// migrations
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
)

func Run() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	l := logger.NewLogger(logger.Config{Logger: cfg.Logger})
	l.Infof("logger initialized successfully")

	db := postgresCommon.InitPsqlDB(postgresCommon.Config{
		Postgres: cfg.Postgres,
	}, l)
	l.Debug("Connected to PostgreSQL")

	clientServiceConn, err := gogrpc.Dial(
		cfg.ClientService.ConnString,
		gogrpc.WithTransportCredentials(insecure.NewCredentials()),
		gogrpc.FailOnNonTempDialError(true),
		gogrpc.WithBlock(),
	)
	if err != nil {
		l.Fatalf("failed to connect to client service: %v", err)
	}
	l.Infof("client service connected successfully")

	servicesServiceConn, err := gogrpc.Dial(
		cfg.ServicesService.ConnString,
		gogrpc.WithTransportCredentials(insecure.NewCredentials()),
		gogrpc.FailOnNonTempDialError(true),
		gogrpc.WithBlock(),
	)
	if err != nil {
		l.Fatalf("failed to connect to services service: %v", err)
	}
	l.Infof("services service connected successfully")

	masterServiceConn, err := gogrpc.Dial(
		cfg.MasterService.ConnString,
		gogrpc.WithTransportCredentials(insecure.NewCredentials()),
		gogrpc.FailOnNonTempDialError(true),
		gogrpc.WithBlock(),
	)
	if err != nil {
		l.Fatalf("failed to connect to master service: %v", err)
	}
	l.Infof("master service connected successfully")

	adminServiceConn, err := gogrpc.Dial(
		cfg.ClientService.ConnString,
		gogrpc.WithTransportCredentials(insecure.NewCredentials()),
		gogrpc.FailOnNonTempDialError(true),
		gogrpc.WithBlock(),
	)
	if err != nil {
		l.Fatalf("failed to connect to client service: %v", err)
	}
	l.Infof("client service connected successfully")

	val := govalidator.New()

	// gateways
	appointmentGateway := gateway.NewAppointmentRepository(db)
	workTimeGateway := gateway.NewWorkTimeRepository(db)
	timeCellGateway := gateway.NewTimeCellRepository(db)
	clientGateway := gateway.NewClientRepository(clientServiceConn)
	servicesGateway := gateway.NewServicesRepository(servicesServiceConn)
	masterGateway := gateway.NewMasterRepository(masterServiceConn)
	adminGateway := gateway.NewAdminRepository(adminServiceConn)
	registryGateway := gateway.NewPGRegistry(db)

	// services
	getService := service.NewGetService(appointmentGateway, workTimeGateway, servicesGateway, timeCellGateway, cfg.SpacerTimeCell.Value)
	deleteService := service.NewDeleteService(registryGateway)
	createService := service.NewCreateService(registryGateway, clientGateway, servicesGateway, masterGateway, cfg.SpacerTimeCell.Value)
	updateService := service.NewUpdateService(registryGateway, clientGateway, masterGateway, adminGateway, servicesGateway, cfg.SpacerTimeCell.Value)

	// controllers
	appointmentController := v1.NewAppointmentController(createService, getService, updateService, val)
	workTimeController := v1.NewWorkTimeController(deleteService, getService, updateService, val)
	timeCellController := v1.NewTimeCellController(getService, val)

	// servers
	grpcServer := gogrpc.NewServer()
	defer grpcServer.GracefulStop()

	gen.RegisterAppointmentServiceV1Server(grpcServer, appointmentController)
	gen.RegisterWorkTimeServiceV1Server(grpcServer, workTimeController)
	gen.RegisterTimeCellServiceV1Server(grpcServer, timeCellController)

	go func() {
		lis, err := net.Listen("tcp", net.JoinHostPort(cfg.GRPC.Host, cfg.GRPC.Port))
		if err != nil {
			log.Fatalf("tcp sock: %s", err.Error())
		}
		defer func(lis net.Listener) {
			err = lis.Close()
			if err != nil {
				l.Error(err)
				return
			}
		}(lis)

		err = grpcServer.Serve(lis)
		if err != nil {
			log.Fatalf("GRPC server: %s", err.Error())
		}
	}()

	l.Debug("Started GRPC server")

	nc, err := nats.Connect(cfg.NATS.ConnString)
	if err != nil {
		if err != nil {
			l.Error(err)
			return
		}
	}

	js, err := nc.JetStream()
	if err != nil {
		if err != nil {
			l.Error(err)
			return
		}
	}

	l.Debug("Connected to NATS")

	natsCtrl := natsController.NewWorkTimeController(js, createService, l)
	if err = natsCtrl.Subscribe(); err != nil {
		if err != nil {
			l.Error(err)
			return
		}
	}

	l.Debug("Registered nats handlers")

	l.Debug("Application has started")

	exit := make(chan os.Signal, 2)

	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	<-exit

	l.Info("Application has been shut down")

	return
}
