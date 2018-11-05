#!/bin/sh

while true; do
  dig "$@"
  sleep 5
done
