# workspace
Demo

Run the below docker command to execute the application

docker build -t my-app .
docker run -p 8080:8080 -it my-app

Sample http requests to test:
http://localhost:8080/currency/all
http://localhost:8080/currency/BTCUSD
http://localhost:8080/currency/ETHBTC