#!/usr/bin/env bash

DIR=videos
MNTPT='/tmp/mntpt'
reVideos='.*.mp4$'

for f in `ls $DIR`;do
  if ! [[ $f =~ $reVideos ]] ; then
	  continue
  fi
  vidname="${f%.*}"
  mkdir -p $MNTPT/$DIR:$vidname
  ls -l $MNTPT
  ffmpeg -i $DIR/$f -profile:v baseline -level 3.0 -s 640x360 -start_number 0 -hls_time 5 -hls_list_size 0 -start_number 1 -hls_segment_filename "$MNTPT/$DIR:$vidname/%d" -f hls "$MNTPT/$DIR:$vidname:index.m3u8"
done
