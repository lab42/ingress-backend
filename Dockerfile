# Setup
FROM alpine:3.21 AS setup
RUN addgroup --gid 10000 -S appgroup && \
    adduser --uid 10000 -S appuser -G appgroup

FROM scratch
COPY --from=setup /etc/passwd /etc/passwd
COPY ingress-backend /ingress-backend
USER appuser
EXPOSE 1234
ENTRYPOINT ["/ingress-backend"]
