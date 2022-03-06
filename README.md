# News Redis

API to manage News and Topic(Tags) With Redis caching.

News and Topic management

- CRUD on News and Tags
- One news can contains multiple tags e.g. "Safe investment" might contains tags "investment", "mutual fund", etc
- One topic has multiple news e.g. "investment" topic might contains "how to start investment", "mutual fund is safe investment type", etc
- Enable filter by news status ("draft", "deleted", "publish")
- Enable filter news by its topics

## API Documentation

- [News and Topic management](https://app.swaggerhub.com/apis-docs/furqonzt99/newsapi/1) - Swagger API Documentation

## Installation

- Clone this project

```
git clone https://github.com/furqonzt99/news-redis.git news-redis
```

- Go to project folder

```
cd news-redis
```

- Install dependencies

```
go mod tidy
```

- Dont Forget to setup .env file on project folder and run this command

```
go run .
```

## Testing

- Dont Forget to setup .env file on test folder and run this command

```
go test -p 1 -v ./test
```

![Test Result](https://github.com/furqonzt99/news-redis/blob/main/docs/test.png)

## Caching Strategies

![Caching Strategis 1](https://github.com/furqonzt99/news-redis/blob/main/docs/cs1.png)
![Caching Strategis 2](https://github.com/furqonzt99/news-redis/blob/main/docs/cs2.png)
