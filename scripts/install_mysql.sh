#!/bin/sh

set -ex

#Latest version and md5sum can be found at https://dev.mysql.com/downloads/repo/apt/
M_APT_VER="0.8.14-1"
sudo mv /var/lib/mysql /var/lib/mysql."$(date +%y%m%d%H%M%S)" || true
[ -f mysql-apt-config_${M_APT_VER}_all.deb ] || wget https://dev.mysql.com/get/mysql-apt-config_${M_APT_VER}_all.deb
echo "5cc94c7720fcd3124449b3e789441b98 mysql-apt-config_${M_APT_VER}_all.deb" | md5sum -c -
sudo MYSQL_SERVER_VERSION="mysql-8.0" DEBIAN_FRONTEND="noninteractive" /bin/sh -c "dpkg -i mysql-apt-config_${M_APT_VER}_all.deb && \
sudo -E apt-get update && \
sudo -E apt-get install -y --force-yes git make mysql-community-server"
sudo /usr/bin/mysqladmin shutdown || true
printf "[mysqld]\nserver-id=1\nbinlog-format=ROW\ngtid_mode=ON\nenforce-gtid-consistency\nlog_bin=/var/log/mysql/mysql-bin.log\nlog_slave_updates=1\ninnodb_flush_log_at_trx_commit=0\ninnodb_flush_log_at_timeout=30"|sudo tee -a /etc/mysql/my.cnf
sudo /usr/share/mysql-8.0/mysql-systemd-start pre
sudo /usr/bin/mysqld_safe --skip-syslog &
while ! sudo /usr/bin/mysqladmin ping; do sleep 1; done
sudo mysql -e "DROP USER IF EXISTS 'storagetapper';CREATE USER 'storagetapper' IDENTIFIED BY 'storagetapper';GRANT ALL ON *.* TO 'storagetapper'@'%';FLUSH PRIVILEGES"
if [ -n "$DOCKER_BUILD" ]; then
       sudo /usr/bin/mysqladmin shutdown || true
fi
