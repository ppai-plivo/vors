#!/usr/bin/env bash

DIR=videos
reVideos='.*.mp4$'

for f in `ls $DIR`;do
  if ! [[ $f =~ $reVideos ]] ; then
	  continue
  fi
  vidname="${f%.*}"
  mkdir -p $DIR/$vidname
  ffmpeg -i $DIR/$f -profile:v baseline -level 3.0 -s 640x360 -start_number 0 -hls_time 5 -hls_list_size 0 -start_number 1 -hls_segment_filename "$DIR/$vidname/%d" -f hls "$DIR/$vidname/index.m3u8"
done
