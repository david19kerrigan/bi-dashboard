#!/bin/sh

sudo mkdir /run/postgresql
sudo chown david /run/postgresql
pg_ctl -D ~/Documents/bi-dashboard/pg start
