FROM alpine:latest

COPY user-web /home
RUN chmod +x /bin/grpc_health_probe
STOPSIGNAL SIGTERM

# CMD ["nginx", "-g", "daemon off;"]
WORKDIR /home/
CMD ./user-web