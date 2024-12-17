#!/bin/bash

# 设置变量
DOWNLOAD_URL="https://github.com/xibolun/jstat_exporter/releases/download/v1.0.0/jstat_exporter"
INSTALL_DIR="/opt/jstat_exporter"


# 创建安装目录
mkdir -p "$INSTALL_DIR"
cd "$INSTALL_DIR"

echo "Downloading ..."
curl -sSL -O "$DOWNLOAD_URL"

# 添加执行权限
chmod 770 jstat_exporter

# 获取用户输入
read -p "Enter Jstat path (default is /usr/bin/jstat): " jstat_path
read -p "Enter JMS exporter port(default is 9010): " port
read -p "Enter target java pid: " target_pid
echo

if [ -z "$jstat_path" ]; then
    jstat_path="/usr/bin/jstat"
fi

if [ -z "$port" ]; then
    port="9010"
fi

if [ -z "$target_pid" ]; then
    ehco "java pid can not be empty"
    exit(1)
fi


# 配置文件
cat << EOF > config.yml
jms_addr: "$jms_url"
jms_token: "$token"
interval: 10
dial_timeout: 3
http_port: $port
EOF

echo "Installation completed successfully."

if command -v "supervisorctl" &> /dev/null; then
cat << EOF > /etc/supervisord/jstat_exporter.conf
[program:jstat_exporter]
directory=$INSTALL_DIR
command=$INSTALL_DIR/jstat_exporter --web.listen-address=":$port" --jstat.path="$jstat_path" --target.pid="$target_pid"
autostart=true
autorestart=true
user=root
EOF
    supervisorctl update jstat_exporter
else
    $INSTALL_DIR/jstat_exporter --web.listen-address=":$port" --jstat.path="$jstat_path" --target.pid="$target_pid"
fi