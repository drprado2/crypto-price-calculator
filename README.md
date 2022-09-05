# Crypto VWAP Calculator

The project contains one WEB API and one Worker exposing the same domain, the domain persist [match](https://docs.cloud.coinbase.com/exchange/docs/websocket-channels#match) events and calculate the VWAP of the latest N trades.
The Worker consumes the [Coinbase's Websocket](https://docs.cloud.coinbase.com/exchange/docs/websocket-overview).
The project contains a docker-compose to deploy all dependencies locally, such as PostgreSQL, Grafana, Jaeger.

## Running API

* Start dependencies with docker
```sh
make start-local-dependencies
```
* Start web API
```sh
make run-web-api
```
* Send a request
```sh
curl --location --request POST 'http://localhost:5050/api/internal/v1/trade-order' \
--header 'Content-Type: application/json' \
--data-raw '{
    "TradeId": 241,
    "MakerOrderId": "df315bb9-c940-4f55-be69-ea165f987e35",
    "TakerOrderId": "6d5be16b-2936-49ad-883f-78860c94891c",
    "Side": "buy",
    "Size": "0.580621",
    "Price": "18924.09",
    "ProductId": "BTC-USD",
    "Sequence": 44805966539,
    "Time": "2022-09-07T14:38:51.473515Z"
}'
```

## Running Worker

* Start dependencies with docker
```sh
make start-local-dependencies
```
* Start Worker
```sh
make run-worker
```

## Running Tests
The command execute the tests, generate the report files and open in the browser, if your browser doesn't open,
you can manually open with the cover.html file created on the project root.
 
* Execute tests
```sh
make test-cover
```
* Clear cover files
```sh
make clear-cover
```

## SQL Schema
* The project uses [PostgreSQL](https://www.postgresql.org/) as DB and [Flyway](https://flywaydb.org/) as migration tool,
the DDL scripts could be find on `eng/sql/structure_sql_scripts`.
* DML script can be added on `data_sql_scripts`.
* The scripts are automatically applied when the docker dependencies are started.

## Tracing
* The application use [Opentelemetry](https://opentelemetry.io/docs/instrumentation/go/) to produce tracers.
* The Jaeger is deployed with the docker dependencies and could be accessed to check the tracings.
* open Jaeger
```sh
make open-jaeger
```

## Logs
* The application use [Logrus](https://pkg.go.dev/github.com/sirupsen/logrus) to produce logs.
* The logs are printed on the `stdout` and also persisted on a file that is shared with the `promtail` container, 
these persisted logs can be accessed through Grafana.
* open Grafana
```sh
make open-grafana
```
* Use the default user `admin` and password `admin` in the first access.
* Click on [`Add your first data source`](http://localhost:3000/datasources/new?utm_source=grafana_gettingstarted) in order to add Loki as data source.
* Filter and select Loki
* Fill the URL input with `http://loki:3100` and click Save and Test, you should see a success message.
* Now you are able to execute queries on the JSON structured logs, or even create some dashboard.




