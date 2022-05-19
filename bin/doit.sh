#!/usr/bin/env sh

wrk --connections 16 --duration 600 --threads 8 -s wrk.lua http://127.0.0.1:3000/events
