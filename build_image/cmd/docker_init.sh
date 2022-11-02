#!/bin/sh
#init pg

pg_config="/w3bstream/build_image/config/postgresql.conf"
pg_data="/var/lib/postgresql_data"
if [ -f "${pg_config}" ]; then
   rm -f /etc/postgresql/13/main/postgresql.conf
   ln -s ${pg_config} /etc/postgresql/13/main/postgresql.conf 
   if [ ! -d "${pg_data}/13" ]; then
      echo "PG data does not exsit!"
      cp -r /var/lib/postgresql/13 ${pg_data}
      rm -rf /var/lib/postgresql/13
      chown -R postgres:postgres ${pg_data}/13
      chmod -R 700 ${pg_data}/13
   else
      echo "PG data exists!"
   fi
   su postgres -c "/usr/lib/postgresql/13/bin/postgres -D /var/lib/postgresql_data/13/main -c config_file=/etc/postgresql/13/main/postgresql.conf"&
else
   su postgres -c "/usr/lib/postgresql/13/bin/postgres -D /var/lib/postgresql/13/main -c config_file=/etc/postgresql/13/main/postgresql.conf"&
fi


#Start postgres
#su postgres -c "/usr/lib/postgresql/13/bin/postgres -D /var/lib/postgresql_data/13/main -c config_file=/etc/postgresql/13/main/postgresql.conf"&
#su postgres sh -c "createuser test_user"
#if [ $? -ne 0 ];then
#   sleep 15
#   su postgres sh -c "createuser test_user"
#fi
#su postgres sh -c "psql -c \"ALTER USER test_user PASSWORD 'test_passwd'\""
##su postgres sh -c "psql -c \"CREATE USER test_user WITH ENCRYPTED PASSWORD 'test_passwd'\""
#su postgres sh -c "psql -c \"CREATE DATABASE test\""
#su postgres sh -c "psql -c \"GRANT ALL PRIVILEGES ON DATABASE test to test_user;;\""

#Start mosquitto
mqtt_config="/w3bstream/build_image/config/mosquitto.conf"
if [ -f "${mqtt_config}" ]; then
   rm -f /etc/mosquitto/mosquitto.conf
   ln -s ${mqtt_config} /etc/mosquitto/mosquitto.conf
fi

/etc/init.d/mosquitto start
#sleep 10

cd /w3bstream/cmd/srv-applet-mgr && ./srv-applet-mgr migrate
#cd /w3bstream/build && ./srv-applet-mgr init_admin &
cd /w3bstream/cmd/srv-applet-mgr && ./srv-applet-mgr &
sleep 3

cd /w3bstream/frontend-build && node server.js

