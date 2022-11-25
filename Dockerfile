FROM alpine:latest

COPY user-web /home
STOPSIGNAL SIGTERM

# CMD ["nginx", "-g", "daemon off;"]
WORKDIR /home/
CMD ./user-web