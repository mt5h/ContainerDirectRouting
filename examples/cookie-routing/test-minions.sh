#!/usr/bin/env bash

baseurl='https://localhost'
max=2

for i in $(seq 1 $max); do 
    echo "testing: ${baseurl}/ with cookie instance=minion-$i"
    curl -L --header "Content-Type: application/json" --insecure \
      --cookie "instance=minion-$i" \
      $baseurl
    echo
done
    echo "Example wrong cookie"
    curl -L --insecure --header "Content-Type: application/json" --cookie "instance=minion-6" $baseurl
    echo    
    echo "Example no cookie"
    curl -L --insecure --header "Content-Type: application/json" $baseurl
    echo
