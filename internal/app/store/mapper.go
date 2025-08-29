package store

import (
	"errors"
	"github.com/google/uuid"
	"github.com/tarantool/go-tarantool/v2/datetime"
)

func mapToTable(dbResponse []interface{}) (*Table, error) {
	if len(dbResponse) == 0 {
		return nil, errors.New("table not found")
	}
	row, ok := dbResponse[0].([]interface{})
	if !ok {
		return nil, errors.New("cannot cast response")
	}

	url, _ := row[1].(string)
	restName, _ := row[2].(string)
	number, _ := row[3].(int8)
	lastCallTnt, _ := row[4].(datetime.Datetime)
	lastCall := lastCallTnt.ToTime()

	return &Table{
		Url:            url,
		RestaurantName: restName,
		Number:         int(number),
		LastCall:       lastCall,
	}, nil
}

func mapFromInterface(row []interface{}) *Order {
	id, _ := row[0].(string)
	restName, _ := row[1].(string)
	number, _ := row[2].(int8)
	createdAtTnt, _ := row[3].(datetime.Datetime)
	updatedAtTnt, _ := row[4].(datetime.Datetime)
	accepted, _ := row[5].(bool)

	createdAt := createdAtTnt.ToTime()
	updatedAt := updatedAtTnt.ToTime()

	return &Order{
		Id:             uuid.MustParse(id),
		RestaurantName: restName,
		TableNumber:    int(number),
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		Accepted:       accepted,
	}
}
