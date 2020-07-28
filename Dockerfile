FROM alpine
RUN apk update && apk add ca-certificates
COPY trip-data-receiver-bot  .
COPY config.toml .
CMD [ "./trip-data-receiver-bot" ]