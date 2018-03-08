FROM    vikings/golang:onbuild-v1.0.5
RUN     apk add --update make
COPY    goagent /goAgent
ENTRYPOINT ["/goAgent"]