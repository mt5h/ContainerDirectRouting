#!/usr/bin/env bash

baseurl='http://localhost'

for i in {1..5}; do 
    echo "testing: ${baseurl}/ with cookie instance=minion-$i"
    curl -L --header "Content-Type: application/json" --cookie "instance=minion-$i" $baseurl
    echo
done
    echo "Example wrong cookie"
    curl -L --header "Content-Type: application/json" --cookie "instance=minion-6" $baseurl
    echo    
    echo "Example no cookie"
    curl -L --header "Content-Type: application/json" $baseurl
    echo
