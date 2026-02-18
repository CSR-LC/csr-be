FROM golang:1.25.7
COPY . /app
WORKDIR /app
RUN make setup \
&&  make generate
USER root
CMD ["go" ,"run" ,"./cmd/swagger/"]
