#!/usr/bin/env sh


# yes | xargs -P 16 bin/post_stuff.sh

export TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYXBpX3VzZXIifQ.W9Fo49rxMbSVnhdK1lzjMwCgf_1MZCPy9GNbt9j10ds"

ab \
	-H "Authorization: Bearer $TOKEN" \
	-p bin/data.json \
	-T application/json \
	-n 17000 \
	-c 2 \
	-s 30 \
	http://127.0.0.1:3000/events
