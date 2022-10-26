package store

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/golang-migrate/migrate/v4"
	ps "github.com/golang-migrate/migrate/v4/database/postgres"

	// only need this for gorm
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type IStore interface {
	Account() IAccountStore
	Driver() IDriverStore
	Location() ILocationStore
	Booking() IBookingStore
	Notification() INotificationStore
}

type Store struct {
	conf *Config
	db   *gorm.DB
}

func New() (IStore, error) {
	conf, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config, err: %v", err)
	}

	db, err := newDB(conf)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get *sql.DB, err: %v", err)
	}

	sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
	sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(conf.ConnLifeTime) * time.Second)

	if err = sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database, err: %v", err)
	}

	s := Store{conf: conf, db: db}

	err = s.migrate()
	if err != nil {
		return nil, fmt.Errorf("failed to migrate, err: %v", err)
	}

	return &s, nil
}

func newDB(conf *Config) (*gorm.DB, error) {
	db, err := gorm.Open(
		postgres.Open(conf.ConnString()),
		&gorm.Config{
			AllowGlobalUpdate: false,
			Logger:            logger.Default.LogMode(logger.LogLevel(conf.LogLevel)),
			NowFunc:           func() time.Time { return time.Now().UTC() },
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection, err: %v", err)
	}

	return db, nil
}

func (s *Store) migrate() (err error) {
	gormDB, err := newDB(s.conf)
	if err != nil {
		return err
	}

	db, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get *sql.DB, err: %v", err)
	}

	driver, err := ps.WithInstance(db, &ps.Config{MigrationsTable: s.conf.MigrateTable})
	if err != nil {
		return fmt.Errorf("failed to prepare driver, err: %v", err)
	}

	mg, err := migrate.NewWithDatabaseInstance(s.conf.MigrateSource, s.conf.Driver, driver)
	if err != nil {
		return fmt.Errorf("failed to create migrator, err: %v", err)
	}

	defer func() {
		sourceErr, dbErr := mg.Close()
		if sourceErr != nil {
			err = fmt.Errorf("failed to close source, err: %s", sourceErr.Error())
		}

		if dbErr != nil {
			err = fmt.Errorf("failed to close database, err: %s", dbErr.Error())
		}
	}()

	if err = mg.Migrate(s.conf.MigrateVersion); err != nil {
		if err == migrate.ErrNoChange {
			err = nil
		} else if err != nil {
			return fmt.Errorf("failed to migrate, err: %v", err)
		}
	}

	return nil
}

func (s Store) Account() IAccountStore {
	return accountStore{s}
}

func (s Store) Driver() IDriverStore {
	return driverStore{s}
}

func (s Store) Location() ILocationStore {
	return locationStore{s}
}

func (s Store) Booking() IBookingStore {
	return bookingStore{s}
}

func (s Store) Notification() INotificationStore {
	return notificationStore{s}
}
