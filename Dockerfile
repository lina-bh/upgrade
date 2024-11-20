FROM docker.io/library/golang:1-alpine AS build
WORKDIR /src
COPY . /src
RUN go build

FROM scratch
COPY --from=build /src/upgrade /upgrade
USER 1000
CMD ["/upgrade"]
