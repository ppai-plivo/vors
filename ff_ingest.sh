#!/usr/bin/env bash

DIR=videos

reSeg='^[0-9]+$'
reIndex='.*m3u8$'

for d in `find $DIR -maxdepth 1 -type d`; do
	vname=$(basename -- "$d")
	for f in `ls $d | sort -h`; do
	  if [[ $f =~ $reIndex ]] ; then
	    cat $d/$f | rscat --set-key videos:$vname:$f
	  fi
	  if [[ $f =~ $reSeg ]] ; then
	    cat $d/$f | rscat --mode=produce --fmt=blob --stream=videos:$vname --id=$f-0 --silent --field-name=blob
	  fi
	done
done
