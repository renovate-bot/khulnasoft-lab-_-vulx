FROM alpine:3.11
RUN apk --no-cache add ca-certificates git rpm
COPY vul /usr/local/bin/vul
COPY contrib/gitlab.tpl contrib/gitlab.tpl
ENTRYPOINT ["vul"]
