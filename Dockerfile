FROM golang:1.24
COPY . /app
WORKDIR /app
RUN make setup \
&&  make generate
USER root
CMD ["go" ,"run" ,"./cmd/swagger/"]
