#!/usr/bin/env bash

baseurl='http://localhost'

for i in {1..5}; do 
    echo "testing: ${baseurl}/minion-$i/"

    # wake up the instance
    response=$(curl -s --header "Content-Type: application/json" "${baseurl}/minion-$i")
    # unable to follow redirect and use the provided cookie atm
    # test installed route
    redirect=$(curl -s --cookie "instance=minion-$i" "${baseurl}/")


  if command -v jq > /dev/null 2>&1; then
    echo $response 
    echo $redirect | jq
  else
    echo "Install jq for a better output"
    echo $response 
    echo $redirect
  fi
done
