FROM reg.docker.alibaba-inc.com/alipay/7u2-common:202202.0T

COPY build/linux/ /app/
RUN chmod +x /app/karbour

ENV PATH="/app/:${PATH}"
ENV LANG=en_US.utf8
ENV GIN_MODE=release
ENV PORT=8080

WORKDIR /app
EXPOSE 8080

LABEL maintainer="yingming.yym@antgroup.com"

ENTRYPOINT ["./karbour"]
