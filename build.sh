# 根目录下执行以下命令
IPFAST_TAG=v1.0.0

docker build -f Dockerfile -t ipfast . --build-arg BUILDDIR=./cmd/api
docker tag ipfast:latest dahangkk/ipfast:$IPFAST_TAG
docker push dahangkk/ipfast:$IPFAST_TAG

# 运行计划服务容器
docker run -d \
  -v /Users/jinxin/Desktop/ipfast/cmd/api/config.yaml:/app/config.yaml \
  -v /Users/jinxin/Desktop/ipfast/cmd/api/locales:/app/locales \
  -v /Users/jinxin/Desktop/ipfast/cmd/api/static:/app/static \
  --name ipfast \
  --restart always \
  -p 40003:40003 \
  dahangkk/ipfast


# # 清理无用Docker资源 停止的容器、未被使用的镜像、未被使用的网络、未被挂载的数据卷
# docker system prune -a
# # go 自带的性能分析工具 网页地址
# http://localhost:6060/debug/pprof/
# # 内存分析
# go tool pprof http://localhost:6060/debug/pprof/heap
# # 分析导出为 pdf
# go tool pprof -pdf http://localhost:6060/debug/pprof/heap > heap.pdf
# # 分析导出为可视化UI
# go tool pprof -http=:8088 /Users/jinxin/pprof/pprof.balance.alloc_objects.alloc_space.inuse_objects.inuse_space.003.pb.gz

# echo 'export http_proxy=http://127.0.0.1:7890' >> ~/.zshrc
# echo 'export https_proxy=http://127.0.0.1:7890' >> ~/.zshrc
# source ~/.zshrc

# curl -I http://www.google.com

#GOOS=linux GOARCH=amd64 go build -o balance

# # 2. 进入 Redis 容器（假设容器名称为 redis_container）
# docker exec -it ad_redis sh

# # 3. 连接到 Redis 实例
# redis-cli

# # 4. 使用 AUTH 命令进行身份验证（假设密码为 yourpassword）
# AUTH myredisdata

# # 5. 查看指定 key 的数据，例如 1_AD_COUNT
# GET 1_AD_COUNT

# TTL 1_AD_COUNT
# KEYS '*'

# 使用 SCAN 命令过滤和限制数量（假设匹配模式为 '1_*'，每次返回 10 个键）
# SCAN 0 MATCH 1_* COUNT 10

# uname -m
# 如果输出是 x86_64，则使用 amd64。
# 如果输出是 armv7l，则使用 arm。
# 如果输出是 aarch64，则使用 arm64。