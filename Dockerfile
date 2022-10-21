#Build go vendor
FROM golang:1.19 AS build-go
COPY . /w3bstream/

WORKDIR /w3bstream
RUN cd cmd/srv-applet-mgr && GOWORK=off GOOS=linux CGO_ENABLED=1 go build -mod vendor
RUN mkdir -p build

#Build noodjs
#FROM node:14.20-alpine AS build-nodejs
#FROM node:16-alpine AS build-nodejs
#
#WORKDIR /w3bstream-nodejs
#
#RUN apk add --no-cache curl
#RUN curl -fsSL "https://github.com/pnpm/pnpm/releases/latest/download/pnpm-linuxstatic-x64" -o /bin/pnpm
#RUN chmod +x /bin/pnpm
#COPY ./frontend .
#RUN npm i pnpm -g
#RUN pnpm i --frozen-lockfile;
#RUN pnpm build:standalone
FROM node:14.20-alpine AS build-nodejs

WORKDIR /w3bstream

COPY . /w3bstream/

RUN apk add --no-cache curl
RUN curl -fsSL "https://github.com/pnpm/pnpm/releases/latest/download/pnpm-linuxstatic-x64" -o /bin/pnpm
RUN chmod +x /bin/pnpm
RUN cd /w3bstream/frontend && pnpm install
RUN cp /w3bstream/frontend/.env.tmpl /w3bstream/frontend/.env
RUN cd /w3bstream/frontend && pnpm build


#run
FROM golang:1.19

EXPOSE 5432
EXPOSE 8888
EXPOSE 1883
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
#COPY . /w3bstream/
#RUN rm -rf /w3bstream/build_image/pgdata
WORKDIR /w3bstream

#RUN /etc/init.d/postgresql start
#Add VOLUMEs to allow backup of config, logs and databases
#VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql", "/var/lib/postgresql_data"]

#Initialize PostgreSQL database
RUN /etc/init.d/postgresql start && \
 su postgres sh -c "psql -c \"CREATE USER test_user WITH ENCRYPTED PASSWORD 'test_passwd'\"" && \
 su postgres sh -c "psql -c \"CREATE DATABASE test\"" && \
 su postgres sh -c "psql -c \"GRANT ALL PRIVILEGES ON DATABASE test to test_user;;\""
RUN mkdir -p /w3bstream/build_image/conf
COPY build_image/conf/postgresql.conf /w3bstream/build_image/conf/postgresql.conf

RUN echo "host all all 0.0.0.0/0 md5" >> /etc/postgresql/13/main/pg_hba.conf
#Install mqtt
#RUN apt-get install add-apt-repository
#RUN add-apt-repository ppa:mosquitto-dev/mosquitto-ppa
RUN apt install mosquitto mosquitto-clients -y
RUN mkdir /var/run/mosquitto -p
RUN rm -f /etc/mosquitto/mosquitto.conf
COPY build_image/conf/mosquitto.conf /etc/mosquitto/mosquitto.conf 
RUN /etc/init.d/mosquitto start

#Install Nginx
RUN apt-get install nginx -y

#Install Nodejs and pnpm
RUN mkdir -p /w3bstream/build_image/packages
COPY build_image/packages/setup_14.x /w3bstream/build_image/packages/setup_14.x
RUN /bin/bash /w3bstream/build_image/packages/setup_14.x
RUN apt install nodejs -y
RUN mkdir -p /w3bstream/build_image/packages
COPY build_image/packages/install.sh /w3bstream/build_image/packages/install.sh
RUN /bin/bash /w3bstream/build_image/packages/install.sh
RUN ln -s /root/.local/share/pnpm/pnpm /usr/bin/pnpm

#WORKDIR /w3bstream
RUN mkdir -p /w3bstream/cmd/srv-applet-mgr
COPY --from=build-go /w3bstream/cmd/srv-applet-mgr/srv-applet-mgr /w3bstream/cmd/srv-applet-mgr/srv-applet-mgr
COPY --from=build-go /w3bstream/cmd/srv-applet-mgr/config /w3bstream/cmd/srv-applet-mgr/config


COPY --from=build-nodejs /w3bstream/frontend /w3bstream/frontend
COPY --from=build-nodejs /w3bstream/frontend/.env.tmpl /w3bstream/frontend/.env

#COPY --from=builder-nodejs /w3bstream-nodejs/public ./frontend-build/public
#COPY --from=builder-nodejs --chown=nextjs:nodejs /w3bstream-nodejs/.next/standalone ./frontend-build
#COPY --from=builder-nodejs --chown=nextjs:nodejs /w3bstream-nodejs/.next/static ./frontend-build/.next/static

#Upload init script
ADD build_image/cmd/docker_init.sh /init.sh
RUN chmod 775 /init.sh
