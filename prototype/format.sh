#!/bin/zsh 

for x in $1/*
do 
    name=$(basename $x .sol)
    a=$(cat $x | grep total | sed 's/total(//g' | sed 's/)//g' | tail -n 1)
    b=$(cat $x | grep Optimization | sed 's/Optimization: //g' | tail -n 1)
    t=$(cat $x | grep '^Time' | sed 's/^Time *: //g' | sed 's/(.*//g')
    echo $name $a $b $t
done
