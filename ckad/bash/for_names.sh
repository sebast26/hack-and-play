#!/bin/bash

names=$(ls -1)

for i in $names; do echo "$i has $(echo -n $i | wc -c)"; done 
