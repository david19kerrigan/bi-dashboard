#!/bin/sh

~/Documents/bi-dashboard/pg/start.sh
tmux send-keys 'cd ~/Documents/bi-dashboard/js ; npm run dev' ENTER
tmux split-window -h ; tmux send-keys 'cd ~/Documents/bi-dashboard/go ; air' ENTER
