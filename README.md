# TapMenu

## API for restaurants

1. Create `config.toml` file by example.
2. Up broker and database.
```bash
cd local-stand
docker-compose -f kafka/docker-compose.yml up -d
docker-compose -f tarantool/docker-compose.yml up -d
```
3. Build and start app.
```bash
cd .. && make
./tapmenu
```
[consumer-part](https://github.com/alex-pvl/go-tapmenu-consumer)