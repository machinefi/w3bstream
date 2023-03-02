#!/bin/sh

passwd=password
rootsubj="/C=CN/ST=GuangDong/L=ShenZhen/O=Company/OU=Gateway/CN=Root CA/"
cakey=ca.key
cacrt=ca.crt

echo "generating root certificates"

openssl genrsa -des3 -passout pass:$passwd -out $cakey 2048 > /dev/null 2>&1
openssl req -new -x509 -days 3650 -key ca.key -out $cacrt -subj "$rootsubj" -passin pass:$passwd > /dev/null 2>&1

echo "\tsubject:  $rootsubj"
echo "\tpassword: $passwd"
echo "\toutputs:  $cakey $cacrt"

echo "generating server certificates"
sub="/C=CN/ST=GuangDong/L=ShenZhen/O=Server/OU=Gateway/CN=localhost"
key=server.key
csr=server.csr
crt=server.crt
pem=server.pem
openssl genrsa -out $key 2048 > /dev/null 2>&1
openssl req -new -out $csr -key $key -subj "$sub" > /dev/null 2>&1
openssl x509 -req -in $csr -CA $cacrt -CAkey $cakey -CAcreateserial -out $crt -days 3650 -passin pass:$passwd > /dev/null 2>&1
openssl x509 -in $crt -out $pem

echo "\tsubject:  $sub"
echo "\toutputs:  $key $csr $crt $pem"

echo "generating client certificates"
sub="/C=CN/ST=GuangDong/L=ShenZhen/O=Server/OU=Gateway/CN=client cert"
key=client.key
csr=client.csr
crt=client.crt
pem=client.pem
openssl genrsa -out $key 2048 > /dev/null 2>&1
openssl req -new -out $csr -key $key -subj "$sub" > /dev/null 2>&1
openssl x509 -req -in $csr -CA $cacrt -CAkey $cakey -CAcreateserial -out $crt -days 3650 -passin pass:$passwd > /dev/null 2>&1
openssl x509 -in $crt -out $pem

echo "\tsubject:  $sub"
echo "\toutputs:  $key $csr $crt $pem"
