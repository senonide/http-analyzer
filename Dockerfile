FROM public.ecr.aws/docker/library/golang:1.23.1-alpine3.19 AS builder
WORKDIR /http-analyzer
COPY . .
RUN go mod download
RUN mkdir build
RUN go build -o build/http-analyzer .

FROM public.ecr.aws/docker/library/alpine:latest AS runner
WORKDIR /http-analyzer
COPY --from=builder /http-analyzer/build/http-analyzer /http-analyzer/build/http-analyzer
COPY --from=builder /http-analyzer/tls.crt /http-analyzer/tls.crt
COPY --from=builder /http-analyzer/tls.key /http-analyzer/tls.key
COPY --from=builder /http-analyzer/.env.example /http-analyzer/.env
EXPOSE 8443
CMD ["/http-analyzer/build/http-analyzer"]
