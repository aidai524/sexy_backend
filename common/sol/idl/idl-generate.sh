#!/usr/bin/env bash

# 进入包含IDL文件的目录
cd $(dirname $0)/../idl

# 遍历所有IDL文件
for filename in *.json; do
  name=$(basename $filename .json)
  echo "Generating Go client for $name"

  # 调用anchor-go工具生成Go代码
  anchor-go --src $filename --dst ./$name
done

echo "All Go clients have been generated."