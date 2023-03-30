FROM nats:scratch
ADD build /client
HEALTHCHECK CMD client test
