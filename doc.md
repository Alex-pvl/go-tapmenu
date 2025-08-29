#### GET /tapmenu/{hash}
Response
```json
{
  "url": "menu.ru",
  "restaurant_name": "rest_name",
  "number": 5,
  "last_call": "2025-09-02T17:13:01.282703Z"
}
```
#### POST /tapmenu/{hash}/call
Response
```json
{
  "id": "9eae72d9-890f-4a11-87f3-d5632d989493",
  "restaurant_name": "rest_name",
  "table_number": 5,
  "created_at": "2025-09-02T17:13:01.282703Z",
  "updated_at": "2025-09-02T17:13:01.282703Z",
  "accepted": false
}
```