#!/bin/ash

echo "FLG{$(echo $RANDOM | md5sum | head -c 20; echo)}" >> flag.txt
python3 -m http.server 80
