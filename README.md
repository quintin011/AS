# AS
pre-requirements
```
awscli
node 20
make
golang >=1.22
```
pre-start
```
aws configure
<input aws accesskey>
<input aws secrets>
<region>
<press enter>
```
if you want to multiple user for divide environment
```
aws configure --profile
```
install serverless & read .env plugin
```
npm i -g serverless
npm i -D serverless-dotenv-plugin
```
for first start, please sign-up serverless and make initialize serverless
```
serverless
```
you can login by browser or using license key for can't login serverless and want to create same organization
please check serverless.yaml to define stage 
when you had multiple aws accesskey on local please define below on Makefile
```
build:
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bootstrap main.go 
deploy: build
	sls deploy --aws-profile <profile>
clean:
	rm -rf ./bin ./vendor Gopkg.lock ./serverless
``` 
for start service
```
make depoly
```
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
  6. /api/v1/user/pos
     - Get Position of user
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
1. /api/v1/login 
     - Login Method to get JWT
     - Body:
```
{
    "email": "test@test.com",
    "password": "123456"
}
``` 
1. /api/v1/order/create 
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
1. /api/v1/order/{Order ID}/cancel 
     - Cancel Order by User. But this is not Remove record on DataBase. 
     - it is change order status to "Cancelled"
2. /api/v1/user/update/bankinfo 
     - Setup bank account
     - Body:
```
{
    "bank": "123",
    "branch": "123",
    "account": "123123123"
}
```
3. /api/v1/user/update/password
     - Change Password
     - Body:
```
{
    "currpwd": "123",
    "newpwd": "123"
}
```
4. /api/v1/user/update/userinfo
     - Edit User informations
     - Body:
```
{
    "name": "fed",
    "email": "abc@def.com",
    "mobile": "98765432"
}
```
5. /api/v1/user/addbalance
   - Add Balance
   - Body:
```
{
   "balance": 10000
}
```
6. /api/v1/user/test/add/positions?symbol=0001&quan=1000
    - just for test trading!!!!
    - symbol: stocks symbol, type: string
    - quan: Quantity, type: int