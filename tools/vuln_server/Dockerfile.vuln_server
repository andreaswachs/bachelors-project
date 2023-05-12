FROM alpine:latest as base

RUN apk add --no-cache python3

FROM base as compose

EXPOSE 80
WORKDIR /app

RUN touch flag.txt
COPY boot.sh boot.sh
CMD [ "./boot.sh" ]
