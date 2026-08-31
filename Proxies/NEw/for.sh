#!/bin/bash
sum=0
for i in {1..10}
do
    sum=$((sum+i))
done
echo "SUM = $sum"