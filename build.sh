#!/bin/bash

# 编译和交叉编译 Go 代码
build_and_package() {
    TARGET_OS=$1
    TARGET_ARCH=$2

    # 交叉编译 Go 代码
    GOOS=$TARGET_OS GOARCH=$TARGET_ARCH go build -o gocopy

    # 移动二进制文件到目标目录
    mv gocopy dist/

    # 进入目标目录
    cd dist

    # 将文件打包成 ZIP 文件
    zip gocopy-$TARGET_OS-$TARGET_ARCH.zip gocopy

    # 可选：删除原始二进制文件，只保留 ZIP 文件
    rm gocopy

    # 返回上级目录
    cd ..
}

# 调用函数，传递不同的平台参数
build_and_package linux amd64
build_and_package linux arm64
build_and_package windows amd64
build_and_package windows arm64
build_and_package darwin amd64
build_and_package darwin arm64
# 等等...
