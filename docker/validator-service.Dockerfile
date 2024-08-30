FROM debian:12-slim
LABEL maintainer="Setlog <info@setlog.com>"
COPY out/validator .
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates
ENTRYPOINT ["./validator", "--act-as-service"]
