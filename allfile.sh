#!/bin/bash

for file in */scripts/*.sh
do
  #cmd [option] "$file" >> results.out
  #sed -i -e 's/\r$//' "$file"
  echo $file
  chmod +x $file
done
