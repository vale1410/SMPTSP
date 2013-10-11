#!/bin/zsh 

mkdir -p tmp
mkdir -p sol

#Runtime=${2-10} # runtime in seconds
#Strategy=${3-1}
#
#echo runtime $Runtime
#echo strategy $Strategy

tmp=tmp/$(basename $1 .dat).tmp

echo tmp $tmp

go run encode.go -f $1 -o $tmp

#Option=' --stat --time-limit='$Runtime' tmp.dat '
##Option=$Option' --time-limit='$Runtime' --restart-on-model --opt-hierarch=2 ' 
##Option=$Option' --time-limit='$Runtime' -t 2,compete --restart-on-model --opt-hierarch=2 '
#
#
#case $Strategy in
#    0) Option=$Option'--configuration=frumpy ';;
#    1) Option=$Option'--configuration=jumpy ';;
#    2) Option=$Option'--configuration=handy ';;
#    3) Option=$Option'--configuration=crafty ';;
#    4) Option=$Option'--configuration=trendy ';;
#    5) Option=$Option'--configuration=chatty ';;
#esac
#
#echo $Option
#

sol=sol/$(basename $1 .dat).sol

gringo independent_set.lp $tmp | clasp --stat --time-limit=20 --configuration=chatty -t 3 > $sol
