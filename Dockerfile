FROM    vikings/alpine:base
LABEL 	maintainer=ztao8607@gmail.com
RUN 	cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    	echo "Asia/Shanghai" > /etc/timezone
COPY    server.key /server.key
COPY    server.crt /server.crt
COPY 	bin/hicd /hicd
EXPOSE 	443
ENTRYPOINT ["/hicd"]