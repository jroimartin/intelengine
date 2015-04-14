#!/bin/sh

ulimit -S -n 4096
exec rabbitmq-server $@
