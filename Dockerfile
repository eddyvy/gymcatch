# Stage 1: Build the frontend
FROM node:22-alpine AS frontend-builder
WORKDIR /app/frontend
COPY . .
RUN npm install && npm run build

# Stage 2: Build the backend
FROM golang:1.22-alpine AS backend-builder
RUN apk update && apk add --no-cache ca-certificates
WORKDIR /app/backend
COPY . .
RUN go mod download
RUN go build -o gymcatch .

# Stage 3: Create the final image
FROM scratch
WORKDIR /app
COPY --from=backend-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=backend-builder /app/backend/gymcatch .
COPY --from=frontend-builder /app/frontend/dist ./dist
ENV PORT=3000
EXPOSE 3000
CMD ["./gymcatch"]
