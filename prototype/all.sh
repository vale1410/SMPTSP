#!/bin/zsh

rm -fr sol

for x in Data_010/*; do ./run2.sh $x; done
mv sol Ind_010
for x in Data_137/*; do ./run2.sh $x; done 
mv sol Ind_137
for x in Data_100/*; do ./run2.sh $x; done 
mv sol Ind_100

for x in Data_010/*; do ./run.sh $x; done
mv sol Sol_010
for x in Data_137/*; do ./run.sh $x; done 
mv sol Sol_137
for x in Data_100/*; do ./run.sh $x; done 
mv sol Sol_100
