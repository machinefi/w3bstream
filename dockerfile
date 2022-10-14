#base image 
FROM golang:1.19 AS build-env

MAINTAINER The Iotex Project <jeson@iotex.io>

EXPOSE 5432
EXPOSE 8888
EXPOSE 3000 
EXPOSE 22
#Create user postgres
RUN mkdir -p /home/postgres
RUN groupadd postgres --gid=999 \
  && useradd -d /home/postgres --gid postgres --uid=999 postgres
RUN chown -R postgres:postgres /home/postgres

#Install postgresql and postgis
RUN apt -y update
RUN apt install postgresql postgresql-contrib -y
RUN apt-get install net-tools -y
RUN apt install postgis postgresql-13-postgis-3 -y
#RUN apt install -y postgresql-13 postgresql-client-13


#Create app directory and upload w3bstream app
RUN mkdir -p /w3bstream
COPY . /w3bstream/
RUN rm -rf /w3bstream/build_image/pgdata
WORKDIR /w3bstream

#RUN /etc/init.d/postgresql start
#Add VOLUMEs to allow backup of config, logs and databases
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql", "/var/lib/postgresql_data"]

#Initialize PostgreSQL database
RUN /etc/init.d/postgresql start && \
 su postgres sh -c "psql -c \"CREATE USER test_user WITH ENCRYPTED PASSWORD 'test_passwd'\"" && \
 su postgres sh -c "psql -c \"CREATE DATABASE test\"" && \
 su postgres sh -c "psql -c \"GRANT ALL PRIVILEGES ON DATABASE test to test_user;;\""

#Install mqtt
#RUN apt-get install add-apt-repository
#RUN add-apt-repository ppa:mosquitto-dev/mosquitto-ppa
RUN apt install mosquitto mosquitto-clients -y
RUN mkdir /var/run/mosquitto -p
RUN rm -f /etc/mosquitto/mosquitto.conf && ln -s /w3bstream/build_image/conf/mosquitto.conf /etc/mosquitto/mosquitto.conf 
RUN /etc/init.d/mosquitto start

#Install Nginx
RUN apt-get install nginx -y

#Install Nodejs and pnpm
RUN /bin/bash /w3bstream/build_image/packages/setup_14.x
RUN apt install nodejs -y
RUN /bin/bash /w3bstream/build_image/packages/install.sh
RUN ln -s /root/.local/share/pnpm/pnpm /usr/bin/pnpm

#Build vendor
WORKDIR /w3bstream
RUN cd cmd/srv-applet-mgr && go build -mod vendor
RUN mkdir -p build
RUN mv cmd/srv-applet-mgr/srv-applet-mgr build
RUN cp -r cmd/srv-applet-mgr/config build/config
RUN echo 'succeed! srv-applet-mgr =>build/srv-applet-mgr*'
RUN echo 'succeed! config =>build/config/'
RUN echo 'modify config/local.yaml to use your server config'

#Build w3bstream front
RUN cd /w3bstream/frontend && pnpm install
RUN cp /w3bstream/frontend/.env.tmpl /w3bstream/frontend/.env
RUN cd /w3bstream/frontend && pnpm build

#
#Upload init script
ADD build_image/cmd/docker_init.sh /init.sh
RUN chmod 775 /init.sh
