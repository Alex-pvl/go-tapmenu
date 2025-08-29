package store

import (
	"context"
	"log"
	"time"

	"github.com/alex-pvl/go-tapmenu/internal/app/config"
	"github.com/tarantool/go-tarantool/v2"
	"github.com/tarantool/go-tarantool/v2/datetime"
)

const (
	tablesSpaceId = 513
	ordersSpaceId = 514
)

type Store struct {
	config *config.Configuration
	conn   *tarantool.Connection
}

func New(config *config.Configuration) *Store {
	s := &Store{config: config}
	if err := s.connect(); err != nil {
		log.Fatal(err)
	}
	return s
}

func (s *Store) connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.Timeout)*time.Second)
	defer cancel()

	dialer := tarantool.NetDialer{
		Address:  s.config.TarantooldbAddress,
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
	selectRequest := tarantool.NewSelectRequest(tablesSpaceId).Key([]interface{}{hash})
	resp, err := s.conn.Do(selectRequest).Get()
	if err != nil {
		log.Fatal(err)
	}
	return mapToTable(resp)
}

func (s *Store) UpdateTable(hash string, table *Table) error {
	lastCallTnt, _ := datetime.MakeDatetime(table.LastCall)
	replaceRequest := tarantool.NewReplaceRequest(tablesSpaceId).Tuple([]interface{}{
		hash, table.Url, table.RestaurantName, table.Number, lastCallTnt,
	})
	_, err := s.conn.Do(replaceRequest).Get()
	return err
}

func (s *Store) FindAndDeleteExistingCall(tableNumber int8) {
	selectRequest := tarantool.NewSelectRequest(ordersSpaceId)
	resp, err := s.conn.Do(selectRequest).Get()
	if err != nil {
		log.Fatal(err)
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

	deleteRequest := tarantool.NewDeleteRequest(ordersSpaceId).Key([]interface{}{idToDelete})

	_, err = s.conn.Do(deleteRequest).Get()
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Store) CreateCall(order Order) error {
	createdAt, _ := datetime.MakeDatetime(order.CreatedAt)
	updatedAt, _ := datetime.MakeDatetime(order.UpdatedAt)

	insertRequest := tarantool.NewReplaceRequest(ordersSpaceId).Tuple([]interface{}{
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
