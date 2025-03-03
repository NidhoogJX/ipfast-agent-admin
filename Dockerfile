# 编译项目为可执行二进制文件
FROM golang:1.22.12-alpine AS builder
ARG BUILDDIR
ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG EXPOSE_PORT=8098
WORKDIR /app
COPY . .
RUN go mod download
ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}
RUN go build -o /app/tocAgentAdmin ${BUILDDIR}
RUN mkdir -p /app/locales
RUN mkdir -p /app/static
COPY ${BUILDDIR}/config.yaml /app/
COPY ${BUILDDIR}/locales/ /app/locales/
COPY ${BUILDDIR}/static/ /app/static/

# 打包镜像
FROM alpine:3.21.2
WORKDIR /app
RUN apk add --no-cache tzdata
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone
RUN mkdir -p /app/log
# 从编译打包阶段镜像里复制编译好的二进制文件
COPY --from=builder /app/tocAgentAdmin /app/
# 从编译打包阶段镜像里复制配置文件,用于在运行容器时挂载配置文件
COPY --from=builder /app/config.yaml /app/
# 从编译打包阶段镜像里复制国际化语言文件夹,用于在运行容器时挂载locales文件夹
COPY --from=builder /app/locales/ /app/
# 从编译打包阶段镜像里复制静态文件夹,用于在运行容器时挂载static文件夹
COPY --from=builder /app/static/ /app/

EXPOSE ${EXPOSE_PORT}
CMD ["/app/tocAgentAdmin"]



