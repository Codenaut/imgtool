#!/bin/sh
#Make sure both imagemagick and ghostscript are installed

input=$1
output=$2
convert -verbose -density 600 -trim $input -quality 100 -sharpen 0x1.0 -background white -flatten $output
