FROM alpine:3.2
ADD router-srv /router-srv
ENTRYPOINT [ "/router-srv" ]
