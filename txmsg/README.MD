下载 protobuf https://github.com/protocolbuffers/protobuf/releases/download/v3.7.0/protobuf-all-3.7.0.zip <p>

跳转到解压后的目录

设置编译目录 ./configure --prefix=/usr/local/protobuf #/usr/local/protobuf/ 为自己配置的编译安装目录

在解压的目录下进行 make sudo make install

配置环境变量 vim ~/.bash_profile export PROTOBUF=/usr/local/protobuf export PATH=$PROTOBUF/bin:$PATH

测试安装结果 输入protoc --version 输出libprotoc 3.7.0

go get github.com/golang/protobuf/proto

protoc --proto_path=./ --go_out=plugins=grpc:./ ./RpcMessage.proto