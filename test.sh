#!/bin/sh
PORT=27960
RUN="go run ."

gnome-terminal --geometry=80x40+50+200 -e "bash -c 'hostname -I && $RUN $PORT';bash"
sleep 1
gnome-terminal --geometry=80x40+850+200 -e "bash -c 'nc localhost $PORT';bash"
sleep 1
gnome-terminal --geometry=80x40+1650+200 -e "bash -c 'nc localhost $PORT';bash"
sleep 10
gnome-terminal --geometry=80x40+850+600 -e "bash -c 'nc localhost $PORT';bash"