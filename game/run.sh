#!/usr/bin/env bash

NOCHECKSFLAG="--nochecks"
ASSEMBLEFLAG="--assemble"

[[ "$1" =~ "$NOCHECKSFLAG" ]] && export NOCHECK=1
[[ "$2" =~ "$NOCHECKSFLAG" ]] && export NOCHECK=1
[[ "$1" =~ "$ASSEMBLEFLAG" ]] && export ASSEMBLE=1
[[ "$2" =~ "$ASSEMBLEFLAG" ]] && export ASSEMBLE=1

if [[ "$NOCHECK" != "1" ]] && [[ "$ASSEMBLE" == "1" ]]; then
    NESASM_EXEC=`command -v nesasm`
    if [[ -z $NESASM_EXEC ]]; then
        echo "Can't find nesasm (https://github.com/camsaul/nesasm) executable" >&2
        exit 1
    fi
fi

[[ "$ASSEMBLE" == "1" ]] && nesasm 2048.asm

if [[ "$NOCHECK" != "1" ]]; then
    MEDNAFEN_EXEC=`command -v mednafen`
    if [[ -z $MEDNAFEN_EXEC ]]; then
        echo "Can't find mednafen installed in this machine!" >&2
        exit 1
    fi
fi

mednafen 2048.nes
