#!/bin/sh

case $1 in
"identity-api") echo 1001
;;
"identity-rpc") echo 1101
;;
"usercenter-api") echo 1002
;;
"usercenter-rpc") echo 1102
;;
"message-mq") echo 1207
;;
"mqueue-rpc") echo 1106
;;
"order-api") echo 1004
;;
"order-mq") echo 1204
;;
"order-rpc") echo 1104
;;
"payment-api") echo 1005
;;
"payment-rpc") echo 1105
;;
"travel-api") echo 1003
;;
"travel-rpc") echo 1103
esac
