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
    exit 1
fi


# 配置文件
cat << EOF > start.sh
#!/bin/bash
target_pid=$target_pid
jstat_path=$jstat_path
/opt/jstat_exporter/jstat_exporter --web.listen-address=":$port" --jstat.path=\$jstat_path --target.pid=\$target_pid
EOF

chmod 775 $INSTALL_DIR/start.sh

echo "Installation completed successfully."

if command -v "supervisorctl" &> /dev/null; then
cat << EOF > /etc/supervisord/jstat_exporter.conf
[program:jstat_exporter]
directory=$INSTALL_DIR
command=/bin/bash -i -c $INSTALL_DIR/start.sh
autostart=true
autorestart=true
user=root
EOF
    supervisorctl update jstat_exporter
else
    bash -c $INSTALL_DIR/start.sh
fi