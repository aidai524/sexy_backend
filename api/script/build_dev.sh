#!/usr/bin/env bash

#使用方法
#cd /Users/maoli/delta-bot-backend/api/script
#sudo ./build.sh

# 从脚本位置获取到项目根目录的相对路径
# 例如，如果脚本在项目根目录下的 "scripts" 文件夹中，
# 根据实际情况调整这个路径
project_root_dir="../../"
toml_path="/root/service/dev/sexy_backend/config/dev/"

# 定义应用名称
app_name='api_dev'
# 定义应用版本
app_version='v1'
# 设置编译后的二进制文件存放的目录和文件名
output_dir="bin"
binary_name="main"

# 使用Docker运行Go容器来构建应用
# 注意调整-v参数，确保挂载的是项目根目录到容器的/app目录
echo '----go build----'

docker run --rm \
    -v "$(pwd)/${project_root_dir}":/app \
    -v go-mod-cache:/go/pkg/mod \
    -w /app \
    golang:1.23.3 go build -o ${output_dir}/${binary_name} ./api/cmd/main/main.go

# 复制docker相关文件到bin目录
cp -r ../docker/Dockerfile_dev ${project_root_dir}bin/Dockerfile
cp -r ${toml_path}api/config.toml ${project_root_dir}bin/


# 切换到包含Dockerfile的目录
cd ${project_root_dir}${output_dir}/

# 提示用户输入解密密钥，并记录下来
echo "请输入解密密钥(key):"
read -s decryption_key

echo '----copy ----'
docker stop ${app_name}
echo '----stop container----'
docker rm ${app_name}
echo '----rm container----'
docker rmi ${app_name}:${app_version}
echo '----rm image----'
mkdir -p ${output_dir}

# 构建Docker镜像
docker build -t ${app_name}:${app_version} .

# （可选）运行Docker镜像
docker run -p 8081:8081 --name ${app_name} -d --restart=on-failure --log-opt max-size=5m -m 250m ${app_name}:${app_version} /app/main -conf=/app/config.toml -key="$decryption_key"