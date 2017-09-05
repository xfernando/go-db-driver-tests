FROM golang:1.9

RUN apt-get update && apt-get install -y libpq-dev libsqlite3-dev freetds-dev
ENV GODEBUG=cgocheck=0
# see https://github.com/jbarham/gopgsqldriver/issues/4
RUN ln -s /usr/include/postgresql/libpq-fe.h /usr/include/ && ln -s /usr/include/postgresql/postgres_ext.h /usr/include/ && ln -s /usr/include/postgresql/pg_config_ext.h /usr/include/

WORKDIR /go/src/app
COPY ./ .
WORKDIR /go/src/app/sqltest
RUN go get -v ./...
# Register gopgsql driver with another name, otherwise we get an error for registering driver with same name twice
RUN sed -i 's/(\"postgres\"/(\"gopgsql\"/' /go/src/github.com/jbarham/gopgsqldriver/pgdriver.go
# Rename sqlite drivers since they all use the same name
RUN sed -i 's/sql.Register(\"sqlite3\"/sql.Register(\"gwenn_sqlite3\"/' /go/src/github.com/gwenn/gosqlite/driver.go && \
   sed -i 's/register(\"sqlite3\")/register(\"mxk_sqlite3\")/' /go/src/github.com/mxk/go-sqlite/sqlite3/sqlite3.go && \
   sed -i 's/sql.Register(\"mssql\"/sql.Register(\"denisenkom_mssql\"/' /go/src/github.com/denisenkom/go-mssqldb/mssql.go
#  sed -i 's/sql.Register(\"sqlite3\"/sql.Register(\"mattn_sqlite3\"/' /go/src/github.com/mattn/go-sqlite3/sqlite3.go
RUN rm -f drivers.go

CMD ["go", "test", "-v"]
