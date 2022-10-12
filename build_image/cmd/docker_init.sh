#!/bin/sh
#init pg
if [ ! -d "/var/lib/postgresql_data/13" ]; then
   echo "PG data does not exsit!"
   cp -r /var/lib/postgresql/13 /var/lib/postgresql_data/
else
   echo "PG data exists!"
fi
chown -R postgres:postgres /var/lib/postgresql_data/13
chmod -R 700 /var/lib/postgresql_data/13
#Start postgres
su postgres -c "/usr/lib/postgresql/13/bin/postgres -D /var/lib/postgresql_data/13/main -c config_file=/etc/postgresql/13/main/postgresql.conf"&
#su postgres -c "/usr/lib/postgresql/13/bin/postgres -D /var/lib/postgresql_data/13/main -c config_file=/etc/postgresql/13/main/postgresql.conf"&
su postgres sh -c "createuser test_user"
if [ $? -ne 0 ];then
   sleep 15
   su postgres sh -c "createuser test_user"
fi
su postgres sh -c "psql -c \"ALTER USER test_user PASSWORD 'test_passwd'\""
#su postgres sh -c "psql -c \"CREATE USER test_user WITH ENCRYPTED PASSWORD 'test_passwd'\""
su postgres sh -c "psql -c \"CREATE DATABASE test\""
su postgres sh -c "psql -c \"GRANT ALL PRIVILEGES ON DATABASE test to test_user;;\""

#Start mosquitto
/etc/init.d/mosquitto start
#sleep 10

cd /w3bstream/build && ./srv-applet-mgr migrate
#cd /w3bstream/build && ./srv-applet-mgr init_admin &
cd /w3bstream/build && ./srv-applet-mgr &

sleep 3

cd /w3bstream/frontend && pnpm start
