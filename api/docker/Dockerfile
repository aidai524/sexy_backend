FROM golang:1.23.3
ENV TZ=UTC
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 在docker的根目录下创建相应的使用目录
WORKDIR "/app"

# 把当前（宿主机上）目录下的文件都复制到docker上刚创建的目录下
COPY /main /app/
COPY /config.toml /app/

EXPOSE 8080

