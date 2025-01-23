FROM node:18 as css-builder

WORKDIR /build
COPY tailwind/package*.json ./
COPY tailwind/tailwind.config.js ./
COPY tailwind/base.css ./
COPY ./assets ./assets

RUN npm install
RUN npm run build:css

FROM golang:1.23-alpine AS builder

ENV ROOT=/go/src/app
WORKDIR ${ROOT}

# Install required packages
RUN apk add --no-cache gcc musl-dev sqlite-dev sqlite-libs ca-certificates

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate templ files
RUN templ generate

# Copy CSS from first stage
RUN mkdir -p public/css
COPY --from=css-builder /build/public/css/output.css ./public/css/

# Setup user and permissions
RUN addgroup -S appuser && adduser -S -G appuser appuser
RUN mkdir -p /go/src/app/data && \
    chown -R appuser:appuser /go/src/app/data && \
    chmod 750 /go/src/app/data

# Build the application
RUN go build -o main ./cmd

# Setup database
RUN touch /go/src/app/data/jobs.db && \
    chown appuser:appuser /go/src/app/data/jobs.db && \
    chmod 640 /go/src/app/data/jobs.db

# Install air for hot reload
RUN go install github.com/air-verse/air@latest

EXPOSE 8082

USER appuser

COPY .air.toml ./
CMD ["air", "-c", ".air.toml"]
