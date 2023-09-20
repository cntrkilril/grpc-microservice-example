package app

import (
	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
	"gitlab.com/d1zero-online-booking/common/pkg/logger"
	"gitlab.com/d1zero-online-booking/common/pkg/postgres"
	"strings"
)

type (
	Config struct {
		Logger          logger.Logger     `koanf:"logger" validate:"required"`
		GRPC            GRPC              `koanf:"grpc" validate:"required"`
		Postgres        postgres.Postgres `koanf:"postgres" validate:"required"`
		NATS            NATS              `koanf:"nats" validate:"required"`
		ClientService   ClientService     `koanf:"clientService" validate:"required"`
		ServicesService ServicesService   `koanf:"servicesService" validate:"required"`
		MasterService   MasterService     `koanf:"masterService" validate:"required"`
		AdminService    AdminService      `koanf:"adminService" validate:"required"`
		SpacerTimeCell  SpacerTimeCell    `koanf:"spacerTimeCell" validate:"required"`
	}

	NATS struct {
		ConnString string `koanf:"connString" validate:"required"`
	}

	GRPC struct {
		Host string `koanf:"host" validate:"required"`
		Port string `koanf:"port" validate:"required"`
	}

	ClientService struct {
		ConnString string `koanf:"connString" validate:"required"`
	}

	ServicesService struct {
		ConnString string `koanf:"connString" validate:"required"`
	}

	MasterService struct {
		ConnString string `koanf:"connString" validate:"required"`
	}

	AdminService struct {
		ConnString string `koanf:"connString" validate:"required"`
	}

	SpacerTimeCell struct {
		Value int `koanf:"value" validate:"required"`
	}
)

func LoadConfig() (*Config, error) {

	k := koanf.New(".")
	defaultLogLevel := int8(-1)
	p := env.Provider("", ".", func(s string) string {
		return strings.ToLower(strings.ReplaceAll(s, "_", "."))
	})

	if err := k.Load(p, nil); err != nil {
		return nil, err
	}

	cfg := &Config{
		GRPC: GRPC{
			Host: "localhost",
			Port: "8075",
		},
		Logger: logger.Logger{
			Level: defaultLogLevel,
		},
		NATS: NATS{
			ConnString: "nats://127.0.0.1:4223",
		},

		Postgres: postgres.Postgres{
			ConnString:      "postgresql://root:pass@127.0.0.1:5432/appointments?sslmode=disable&application_name=appointment-service",
			MaxOpenConns:    10,
			ConnMaxLifetime: 20,
			MaxIdleConns:    15,
			ConnMaxIdleTime: 30,
			AutoMigrate:     true,
			MigrationsPath:  "db/migration",
			DBName:          "appointments",
		},
		ClientService: ClientService{
			ConnString: "127.0.0.1:8070",
		},
		ServicesService: ServicesService{
			ConnString: "127.0.0.1:8010",
		},
		MasterService: MasterService{
			ConnString: "127.0.0.1:8065",
		},
		AdminService: AdminService{
			ConnString: "127.0.0.1:8000",
		},
		SpacerTimeCell: SpacerTimeCell{
			Value: 15,
		},
	}

	if err := k.Unmarshal("", cfg); err != nil {
		return nil, err
	}

	err := validator.New().Struct(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
