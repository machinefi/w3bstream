#!/bin/sh

/etc/init.d/postgresql start
sleep 10
/etc/init.d/mosquitto start
sleep 10

cd /w3bstream/build && ./srv-applet-mgr migrate
cd /w3bstream/build && ./srv-applet-mgr init_admin &>/w3bdata/password.txt &
cd /w3bstream/build && ./srv-applet-mgr &

cd /w3bstream/studio && pnpm start
