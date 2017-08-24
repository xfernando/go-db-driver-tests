<img src="hugging-docker.svg" width="40%">[1]

## go-db-driver-tests

This project is based on Brad Fitz's [earlier work](https://github.com/bradfitz/go-sql-test) and aims to provide a testing suite for go [database drivers
](https://github.com/golang/go/wiki/SQLDrivers) using docker to start the database servers needed to run all the tests.

Progress:

MySQL drivers tested: 2 out of 2

Postgres drivers teste: 2 out of 3

## Usage

We use `docker-compose` to build the image with the test code and start all the database containers needed for the tests to run:

```bash
docker-compose build
docker-compose up -d
```

Then, you can see the output of the test run with:
```bash
docker-compose logs -f godbtests
```

When the tests are completed, run the following to stop the database containers:
```bash
docker-compose down
```


[1] Image copied from https://github.com/egonelbre/gophers