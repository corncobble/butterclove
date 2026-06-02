FROM golang:1.26.3-trixie AS build
WORKDIR /app
COPY . ./
RUN make build

FROM gcr.io/distroless/static-debian13
WORKDIR /
COPY --from=build /app/butterclove /
EXPOSE 7590
CMD ["/butterclove"]
