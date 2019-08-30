# nes-emulator
2s2019 MC861/MC871 A

by november 2019 this repository will hopefully be home to a working nes emulator
and a 2048 clone game for nes

# Assembling and playing the game

This project was made using NESASM 3 by bunnyboy, so it needs to be installed in your machine if you want to assembly it yourself.

The `run.sh` script must be run in the project root folder (as it will try to `cd game`).

If you have `nesasm` and `mednafen` commands in your machine but the script keeps complaining, just run it with the `--nochecks` flag as shown below.

```
./run.sh --nochecks
```

The script can also assemble the game before running it by using the `--assemble` flag.

# Collaborators

- 170600 - Erick Seiji Furukawa
- 169176 - Guilherme Tiaki Sassai Sato
- 176353 - João Vitor Araki Gonçalves
- 182851 - Lucy Narita
