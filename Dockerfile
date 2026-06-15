FROM node:26-alpine AS frontend
WORKDIR /build
COPY frontend/ .
RUN npm ci && npm run build

FROM golang:1.26-alpine AS builder
WORKDIR /build
COPY api/ .
COPY --from=frontend /build/dist/ internal/webui/dist/
RUN go build -o /mailfyke .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /mailfyke /mailfyke
EXPOSE 5789 2525
ENTRYPOINT ["/mailfyke"]
