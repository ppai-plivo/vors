#!/usr/bin/env bash

DIR=redis

reSeg='^[0-9]+$'
reIndex='.*m3u8$'

for f in `ls $DIR | sort -h`; do
  if [[ $f =~ $reIndex ]] ; then
	cat $DIR/$f | rscat --set-key videos:$DIR:$f
  fi
  if [[ $f =~ $reSeg ]] ; then
    cat $DIR/$f | rscat --mode=produce --fmt=blob --stream=videos:$DIR --id=$f-0 --silent --field-name=blob
  fi
done
