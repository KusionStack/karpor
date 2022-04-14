FROM reg.docker.alibaba-inc.com/alipay/7u2-common:202107.0T

COPY build/linux/ /work/
RUN chmod +x /work/iactestpolicy

ENV PATH="/work/:${PATH}"
ENV LANG=en_US.utf8
ENV GIN_MODE=release
ENV PORT=8080

WORKDIR /work
EXPOSE 8080

LABEL maintainer="yingming.yym@antgroup.com"

ENTRYPOINT ["./iactestpolicy"]
