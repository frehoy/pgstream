#!/usr/bin/env sh

export TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYXBpX3VzZXIifQ.W9Fo49rxMbSVnhdK1lzjMwCgf_1MZCPy9GNbt9j10ds"


echo "hello"
while true
do
	curl http://localhost:3000/events -X POST \
	     -H "Authorization: Bearer $TOKEN"   \
	     -H "Content-Type: application/json" \
	     -d '{"message": "learn how to auth"}'	
	curl http://localhost:3001/events -X POST \
	     -H "Authorization: Bearer $TOKEN"   \
	     -H "Content-Type: application/json" \
	     -d '{"message": "learn how to auth"}'	
done
echo "done"