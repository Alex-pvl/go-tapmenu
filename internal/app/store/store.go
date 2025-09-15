package store

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/alex-pvl/go-tapmenu/internal/app/config"
	"github.com/tarantool/go-tarantool/v2"
	"github.com/tarantool/go-tarantool/v2/datetime"
)

type Store struct {
	config *config.Configuration
	conn   *tarantool.Connection
	logger *logrus.Logger
}

func New(config *config.Configuration, logger *logrus.Logger) *Store {
	s := &Store{config: config, logger: logger}
	if err := s.connect(); err != nil {
		logger.Error(err)
	}
	return s
}

func (s *Store) connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.Timeout)*time.Second)
	defer cancel()

	dialer := tarantool.NetDialer{
		Address:  s.config.TarantoolAddress,
		User:     s.config.Username,
		Password: s.config.Password,
	}
	opts := tarantool.Opts{
		Timeout: time.Duration(s.config.Timeout) * time.Second,
	}

	conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		return err
	}

	s.conn = conn
	return nil
}

func (s *Store) GetTable(hash string) (*Table, error) {
	selectRequest := tarantool.NewSelectRequest(s.config.TablesSpaceId).Key([]interface{}{hash})
	resp, err := s.conn.Do(selectRequest).Get()
	if err != nil {
		s.logger.Error(err)
	}
	return mapToTable(resp)
}

func (s *Store) UpdateTable(hash string, table *Table) error {
	lastCallTnt, _ := datetime.MakeDatetime(table.LastCall)
	replaceRequest := tarantool.NewReplaceRequest(s.config.TablesSpaceId).Tuple([]interface{}{
		hash, table.Url, table.RestaurantName, table.Number, lastCallTnt,
	})
	_, err := s.conn.Do(replaceRequest).Get()
	return err
}

func (s *Store) FindAndDeleteExistingCall(tableNumber int8) {
	selectRequest := tarantool.NewSelectRequest(s.config.OrdersSpaceId)
	resp, err := s.conn.Do(selectRequest).Get()
	if err != nil {
		s.logger.Error(err)
	}

	var idToDelete string
	for _, line := range resp {
		order := *mapFromInterface(line.([]interface{}))

		if order.Accepted {
			continue
		}

		if order.TableNumber == int(tableNumber) {
			idToDelete = order.Id.String()
			break
		}
	}

	if idToDelete == "" {
		return
	}

	deleteRequest := tarantool.NewDeleteRequest(s.config.OrdersSpaceId).Key([]interface{}{idToDelete})

	_, err = s.conn.Do(deleteRequest).Get()
	if err != nil {
		s.logger.Error(err)
	}
}

func (s *Store) CreateCall(order Order) error {
	createdAt, _ := datetime.MakeDatetime(order.CreatedAt)
	updatedAt, _ := datetime.MakeDatetime(order.UpdatedAt)

	insertRequest := tarantool.NewReplaceRequest(s.config.OrdersSpaceId).Tuple([]interface{}{
		order.Id.String(),
		order.RestaurantName,
		order.TableNumber,
		createdAt,
		updatedAt,
		order.Accepted,
	})
	_, err := s.conn.Do(insertRequest).Get()
	return err
}
