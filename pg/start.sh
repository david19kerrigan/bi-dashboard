#!/bin/sh

sudo mkdir /run/postgresql
sudo chown david /run/postgresql
pg_ctl -D . start
