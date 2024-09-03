# AS
```
.env
----
DB_HOST=
DB_PORT=5432
DB_USER=postgres
DB_NAME=postgres
DB_PASSWORD=
```
| Name          | Description                                   |
| ------------- | --------------------------------------------- |
| DB_HOST       | PostgreSQL Database IP/Hostname               |
| DB_PORT       | PostgreSQL Database Port. Default: 5432       |
| DB_USER       | PostgreSQL Operations User. Default: postgres |
| DB_NAME       | Schema Name. Default: postgres                |
| DB_PASSWORD   | Pasword of PostgreSQL Operations User         |

For Stocks data, please self define marketdata.json file in stocks directory.
Services will make new stocks directory when stocks is not exist on working directory.

