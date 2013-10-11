#!/bin/zsh

gringo prototype/combined.lp instances/Data_010/data_1_50_258_20_90_100_200_300.lp | clasp --stat --configuration=chatty -t 6 > output.txt
