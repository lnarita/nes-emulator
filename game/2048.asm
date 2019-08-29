;--------------------------------------------------------------------
; iNES Header
;--------------------------------------------------------------------
    .inesprg 1   ; 1x 16KB PRG code
    .ineschr 1   ; 1x  8KB CHR data
    .inesmap 0   ; mapper 0 = NROM, no bank swapping
    .inesmir 1   ; background mirroring

;--------------------------------------------------------------------
; variables
;--------------------------------------------------------------------
    .rsset $0000   ; start variables at ram location 0

gamestate   .rs 1  ; .rs 1 means reserve one byte of space
pointerLo   .rs 1  ; pointer variables are declared in RAM
pointerHi   .rs 1  ; low byte first, high byte immediately after
buttons1    .rs 1  ; player 1 gamepad buttons, one bit per button
bgTileLo    .rs 1
bgTileHi    .rs 1
lastPressed .rs 1
tiles       .rs 16
currTile    .rs 1
random      .rs 1
soundTimer  .rs 1
beginTile   .rs 1
tempWord    .rs 2
scoreLo     .rs 1
scoreHi     .rs 1
scoreDig    .rs 5
scoreTempHi .rs 1
scoreTempLo .rs 1
score2Lo    .rs 1
score2Hi    .rs 1

;--------------------------------------------------------------------
; constants
;--------------------------------------------------------------------

STATETITLE     = $00  ; displaying title screen
STATEPLAYING   = $01  ; move paddles/ball, check for collisions
STATEGAMEOVER  = $02  ; displaying game over screen

;; background tiles
GUL = $0e ; playing tile up left corner
GUR = $0f ; playing tile up right corner
GDL = $10 ; playing tile down left corner
GDR = $11 ; playing tile down right corner
GLB = $0a ; playing tile left border
GRB = $0b ; playing tile right border
GUB = $0c ; playing tile up border
GDB = $0d ; playing tile down border
GBL = $26 ; playing tile blank fill
GBG = $24 ; screen solid background
GLG = $12 ; outer left border
GRG = $13 ; outer right border
GUG = $14 ; outer up border
GDG = $15 ; outer down border
; $e0, $da, $e6, $de, GBG, GBG, $e8, $ef, $de, $eb,
;; Char tiles
DC0 = $d0
DC1 = $d1
DC2 = $d2
DC3 = $d3
DC4 = $d4
DC5 = $d5
DC6 = $d6
DC7 = $d7
DC8 = $d8
DC9 = $d9
DCA = $da
DCB = $db
DCC = $dc
DCD = $dd
DCE = $de
DCF = $df
DCG = $e0
DCH = $e1
DCI = $e2
DCJ = $e3
DCK = $e4
DCL = $e5
DCM = $e6
DCN = $e7
DCO = $e8
DCP = $e9
DCQ = $ea
DCR = $eb
DCS = $ec
DCT = $ed
DCU = $ee
DCV = $ef
DCW = $f0
DCX = $f1
DCY = $f2
DCZ = $f3

;; buttons
GAMEPAD_A      = %10000000
GAMEPAD_B      = %01000000
GAMEPAD_SELECT = %00100000
GAMEPAD_START  = %00010000
GAMEPAD_UP     = %00001000
GAMEPAD_DOWN   = %00000100
GAMEPAD_LEFT   = %00000010
GAMEPAD_RIGHT  = %00000001
GAMEPAD_ANY_PRESSED = %11111111

;;sound
BEEP_DURATION = $03

;--------------------------------------------------------------------
; program bank(s)
;--------------------------------------------------------------------

    .bank 0
    .org $C000

vblankwait:      ; wait for vblank
    BIT $2002
    BPL vblankwait
    RTS

RESET:
    SEI          ; disable IRQs
    CLD          ; disable decimal mode
    LDX #$40
    STX $4017    ; disable APU frame IRQ
    LDX #$FF
    TXS          ; Set up stack
    INX          ; now X = 0
    STX $2000    ; disable NMI
    STX $2001    ; disable rendering
    STX $4010    ; disable DMC IRQs
    JSR vblankwait

clrmem:
    LDA #$00
    STA $0000, x
    STA $0100, x
    STA $0300, x
    STA $0400, x
    STA $0500, x
    STA $0600, x
    STA $0700, x
    LDA #$FE
    STA $0200, x
    INX
    BNE clrmem
    JSR vblankwait

LoadPalettes:
    LDA $2002             ; read PPU status to reset the high/low latch
    LDA #$3F
    STA $2006             ; write the high byte of $3F00 address
    LDA #$00
    STA $2006             ; write the low byte of $3F00 address
    LDX #$00              ; start out at 0
LoadPalettesLoop:
    LDA palette, x        ; load data from address (palette + the value in x)
    STA $2007             ; write to PPU
    INX                   ; X = X + 1
    CPX #$20              ; Compare X to hex $10, decimal 16 - copying 16 bytes = 4 sprites
    BNE LoadPalettesLoop  ; Branch to LoadPalettesLoop if compare was Not Equal to zero

LoadMenuBackground:
    LDA $2002             ; read PPU status to reset the high/low latch
    LDA #$20
    STA $2006             ; write the high byte of $2000 address
    LDA #$00
    STA $2006             ; write the low byte of $2000 address
    LDA #$00
    JSR LoadNametable
    LDA #%10010000        ; enable NMI, sprites from Pattern Table 0, background from Pattern Table 1
    STA $2000
    LDA #%00011110        ; enable sprites, enable background, no clipping on left side
    STA $2001


soundConfig:
    LDA #%00000001
    STA $4015             ; enable square 1
    LDA #$C9              ; 0C9 is a C# in NTSC mode
    LDA #%10110000        ; Duty 10, Volume 0
    STA $4000
    STA $4002
    LDA #$00
    STA $4003

Forever:
    JMP Forever           ; jump back to Forever, infinite loop, waiting for NMI

;;; INIT TILES ;;;

initTiles:
    LDX #$0
loopInitTiles:
    TXA
    CMP #$10
    BEQ loopInitTilesDone
    LDA #$0
    STA tiles,x
    INX
    JMP loopInitTiles
loopInitTilesDone:
    LDA random
    LSR A
    LSR A
    LSR A
    LSR A
    AND #$0f ; mod 16
    TAX      ; transfer random value to X
    LDA #$01
    STA tiles,x

findEmptyInit:
    JSR updateRandom
    LDA random
    LSR A
    LSR A
    LSR A
    LSR A
    AND #$0f          ; mod 16
    TAX               ; transfer random value to X
    LDA tiles,x       ; load x-th tile
    CMP #$00
    BEQ initTwo       ; if tile is empty, fill
    JMP findEmptyInit ; else, try again

initTwo:
    LDA #$01
    STA tiles,x
    RTS

;;; END INIT TILES ;;;

NMI:
    LDA #$00
    STA $2003           ; set the low byte (00) of the RAM address
    LDA #$02
    STA $4014           ; set the high byte (02) of the RAM address, start the transfer


    ;; This is the PPU clean up section, so rendering the next frame starts properly.
    LDA #%10010000      ; enable NMI, sprites from Pattern Table 0, background from Pattern Table 1
    STA $2000
    LDA #%00011110      ; enable sprites, enable background, no clipping on left side
    STA $2001
    LDA #$00            ; tell the ppu there is no background scrolling
    STA $2005
    STA $2005

    ;; all graphics updates done by here, run game engine
    JSR ReadController1 ; get the current button data for player 1

GameEngine:
    JSR soundCheck
    JSR updateRandom

    LDA gamestate
    CMP #STATETITLE
    BEQ EngineTitle     ; game is displaying title screen

    LDA gamestate
    CMP #STATEGAMEOVER
    BEQ EngineGameOver  ; game is displaying ending screen


    LDA gamestate
    CMP #STATEPLAYING
    BEQ EnginePlaying   ; game is playing
GameEngineDone:
    RTI             ; return from interrupt




;;;;;;;;

EngineTitle:
    LDA buttons1
    AND #GAMEPAD_START
    CMP #GAMEPAD_START

    BNE GameEngineDone
    LDA #STATEPLAYING
    STA gamestate

    LDA #%00000000        ;Turn the screen off
    STA $2000
    STA $2001
    JSR LoadNametable
    LDA #%10001000        ;Turn the screen on
    STA $2000
    JSR initTiles
    JSR UpdateSprites

    JMP GameEngineDone

;;;;;;;;;

EngineGameOver:
    LDA buttons1
    AND #GAMEPAD_START
    CMP #GAMEPAD_START

    BNE GameEngineDone
    LDA #STATEPLAYING
    STA gamestate

    JSR initTiles

    LDA #%00000000        ;Turn the screen off
    STA $2000
    STA $2001
    JSR LoadNametable
    LDA #%10001000        ;Turn the screen on
    STA $2000
    JSR UpdateSprites

    JMP GameEngineDone

;;;;;;;;;;;

EnginePlaying:
    LDA buttons1
    AND #GAMEPAD_UP
    BEQ MPU1Done
    LDA lastPressed
    CMP #GAMEPAD_UP
    BEQ MPU1Done

doMvUp:
    JSR validUpMove
    CMP #$01
    BEQ MPU1Done
    JSR moveUp
    JSR moveUp
    JSR moveUp
    JSR mergeUp
    JSR moveUp
    JSR addTile
    JSR UpdateSprites

    LDA buttons1
    AND #GAMEPAD_UP
    STA lastPressed

    JSR playSound
MPU1Done:

    LDA buttons1
    AND #GAMEPAD_DOWN
    BEQ MPD1Done
    LDA lastPressed
    CMP #GAMEPAD_DOWN
    BEQ MPD1Done

doMvDown:
    JSR validDownMove
    CMP #$01
    BEQ MPD1Done
    JSR moveDown
    JSR moveDown
    JSR moveDown
    JSR mergeDown
    JSR moveDown
    JSR addTile
    JSR UpdateSprites
    LDA buttons1
    AND #GAMEPAD_DOWN
    STA lastPressed

    JSR playSound
MPD1Done:

    LDA buttons1
    AND #GAMEPAD_LEFT
    BEQ MPL1Done
    LDA lastPressed
    CMP #GAMEPAD_LEFT
    BEQ MPL1Done

doMvLeft:
    JSR validLeftMove
    CMP #$01
    BEQ MPL1Done
    JSR moveLeft
    JSR moveLeft
    JSR moveLeft
    JSR mergeLeft
    JSR moveLeft
    JSR addTile
    JSR UpdateSprites
    LDA buttons1
    AND #GAMEPAD_LEFT
    STA lastPressed

    JSR playSound
MPL1Done:

    LDA buttons1
    AND #GAMEPAD_RIGHT
    BEQ MPR1Done
    LDA lastPressed
    CMP #GAMEPAD_RIGHT
    BEQ MPR1Done

doMvRight:
    JSR validRightMove
    CMP #$01
    BEQ MPR1Done
    JSR moveRight
    JSR moveRight
    JSR moveRight
    JSR mergeRight
    JSR moveRight
    JSR addTile
    JSR UpdateSprites
    LDA buttons1
    AND #GAMEPAD_RIGHT
    STA lastPressed

    JSR playSound
    MPR1Done:
    JSR checkAnyMovesLeft

checkNonePressed:
    LDA buttons1
    AND #GAMEPAD_ANY_PRESSED
    BNE checkNonePressedDone
    LDA #%00000000
    STA lastPressed
checkNonePressedDone:
    JMP GameEngineDone

UpdateSprites:
    LDA #%00000000        ; Turn the screen off
    STA $2000
    STA $2001
    JSR LoadNametable
spriteLoop:
    LDA #$A5
    STA bgTileLo
    LDA #$20
    STA bgTileHi       ; draws the background from memory pos 2000

    ; calcula posicao de memoria do background da x-esima tile
    ; cada tile2048 -> 6x6 tiles do NES
    TXA
    AND #%00000011     ; A%4
    TAY
    horizLoop:
    CPY #$00
    BEQ horizLoopDone

    LDA bgTileLo     ; load low 8 bits of 16 bit value
    CLC              ; clear carry
    ADC #$06         ; add 6, as one tile2048 is 6 tiles wide
    STA bgTileLo     ; done with low bits, save back
    LDA bgTileHi     ; load upper 8 bits
    ADC #$00         ; add 0 and carry from previous add
    STA bgTileHi     ; save back

    DEY
    JMP horizLoop
horizLoopDone:

    TXA
    LSR A
    LSR A; A/4
    TAY
vertLoop:
    CPY #$00
    BEQ vertLoopDone

    LDA bgTileLo
    CLC
    ADC #$C0 ; add 6*32, as one tile2048 is 6 tiles tall, and one row has 32 tiles
    STA bgTileLo
    LDA bgTileHi
    ADC #$00
    STA bgTileHi

    DEY
    JMP vertLoop
vertLoopDone:

    LDA $2002             ; read PPU status to reset the high/low latch
    LDA bgTileHi
    STA $2006             ; write the high byte of $2000 address
    LDA bgTileLo
    STA $2006             ; write the low byte of $2000 address

    JSR DrawTile

    INX
    CPX #$10 ; 10 em hex eh 16 em dec
    BNE spriteLoop

    JSR calculateScore
    JSR drawScore

    LDA #%10001000        ; Turn the screen on
    STA $2000
    LDX #$00
    RTS

ReadController1:
    LDA #$01
    STA $4016
    LDA #$00
    STA $4016
    LDX #$08
ReadController1Loop:
    LDA $4016
    LSR A            ; bit0 -> Carry
    ROL buttons1     ; bit0 <- Carry
    DEX
    BNE ReadController1Loop
    RTS

DrawTile:
    LDA tiles,x ; load em A, o valor da x-esima tile
dt0:
    CMP #$00
    BNE dt2
    JSR tile0
dt2:
    CMP #$01
    BNE dt4
    JSR tile2
dt4:
    CMP #$02
    BNE dt8
    JSR tile4
dt8:
    CMP #$03
    BNE dt16
    JSR tile8
dt16:
    CMP #$04
    BNE dt32
    JSR tile16
dt32:
    CMP #$05
    BNE dt64
    JSR tile32
dt64:
    CMP #$06
    BNE dt128
    JSR tile64
dt128:
    CMP #$07
    BNE dt256
    JSR tile128
dt256:
    CMP #$08
    BNE dt512
    JSR tile256
dt512:
    CMP #$09
    BNE dt1024
    JSR tile512
dt1024:
    CMP #$0A
    BNE dt2048
    JSR tile1024
dt2048:
    CMP #$0B
    BNE tileDrawDone
    JSR tile2048
tileDrawDone:
    RTS

tile0:
    LDA #GBL
    STA $2007
    LDA #GBL
    STA $2007
    LDA #GBL
    STA $2007
    LDA #GBL
    STA $2007
    LDA #$FF ; um valor aleatorio pra n cair nas outras condicionais
    RTS
tile2:
    LDA #GBL
    STA $2007
    LDA #GBL
    STA $2007
    LDA #$02
    STA $2007
    LDA #GBL
    STA $2007
    LDA #$FF ; um valor aleatorio pra n cair nas outras condicionais
    RTS
tile4:
    LDA #GBL
    STA $2007
    LDA #GBL
    STA $2007
    LDA #$04
    STA $2007
    LDA #GBL
    STA $2007
    LDA #$FF ; um valor aleatorio pra n cair nas outras condicionais
    RTS
tile8:
    LDA #GBL
    STA $2007
    LDA #GBL
    STA $2007
    LDA #$08
    STA $2007
    LDA #GBL
    STA $2007
    LDA #$FF ; um valor aleatorio pra n cair nas outras condicionais
    RTS
tile16:
    LDA #GBL
    STA $2007
    LDA #$01
    STA $2007
    LDA #$06
    STA $2007
    LDA #GBL
    STA $2007
    LDA #$FF ; um valor aleatorio pra n cair nas outras condicionais
    RTS
tile32:
    LDA #GBL
    STA $2007
    LDA #$03
    STA $2007
    LDA #$02
    STA $2007
    LDA #GBL
    STA $2007
    LDA #$FF ; um valor aleatorio pra n cair nas outras condicionais
    RTS
tile64:
    LDA #GBL
    STA $2007
    LDA #$06
    STA $2007
    LDA #$04
    STA $2007
    LDA #GBL
    STA $2007
    LDA #$FF ; um valor aleatorio pra n cair nas outras condicionais
    RTS
tile128:
    LDA #GBL
    STA $2007
    LDA #$01
    STA $2007
    LDA #$02
    STA $2007
    LDA #$08
    STA $2007
    LDA #$FF ; um valor aleatorio pra n cair nas outras condicionais
    RTS
tile256:
    LDA #GBL
    STA $2007
    LDA #$02
    STA $2007
    LDA #$05
    STA $2007
    LDA #$06
    STA $2007
    LDA #$FF ; um valor aleatorio pra n cair nas outras condicionais
    RTS
tile512:
    LDA #GBL
    STA $2007
    LDA #$05
    STA $2007
    LDA #$01
    STA $2007
    LDA #$02
    STA $2007
    LDA #$FF ; um valor aleatorio pra n cair nas outras condicionais
    RTS
tile1024:
    LDA #$01
    STA $2007
    LDA #$00
    STA $2007
    LDA #$02
    STA $2007
    LDA #$04
    STA $2007
    LDA #$FF ; um valor aleatorio pra n cair nas outras condicionais
    RTS
tile2048:
    LDA #$02
    STA $2007
    LDA #$00
    STA $2007
    LDA #$04
    STA $2007
    LDA #$08
    STA $2007
    LDA #$FF ; um valor aleatorio pra n cair nas outras condicionais
    RTS


playSound:
    LDA #$00
    STA soundTimer
    LDA #%10111111 ;Duty 10, Volume F
    STA $4000
    RTS

soundCheck:
    LDA soundTimer
    TAX
    INX
    TXA
    STA soundTimer

    CMP #BEEP_DURATION
    BNE soundCheckDone
    LDA #%10110000 ;Duty 10, Volume 0
    STA $4000
soundCheckDone:
    RTS

;;;; MOVE RIGHT ;;;;
moveRight:
    LDX #$0      ; initialize indexes
    LDY #$0
loopMoveRight:
    CPX #$F      ; check if looked at all tiles, no need to check last one
    BEQ DONEmoveRight
    TXA          ; if not in the end
    AND #$03     ; mod 4 check if at the end of a tile row, then no need to check current tile
    CMP #$03
    BEQ SKIPMoveRight

    LDA tiles, x ; load the value of the current tile

    CMP #$00     ; if the tile is 0, no need to do anything
    BEQ SKIPMoveRight
    ;else
    INX          ; now we will check the next tile
    LDA tiles, x
    DEX
    CMP #$00     ; if the next tile is zero we can make the move, else there is nothing to be done
    BNE SKIPMoveRight
    ; else current not 0 and next 0 then swap
    LDA tiles, x ; load current tile again
    TAY          ; the current tile will be replaced with the value 0
    LDA #$0
    STA tiles, x ; save the value to zero
    TYA
    INX          ; now we'll make the swap
    STA tiles, x
    DEX
SKIPMoveRight:
    INX          ; if there's no swap to be done we just increment the pointer
    JMP loopMoveRight
DONEmoveRight:
    RTS
;;;; END MOVE RIGHT ;;;;

;;;; MOVE LEFT ;;;;
moveLeft:
    LDX #$f      ; initialize indexes (starting from left to right)
    LDY #$0
loopMoveLeft:
    CPX #$0      ; check if looked at all tiles, no need to check last one
    BEQ DONEmoveLeft
    TXA          ; if not in the end
    AND #$03     ; mod 4 check if at the end of a tile row, then no need to check current tile
    CMP #$00
    BEQ SKIPMoveLeft

    LDA tiles, x ; load the value of the current tile

    CMP #$00     ; if the tile is 0, no need to do anything
    BEQ SKIPMoveLeft
    ; else
    DEX          ; now we will check the next tile
    LDA tiles, x
    INX
    CMP #$00     ; if the next tile is zero we can make the move, else there is nothing to be done
    BNE SKIPMoveLeft
    ; else current not 0 and next 0 then swap
    LDA tiles, x ; load current tile again
    TAY          ; the current tile will be replaced with the value 0
    LDA #$0
    STA tiles, x ; save the value to zero
    TYA
    DEX          ; now we'll make the swap
    STA tiles, x
    INX
SKIPMoveLeft:
    DEX          ; if there's no swap to be done we just decrement the pointer
    JMP loopMoveLeft
DONEmoveLeft:
    RTS
;;;; END MOVE LEFT ;;;;


;;;; MOVE DOWN ;;;;
moveDown:
    LDY #$0      ; initialize indexes
    LDX #$0
loopDownOuter:
    TYA
    TAX          ; initialize X with Y
    CPY #04      ; done checking tiles
    BEQ DONEmoveDown
loopDownInner:
    TXA          ; if not in the end
    AND #%1100   ; check if its on the last line
    CMP #%1100
    BEQ doneLoopDownOuter

    LDA tiles, x ; load the value of the current tile

    CMP #$00     ; if the tile is 0, no need to do anything
    BEQ doneLoopDownInner
    ; else
    INX
    INX
    INX
    INX          ; now we will check the next tile

    LDA tiles, x
    DEX
    DEX
    DEX
    DEX
    CMP #$00     ; if the next tile is zero we can make the move, else there is nothing to be done
    BNE doneLoopDownInner
    ; else current not 0 and next 0 then swap
    LDA tiles, x ; load current tile again
    PHA          ; save current value to stack
    LDA #$0
    STA tiles, x ; save zero to current position
    PLA          ; retrieve previous value from stack
    INX			 ; now we'll make the swap
    INX
    INX
    INX
    STA tiles, x ; make the swap
    DEX
    DEX
    DEX
    DEX
    JMP doneLoopDownInner
doneLoopDownOuter:
    INY
    JMP loopDownOuter
doneLoopDownInner:
    INX          ; go to tile bellow
    INX
    INX
    INX
    JMP loopDownInner
DONEmoveDown:
    RTS
;;;; END MOVE DOWN ;;;;



;;;; MOVE UP ;;;;
moveUp:
    LDY #$3      ; initialize indexes
    TYA
    CLC
    ADC #$0C     ; add 12 (15)
    TAX
loopUpOuter:
    TYA
    CLC
    ADC #$0C
    TAX          ; initialize X with Y
    CPY #$FF     ; done checking tiles
    BEQ DONEmoveUp
loopUpInner:
    TXA          ; if not in the end
    AND #%1100   ; check if its on the last line
    CMP #$0
    BEQ doneLoopUpOuter

    LDA tiles, x ; load the value of the current tile

    CMP #$00     ; if the tile is 0, no need to do anything
    BEQ doneLoopUpInner
    ; else
    DEX
    DEX
    DEX
    DEX          ; now we will check the next tile

    LDA tiles, x
    INX
    INX
    INX
    INX
    CMP #$00     ; if the next tile is zero we can make the move, else there is nothing to be done
    BNE doneLoopUpInner
    ; else current not 0 and next 0 then swap
    LDA tiles, x ; load current tile again
    PHA          ; save current value to stack
    LDA #$0
    STA tiles, x ; save zero to current position
    PLA          ; retrieve previous value from stack
    DEX          ; now we'll make the swap
    DEX
    DEX
    DEX
    STA tiles, x ; make the swap
    INX
    INX
    INX
    INX
    JMP doneLoopUpInner
doneLoopUpOuter:
    DEY
    JMP loopUpOuter
doneLoopUpInner:
    DEX          ; go to tile bellow
    DEX
    DEX
    DEX
    JMP loopUpInner
DONEmoveUp:
    RTS
;;;; END MOVE UP ;;;;

;;; UPDATE RANDOM ;;;

; Linear congruential pseudo-random number generator
updateRandom:
    LDA random
    ASL A ; multiply by 5
    ASL A
    CLC
    ADC random
    CLC   ; add 17
    ADC #$17
    STA random
    RTS

;;; END UPDATE RANDOM ;;;

;;; ADD TILE ;;;

addTile:
    LDA random
    LSR A
    LSR A
    LSR A
    LSR A
    AND #$0f      ; mod 16
    STA beginTile ; store begin tile to check full cycle
    TAX           ; transfer random value to X

findEmpty:
    LDA tiles,x   ; load x-th tile
    CMP #$00
    BEQ twoORfour ; if tile is empty, fill
    INX           ; increment X
    TXA           ; transfer X to A
    AND #$0f      ; mod 16
    CMP beginTile ; if tried all tiles and none is empty, check if any moves left
    BEQ checkAnyMovesLeft
    TAX           ; transfer A to X
    JMP findEmpty

twoORfour:
    LDA random
    LSR A
    LSR A
    LSR A
    LSR A
    LSR A
    AND #$07 ; mod 8
    BEQ newFour ; if zero, draw four

newTwo:
    LDA #$01
    STA tiles,x
    RTS

newFour:
    LDA #$02
    STA tiles,x
    RTS

;;; END ADD NEW TILE ;;;

;;; CHECK ANY MOVES LEFT ;;;

checkAnyMovesLeft:
allFilled:
    LDX #$00
allFilledLoop:
    LDA tiles,x
    CMP #$00
    BEQ doneAllFilled ; at least one space is empty, therefore there are still moves left
    INX
    TXA
    CMP #$10
    BEQ compare01     ; if all filled, compare neighbors
    JMP allFilledLoop
doneAllFilled:
    RTS

compare01:
    LDX #$00
    LDA tiles,x
    INX
    CMP tiles,x ; compare tiles 0 and 1
    BNE compare12
    JMP doneCheckAnyMovesLeft
compare12:
    LDX #$01
    LDA tiles,x
    INX
    CMP tiles,x ; compare tiles 1 and 2
    BNE compare23
    JMP doneCheckAnyMovesLeft
compare23:
    LDX #$02
    LDA tiles,x
    INX
    CMP tiles,x ; compare tiles 2 and 3
    BNE compare45
    JMP doneCheckAnyMovesLeft
compare45:
    LDX #$04
    LDA tiles,x
    INX
    CMP tiles,x ; compare tiles 4 and 5
    BNE compare56
    JMP doneCheckAnyMovesLeft
compare56:
    LDX #$05
    LDA tiles,x
    INX
    CMP tiles,x ; compare tiles 5 and 6
    BNE compare67
    JMP doneCheckAnyMovesLeft
compare67:
    LDX #$06
    LDA tiles,x
    INX
    CMP tiles,x ; compare tiles 6 and 7
    BNE compare89
    JMP doneCheckAnyMovesLeft
compare89:
    LDX #$08
    LDA tiles,x
    INX
    CMP tiles,x ; compare tiles 8 and 9
    BNE compare9a
    JMP doneCheckAnyMovesLeft
compare9a:
    LDX #$09
    LDA tiles,x
    INX
    CMP tiles,x ; compare tiles 9 and 10
    BNE compareab
    JMP doneCheckAnyMovesLeft
compareab:
    LDX #$0a
    LDA tiles,x
    INX
    CMP tiles,x ; compare tiles 10 and 11
    BNE comparecd
    JMP doneCheckAnyMovesLeft
comparecd:
    LDX #$0c
    LDA tiles,x
    INX
    CMP tiles,x ; compare tiles 12 and 13
    BNE comparede
    JMP doneCheckAnyMovesLeft
comparede:
    LDX #$0d
    LDA tiles,x
    INX
    CMP tiles,x ; compare tiles 13 and 14
    BNE compareef
    JMP doneCheckAnyMovesLeft
compareef:
    LDX #$0e
    LDA tiles,x
    INX
    CMP tiles,x ; compare tiles 14 and 15
    BNE compare04
    JMP doneCheckAnyMovesLeft
compare04:
    LDX #$00
    LDA tiles,x
    INX
    INX
    INX
    INX
    CMP tiles,x ; compare tiles 0 and 4
    BNE compare48
    JMP doneCheckAnyMovesLeft
compare48:
    LDX #$04
    LDA tiles,x
    INX
    INX
    INX
    INX
    CMP tiles,x ; compare tiles 4 and 8
    BNE compare8c
    JMP doneCheckAnyMovesLeft
compare8c:
    LDX #$08
    LDA tiles,x
    INX
    INX
    INX
    INX
    CMP tiles,x ; compare tiles 8 and 12
    BNE compare15
    JMP doneCheckAnyMovesLeft
compare15:
    LDX #$01
    LDA tiles,x
    INX
    INX
    INX
    INX
    CMP tiles,x ; compare tiles 1 and 5
    BNE compare59
    JMP doneCheckAnyMovesLeft
compare59:
    LDX #$05
    LDA tiles,x
    INX
    INX
    INX
    INX
    CMP tiles,x ; compare tiles 5 and 9
    BNE compare9d
    JMP doneCheckAnyMovesLeft
compare9d:
    LDX #$09
    LDA tiles,x
    INX
    INX
    INX
    INX
    CMP tiles,x ; compare tiles 9 and 13
    BNE compare26
    JMP doneCheckAnyMovesLeft
compare26:
    LDX #$02
    LDA tiles,x
    INX
    INX
    INX
    INX
    CMP tiles,x ; compare tiles 2 and 6
    BNE compare6a
    JMP doneCheckAnyMovesLeft
compare6a:
    LDX #$06
    LDA tiles,x
    INX
    INX
    INX
    INX
    CMP tiles,x ; compare tiles 6 and 10
    BNE compareae
    JMP doneCheckAnyMovesLeft
compareae:
    LDX #$0a
    LDA tiles,x
    INX
    INX
    INX
    INX
    CMP tiles,x ; compare tiles 10 and 14
    BNE compare37
    JMP doneCheckAnyMovesLeft
compare37:
    LDX #$03
    LDA tiles,x
    INX
    INX
    INX
    INX
    CMP tiles,x ; compare tiles 3 and 7
    BNE compare7b
    JMP doneCheckAnyMovesLeft
compare7b:
    LDX #$07
    LDA tiles,x
    INX
    INX
    INX
    INX
    CMP tiles,x ; compare tiles 7 and 11
    BNE comparebf
    JMP doneCheckAnyMovesLeft
comparebf:
    LDX #$0b
    LDA tiles,x
    INX
    INX
    INX
    INX
    CMP tiles,x ; compare tiles 11 and 15
    BNE gameOver
    JMP doneCheckAnyMovesLeft

gameOver:
    JSR calculateScore
    LDA #STATEGAMEOVER
    STA gamestate
    LDA #%00000000        ;Turn the screen off
    STA $2000
    STA $2001
    JSR LoadNametable
	JSR drawScore
    LDA #%10001000        ;Turn the screen on
    STA $2000

GameOverForever:
    JMP GameOverForever

doneCheckAnyMovesLeft:
    RTS

;;; END CHECK ANY MOVES LEFT ;;;

;;; CALCULATE SCORE ;;;

calculateScore:
    LDA #$00
    STA scoreLo
    LDA #$00
    STA scoreHi

    LDX #$00
calculateScoreLoop:
    LDA tiles,x
    CMP #$02
    BEQ sum4
    CMP #$03
    BEQ sum16
    CMP #$04
    BEQ sum48
    CMP #$05
    BEQ sum128
    CMP #$06
    BEQ sum320
    JMP jumpLargeSum
returnCalculateScore:
    INX
    TXA
    CMP   #$10
    BEQ doneCalculateScoreLoop
    TAX
    JMP calculateScoreLoop
doneCalculateScoreLoop:
    RTS
sum4:
    CLC
    LDA scoreLo
    ADC #$04
    STA scoreLo
    LDA scoreHi
    ADC #$00
    STA scoreHi
    LDA #$FF
    JMP returnCalculateScore
sum16:
    CLC
    LDA scoreLo
    ADC #$10
    STA scoreLo
    LDA scoreHi
    ADC #$00
    STA scoreHi
    LDA #$FF
    JMP returnCalculateScore
sum48:
    CLC
    LDA scoreLo
    ADC #$30
    STA scoreLo
    LDA scoreHi
    ADC #$00
    STA scoreHi
    LDA #$FF
    JMP returnCalculateScore
sum128:
    CLC
    LDA scoreLo
    ADC #$80
    STA scoreLo
    LDA scoreHi
    ADC #$00
    STA scoreHi
    LDA #$FF
    JMP returnCalculateScore
sum320:
    CLC
    LDA scoreLo
    ADC #$40
    STA scoreLo
    LDA scoreHi
    ADC #$01
    STA scoreHi
    LDA #$FF
    JMP returnCalculateScore
jumpLargeSum:
    CMP #$07
    BEQ sum768
    CMP #$08
    BEQ sum1855
    CMP #$09
    BEQ sum4096
    CMP #$0a
    BEQ sum9216
    CMP #$0b
    BEQ sum20480
    JMP returnCalculateScore
sum768:
    CLC
    LDA scoreHi
    ADC #$03
    STA scoreHi
    LDA #$FF
    JMP returnCalculateScore
sum1855:
    CLC
    LDA scoreLo
    ADC #$3F
    STA scoreLo
    LDA scoreHi
    ADC #$07
    STA scoreHi
    LDA #$FF
    JMP returnCalculateScore
sum4096:
    CLC
    LDA scoreHi
    ADC #$10
    STA scoreHi
    LDA #$FF
    JMP returnCalculateScore
sum9216:
    CLC
    LDA scoreHi
    ADC #$24
    STA scoreHi
    LDA #$FF
    JMP returnCalculateScore
sum20480:
    CLC
    LDA scoreHi
    ADC #$50
    STA scoreHi
    LDA #$FF
    JMP returnCalculateScore

;;; END CALCULATE SCORE ;;;

;;; MERGE ;;;

;;  0,  1,  2,  3
;;  4,  5,  6,  7
;;  8,  9, 10, 11
;; 12, 13, 14, 15

mergeUp:
    LDX #$00        ; X <= the current tile idx
    LDY #$04        ; Y <= the idx of the tile to compare for the merge
mergeUpLoop:
    LDA tiles, X    ; get current tile value
    STA currTile    ; store it in a variable
    CMP #$00        ; if the current tile is 00 (empty) we don't even bother, just go to next
    BEQ mergeUpNext
    LDA tiles, Y    ; get the value of the tile to compare
    EOR currTile    ; check equals
    BNE mergeUpNext
    INC currTile
    LDA currTile
    STA tiles, X    ; replace current value in the tiles array
    LDA #$00        ; empty the tile which was used in merge
    STA tiles, Y
mergeUpNext:
    INX             ; go to next
    INY
    CPY #$10        ; if we've already checked every tile, stop
    BNE mergeUpLoop
mergeUpDone:
    RTS

mergeDown:
    LDX #$0B        ; X <= the idx of the tile to compare
    LDY #$0F        ; Y <= current tile idx
mergeDownLoop:
    LDA tiles, Y
    STA currTile
    CMP #$00
    BEQ mergeDownNext
    LDA tiles, X
    EOR currTile
    BNE mergeDownNext
    INC currTile
    LDA currTile
    STA tiles, Y
    LDA #$00
    STA tiles, X
mergeDownNext:
    DEX
    DEY
    CPY #$03
    BNE mergeDownLoop
mergeDownDone:
    RTS


mergeLeft:
    LDX #$00        ; X <= current tile idx
    LDY #$01        ; Y <= tile idx to compare
mergeLeftLoop:
    CPY #$04
    BEQ mergeLeftNext
    CPY #$08
    BEQ mergeLeftNext
    CPY #$0C
    BEQ mergeLeftNext
    LDA tiles, X
    STA currTile
    CMP #$00
    BEQ mergeLeftNext
    LDA tiles, Y
    EOR currTile
    BNE mergeLeftNext
    INC currTile
    LDA currTile
    STA tiles, X
    LDA #$00
    STA tiles, Y
mergeLeftNext:
    INX
    INY
    CPY #$10
    BNE mergeLeftLoop
mergeLeftDone:
    RTS

mergeRight:
    LDX #$0E        ; X <= tile idx to compare
    LDY #$0F        ; Y <= current tile
mergeRightLoop:
    CPY #$04
    BEQ mergeRightNext
    CPY #$08
    BEQ mergeRightNext
    CPY #$0C
    BEQ mergeRightNext
    LDA tiles, Y
    STA currTile
    CMP #$00
    BEQ mergeRightNext
    LDA tiles, X
    EOR currTile
    BNE mergeRightNext
    INC currTile
    LDA currTile
    STA tiles, Y
    LDA #$00
    STA tiles, X
mergeRightNext:
    DEX
    DEY
    CPY #$00
    BNE mergeRightLoop
mergeRightDone:
    RTS

;;; END MERGE ;;;

drawScore:
    LDA scoreLo
    STA scoreTempLo
    LDA scoreHi
    STA scoreTempHi

    JSR scoreDig5
    JSR scoreDig4
    JSR scoreDig3
    JSR scoreDig2
    JSR scoreDig1
    ; draw tiles

    LDA $2002             ; read PPU status to reset the high/low latch
    LDA #$23
    STA $2006             ; write the high byte of $3F00 address
    LDA #$6D
    STA $2006             ; write the low byte of $3F00 address

    LDA #$00
    TAX

    LDA scoreDig,X
    CLC
    ADC #$D0
    STA $2007
    INX

    LDA scoreDig,X
    CLC
    ADC #$D0	
    STA $2007
    INX

    LDA scoreDig,X
    CLC
    ADC #$D0
    STA $2007
    INX

    LDA scoreDig,X
    CLC
    ADC #$D0
    STA $2007
    INX

    LDA scoreDig,X
    CLC
    ADC #$D0
    STA $2007

    RTS

scoreDig5:
    LDX #$00
    LDA scoreTempLo  ; backup when it goes negative
    STA score2Lo
    LDA scoreTempHi
    STA score2Hi

    scoreDig5Loop:   ; 10000 is $2710
    LDA scoreTempLo  ; load low 8 bits of 16 bit value
    SEC              ; set carry
    SBC #$10         ; subtract
    STA scoreTempLo  ; done with low bits, save back
    LDA scoreTempHi  ; load upper 8 bits
    SBC #$27         ; subtract
    STA scoreTempHi  ; save back

    ; if not negative
    BMI scoreDig5Done

    INX
    LDA scoreTempLo  ; backup when it goes negative
    STA score2Lo
    LDA scoreTempHi
    STA score2Hi
    JMP scoreDig5Loop

scoreDig5Done:
    LDA score2Lo     ; restores, because it went negative
    STA scoreTempLo
    LDA score2Hi
    STA scoreTempHi

    LDY #$00
    TXA
    STA scoreDig,Y   ; save the digit

    RTS

scoreDig4:
    LDX #$00
    LDA scoreTempLo  ; backup when it goes negative
    STA score2Lo
    LDA scoreTempHi
    STA score2Hi

scoreDig4Loop:       ; 1000 is $03E8
    LDA scoreTempLo  ; load low 8 bits of 16 bit value
    SEC              ; set carry
    SBC #$E8         ; subtract
    STA scoreTempLo  ; done with low bits, save back
    LDA scoreTempHi  ; load upper 8 bits
    SBC #$03         ; subtract
    STA scoreTempHi  ; save back

    ; if not negative
    BMI scoreDig4Done

    INX
    LDA scoreTempLo  ; backup when it goes negative
    STA score2Lo
    LDA scoreTempHi
    STA score2Hi
    JMP scoreDig4Loop

scoreDig4Done:
    LDA score2Lo     ; restores, because it went negative
    STA scoreTempLo
    LDA score2Hi
    STA scoreTempHi

    LDY #$01
    TXA
    STA scoreDig,Y   ; save the digit

    RTS

scoreDig3:
    LDX #$00
    LDA scoreTempLo  ; backup when it goes negative
    STA score2Lo
    LDA scoreTempHi
    STA score2Hi

scoreDig3Loop:       ; 100 is $0064
    LDA scoreTempLo  ; load low 8 bits of 16 bit value
    SEC              ; set carry
    SBC #$64         ; subtract
    STA scoreTempLo  ; done with low bits, save back
    LDA scoreTempHi  ; load upper 8 bits
    SBC #$00         ; subtract
    STA scoreTempHi  ; save back

    ; if not negative
    BMI scoreDig3Done

    INX
    LDA scoreTempLo  ; backup when it goes negative
    STA score2Lo
    LDA scoreTempHi
    STA score2Hi
    JMP scoreDig3Loop

scoreDig3Done:
    LDA score2Lo     ; restores, because it went negative
    STA scoreTempLo
    LDA score2Hi
    STA scoreTempHi

    LDY #$02
    TXA
    STA scoreDig,Y   ; save the digit

    RTS

scoreDig2:
    LDX #$00
    LDA scoreTempLo  ; backup when it goes negative
    STA score2Lo
    LDA scoreTempHi
    STA score2Hi

scoreDig2Loop:       ; 10 is $000A
    LDA scoreTempLo  ; load low 8 bits of 16 bit value
    SEC              ; set carry
    SBC #$0A         ; subtract
    STA scoreTempLo  ; done with low bits, save back
    LDA scoreTempHi  ; load upper 8 bits
    SBC #$00         ; subtract
    STA scoreTempHi  ; save back

    ; if not negative
    BMI scoreDig2Done

    INX
    LDA scoreTempLo  ; backup when it goes negative
    STA score2Lo
    LDA scoreTempHi
    STA score2Hi
    JMP scoreDig2Loop

scoreDig2Done:
    LDA score2Lo     ; restores, because it went negative
    STA scoreTempLo
    LDA score2Hi
    STA scoreTempHi

    LDY #$03
    TXA
    STA scoreDig,Y   ; save the digit

    RTS

scoreDig1:
    LDX #$00
    LDA scoreTempLo  ; backup when it goes negative
    STA score2Lo
    LDA scoreTempHi
    STA score2Hi

scoreDig1Loop:       ; 1 is $0001
    LDA scoreTempLo  ; load low 8 bits of 16 bit value
    SEC              ; set carry
    SBC #$01         ; subtract
    STA scoreTempLo  ; done with low bits, save back
    LDA scoreTempHi  ; load upper 8 bits
    SBC #$00         ; subtract
    STA scoreTempHi  ; save back

    ; if not negative
    BMI scoreDig1Done

    INX
    LDA scoreTempLo ; backup when it goes negative
    STA score2Lo
    LDA scoreTempHi
    STA score2Hi
    JMP scoreDig1Loop

scoreDig1Done:
    LDA score2Lo ; restores, because it went negative
    STA scoreTempLo
    LDA score2Hi
    STA scoreTempHi

    LDY #$04
    TXA
    STA scoreDig,Y ; save the digit

    RTS

;;; VALID UP MOVE ;;;

; Returns 0 in A if valid, 1 otherwise
validUpMove:
    LDX #$00
    LDY #$04
    validUpMoveLoop:
    LDA tiles,x
    CMP #$00
    BEQ checkBottomTile ; if top tile is zero, check if bottom isn't
    ; else, check if both are equal
    CMP tiles,y
    BEQ upMoveOK        ; if tiles,x == tiles,y != 0, OK
    ; else, check next
    INX
    INY
    TYA
    CMP #$10
    BEQ notValidUpMove  ; no more pair of tiles to check
    JMP validUpMoveLoop
checkBottomTile:
    LDA tiles,y
    CMP #$00 ;
    BNE upMoveOK        ; if tiles,x == 0 and tiles,y != 0, OK
    ; else, check next
    INX
    INY
    TYA
    CMP #$10
    BEQ notValidUpMove  ; no more pair of tiles to check
    JMP validUpMoveLoop
upMoveOK:
    LDA #$00
    RTS
notValidUpMove:
    LDA #$01
    RTS

;;; END VALID UP MOVE ;;;

;;; VALID DOWN MOVE ;;;

; Returns 0 in A if valid, 1 otherwise
validDownMove:
    LDX #$04
    LDY #$00
validDownMoveLoop:
    LDA tiles,x
    CMP #$00
    BEQ checkTopTile     ; if bottom tile is zero, check if top isn't
    ; else, check if both are equal
    CMP tiles,y
    BEQ downMoveOK       ; if tiles,x == tiles,y != 0, OK
    ; else, check next
    INX
    INY
    TXA
    CMP #$10
    BEQ notValidDownMove ; no more pair of tiles to check
    JMP validDownMoveLoop
checkTopTile:
    LDA tiles,y
    CMP #$00 ;
    BNE downMoveOK       ; if tiles,x == 0 and tiles,y != 0, OK
    ; else, check next
    INX
    INY
    TXA
    CMP #$10
    BEQ notValidDownMove ; no more pair of tiles to check
    JMP validDownMoveLoop
downMoveOK:
    LDA #$00
    RTS
notValidDownMove:
    LDA #$01
    RTS

;;; END VALID DOWN MOVE ;;;

;;; VALID LEFT MOVE ;;;

; Returns 0 in A if valid, 1 otherwise
validLeftMove:
    LDX #$00
    LDY #$01
validLeftMoveLoop:
    LDA tiles,x
    CMP #$00
    BEQ checkRightTile   ; if left tile is zero, check if right isn't
    ; else, check if both are equal
    CMP tiles,y
    BEQ leftMoveOK       ; if tiles,x == tiles,y != 0, OK
    ; else, check next
    INX
    INY
    TXA
    AND #$03
    CMP #$03
    BNE noLineChange1    ; increases index twice when line ends so as not to compare the last tile of the row with the first tile of the next row
    INX
    INY
    TXA
noLineChange1:
    CMP #$10
    BEQ notValidLeftMove ; no more pair of tiles to check
    JMP validLeftMoveLoop
    checkRightTile:
    LDA tiles,y
    CMP #$00 ;
    BNE leftMoveOK       ; if tiles,x == 0 and tiles,y != 0, OK
    ; else, check next
    INX
    INY
    TXA
    AND #$03
    CMP #$03
    BNE noLineChange2    ; increases index twice when line ends so as not to compare the last tile of the row with the first tile of the next row
    INX
    INY
    TXA
noLineChange2:
    CMP #$10
    BEQ notValidLeftMove ; no more pair of tiles to check
    JMP validLeftMoveLoop
leftMoveOK:
    LDA #$00
    RTS
notValidLeftMove:
    LDA #$01
    RTS

;;; END VALID LEFT MOVE ;;;

;;; VALID RIGHT MOVE ;;;

; Returns 0 in A if valid, 1 otherwise
validRightMove:
    LDX #$01
    LDY #$00
validRightMoveLoop:
    LDA tiles,x
    CMP #$00
    BEQ checkLeftTile     ; if right tile is zero, check if left isn't
    ; else, check if both are equal
    CMP tiles,y
    BEQ rightMoveOK       ; if tiles,x == tiles,y != 0, OK
    ; else, check next
    INX
    INY
    TYA
    AND #$03
    CMP #$03
    BNE noLineChange3     ; increases index twice when line ends so as not to compare the last tile of the row with the first tile of the next row
    INX
    INY
    TYA
noLineChange3:
    CMP #$10
    BEQ notValidRightMove ; no more pair of tiles to check
    JMP validRightMoveLoop
checkLeftTile:
    LDA tiles,y
    CMP #$00 ;
    BNE rightMoveOK       ; if tiles,x == 0 and tiles,y != 0, OK
    ; else, check next
    INX
    INY
    TYA
    AND #$03
    CMP #$03
    BNE noLineChange4     ; increases index twice when line ends so as not to compare the last tile of the row with the first tile of the next row
    INX
    INY
    TYA
noLineChange4:
    CMP #$10
    BEQ notValidRightMove ; no more pair of tiles to check
    JMP validRightMoveLoop
rightMoveOK:
    LDA #$00
    RTS
notValidRightMove:
    LDA #$01
    RTS

;;; END VALID RIGHT MOVE ;;;

LoadNametable:
    LDA $2002     ; read PPU status to reset the high/low latch
    LDA #$20
    STA $2006
    LDA #$00
    STA $2006

    ; Set the nametable
    LDA gamestate
    ASL A
    TAY
    LDA NametablePointerTable, Y
    STA tempWord
    LDA NametablePointerTable+1, Y
    STA tempWord+1

    LDX #$04
    LDY #$00
    ; Load the nametable (change this, attribute is being set unnecessarily)
.loadNametableLoop:
    LDA [tempWord], Y              ; load nametable
    STA $2007                      ; draw tile
    INY
    BNE .loadNametableLoop
    INC tempWord+1
    DEX
    BNE .loadNametableLoop
    RTS

;;;;;;;;;;;;;;

    .bank 1
    .org $E000    ;;align the background data so the lower address is $00

NametablePointerTable:
    .dw menuBackground     ;STATE_TITLE
    .dw background   ;STATE_PLAYING
    .dw gameoverBackground ;STATE_GAMEOVER

menuBackground:
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, DC2, DC0, DC4, DC8, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, DCP, DCR, DCE, DCS, DCS, GBG, DCS, DCT, DCA, DCR, DCT,  GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG

background:
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GDG, GDG, GDG, GDG, GDG, GDG, GDG, GDG, GDG, GDG, GDG, GDG, GDG, GDG, GDG, GDG, GDG, GDG, GDG, GDG, GDG, GDG, GDG, GDG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GRG, GUL, GUB, GUB, GUB, GUB, GUR, GUL, GUB, GUB, GUB, GUB, GUR, GUL, GUB, GUB, GUB, GUB, GUR, GUL, GUB, GUB, GUB, GUB, GUR, GLG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GRG, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GRG, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GRG, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GRG, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GRG, GDL, GDB, GDB, GDB, GDB, GDR, GDL, GDB, GDB, GDB, GDB, GDR, GDL, GDB, GDB, GDB, GDB, GDR, GDL, GDB, GDB, GDB, GDB, GDR, GLG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GRG, GUL, GUB, GUB, GUB, GUB, GUR, GUL, GUB, GUB, GUB, GUB, GUR, GUL, GUB, GUB, GUB, GUB, GUR, GUL, GUB, GUB, GUB, GUB, GUR, GLG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GRG, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GRG, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GRG, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GRG, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GRG, GDL, GDB, GDB, GDB, GDB, GDR, GDL, GDB, GDB, GDB, GDB, GDR, GDL, GDB, GDB, GDB, GDB, GDR, GDL, GDB, GDB, GDB, GDB, GDR, GLG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GRG, GUL, GUB, GUB, GUB, GUB, GUR, GUL, GUB, GUB, GUB, GUB, GUR, GUL, GUB, GUB, GUB, GUB, GUR, GUL, GUB, GUB, GUB, GUB, GUR, GLG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GRG, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GRG, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GRG, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GRG, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GRG, GDL, GDB, GDB, GDB, GDB, GDR, GDL, GDB, GDB, GDB, GDB, GDR, GDL, GDB, GDB, GDB, GDB, GDR, GDL, GDB, GDB, GDB, GDB, GDR, GLG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GRG, GUL, GUB, GUB, GUB, GUB, GUR, GUL, GUB, GUB, GUB, GUB, GUR, GUL, GUB, GUB, GUB, GUB, GUR, GUL, GUB, GUB, GUB, GUB, GUR, GLG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GRG, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GRG, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GRG, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GRG, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLB, GBL, GBL, GBL, GBL, GRB, GLG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GRG, GDL, GDB, GDB, GDB, GDB, GDR, GDL, GDB, GDB, GDB, GDB, GDR, GDL, GDB, GDB, GDB, GDB, GDR, GDL, GDB, GDB, GDB, GDB, GDR, GLG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GUG, GUG, GUG, GUG, GUG, GUG, GUG, GUG, GUG, GUG, GUG, GUG, GUG, GUG, GUG, GUG, GUG, GUG, GUG, GUG, GUG, GUG, GUG, GUG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG

gameoverBackground:
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, DC2, DC0, DC4, DC8, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, DCG, DCA, DCM, DCE, GBG, GBG, DCO, DCV, DCE, DCR, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, DCP, DCR, DCE, DCS, DCS, GBG, DCS, DCT, DCA, DCR, DCT, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, DCF, DCI, DCN, DCA, DCL, GBG, DCS, DCC, DCO, DCR, DCE,  GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG
    .db GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG, GBG

attributes:  ;8 x 8 = 64 bytes
    .db %00000000, %00000000, %00000000, %00000000, %00000000, %00000000, %00000000, %00000000
    .db %00000000, %00000000, %00000000, %00000000, %00000000, %00000000, %00000000, %00000000
    .db %00000000, %00000000, %00000000, %00000000, %00000000, %00000000, %00000000, %00000000
    .db %00000000, %00000000, %00000000, %00000000, %00000000, %00000000, %00000000, %00000000
    .db %00000000, %00000000, %00000000, %00000000, %00000000, %00000000, %00000000, %00000000
    .db %00000000, %00000000, %00000000, %00000000, %00000000, %00000000, %00000000, %00000000
    .db %00000000, %00000000, %00000000, %00000000, %00000000, %00000000, %00000000, %00000000
    .db %00000000, %00000000, %00000000, %00000000, %00000000, %00000000, %00000000, %00000000


palette:
    .db $22,$29,$1A,$0F,  $22,$29,$1A,$0F,  $22,$29,$1A,$0F,  $22,$29,$1A,$0F   ;;background palette
    .db $22,$1C,$15,$14,  $22,$02,$38,$3C,  $22,$1C,$15,$14,  $22,$02,$38,$3C   ;;sprite palette

sprites:
    ;vert tile attr horiz
    .db $80, $32, $00, $80   ;sprite 0
    .db $80, $24, $00, $88   ;sprite 1
    .db $88, $24, $00, $80   ;sprite 2
    .db $88, $24, $00, $88   ;sprite 3

    .org $FFFA     ;first of the three vectors starts here
    .dw NMI        ;when an NMI happens (once per frame if enabled) the
    ;processor will jump to the label NMI:
    .dw RESET      ;when the processor first turns on or is reset, it will jump
    ;to the label RESET:
    .dw 0          ;external interrupt IRQ is not used in this tutorial
;;;;;;;;;;;;;;
    .bank 2
    .org $0000
    .incbin "sprite.chr"   ;includes 8KB graphics file from SMB1

