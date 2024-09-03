# AS
For connection of PostgreSQL Database,
Please define environment file under working directory.
> .env
```
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

> marketdata.json
```
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
- GET
  1. /api/hc 
     - Health Check of API
  2. /api/v1/jwt/refresh 
     - Refresh JWT timelive
  3. /api/v1/order 
     - List Order by User
  4. /api/v1/order/{Order ID} 
     - Get Specific Order 
  5. /api/v1/user 
     - Get User infomations
  6. /api/v1/stock 
     - List Stock on Market
  7. /api/v1/stock/{Symbol} 
     - Get Specific Stock
- POST
1. /api/v1/register 
    - User Registration
    - Body:
```
{   
    "name": "test",
    "mobile": "12345678",
    "email": "test@test1.com",
    "password": "123456",
    "hkid": "Y1234567"
}
```
  2. /api/v1/login 
     - Login Method to get JWT
     - Body:
```
{
    "email": "test@test.com",
    "password": "123456"
}
``` 
  3. /api/v1/order/create 
     - Create Order by User
     - Body: 
```
{
        "method":"buy/sell",
        "order":"limit/price",
        "place":"standard/bid",
        "symbol":"0001",
        "price":123,
        quantity:123
}
```
  4. /api/v1/order/{Order ID}/cancel 
     - Cancel Order by User. But this is not Remove record on DataBase. 
     - it is change order status to "Cancelled"
  5. /api/v1/user/update/bankinfo 
     - Setup bank account
     - Body:
```
{
    "bank": "123",
    "branch": "123",
    "account": "123123123"
}
```
  6. /api/v1/user/update/password
     - Change Password
     - Body:
```
{
    "currpwd": "123",
    "newpwd": "123"
}
```
  7. /api/v1/user/update/userinfo
     - Edit User informations
     - Body:
```
{
    "name": "fed",
    "email": "abc@def.com",
    "mobile": "98765432"
}
```