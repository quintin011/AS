# AS
For connection of PostgreSQL Database,
Please define environment file under working directory.
```
.env
----
SERVICE_PORT=10001
DB_HOST=
DB_PORT=5432
DB_USER=postgres
DB_NAME=postgres
DB_PASSWORD=
```
| Name          | Description                                   |
| ------------- | --------------------------------------------- |
| SERVICE_PORT  | API service port. Default: 10001              |
| DB_HOST       | PostgreSQL Database IP/Hostname               |
| DB_PORT       | PostgreSQL Database Port. Default: 5432       |
| DB_USER       | PostgreSQL Operations User. Default: postgres |
| DB_NAME       | Schema Name. Default: postgres                |
| DB_PASSWORD   | Password of PostgreSQL Operations User        |

For Stocks data, please self define marketdata.json file in stocks directory.
Services will make new stocks directory when stocks is not exist on working directory.
```
marketdata.json
---
[
    {
        "symbol": string,
        "updated_at": string,
        "currbid": float32,
        "currask": float32,
        "lasttrade": float32,
        "high_price": float32,
        "low_price": float32,
        "vol": int
    },
    ...
]
```
| Name          | Description                                                |
| ------------- | ---------------------------------------------------------- |
| symbol        | Stock ID, Size: 4                                          |
| updated_at    | Update time, Timestamp, format: yyyy-mm-ddTHH:MM:SS.ssssss |
| currbid       | Current Bid Price                                          |
| currask       | Current Ask Price                                          |
| lasttrade     | Last Transaction Price                                     |
| high_price    | Highest Transcation Price                                  |
| low_price     | Lowest Transcation Price                                   |
| vol           | Stock Quantity                                             |
```