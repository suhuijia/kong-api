# launch KGK for Handwriting Mathematical Expression Recognition
#!/usr/bin/env bash

if [ ! -d "log" ]; then
    mkdir "log"
fi

nohup ./KGK_HWER 8004 &
