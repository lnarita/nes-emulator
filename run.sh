#!/usr/bin/env bash

cd game

[[ "$1" =~ "--nochecks" ]] && export NOCHECK=1

if [[ "$NOCHECK" != "1" ]]; then
    NESASM_EXEC=`command -v nesasm`
    if [[ -z $NESASM_EXEC ]]; then
        echo "Can't find nesasm (https://github.com/camsaul/nesasm) executable" >&2
        exit 1
    fi

    MEDNAFEN_EXEC=`command -v mednafen`
    if [[ -z $MEDNAFEN_EXEC ]]; then
        echo "Can't find mednafen installed in this machine!" >&2
        exit 1
    fi
fi

nesasm 2048.asm
mednafen 2048.nes
