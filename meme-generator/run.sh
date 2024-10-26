#!/bin/bash

while true; do
    if pg_isready; then
        break
    fi
    sleep 1
done

exec /app/meme-generator.bin
