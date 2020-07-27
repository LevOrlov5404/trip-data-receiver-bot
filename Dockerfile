FROM alpine
RUN apk update && apk add ca-certificates
ADD trip-data-receiver-bot  trip-data-receiver-bot
CMD [ "./trip-data-receiver-bot" ]