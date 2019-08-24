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
	.rsset $0000  ;;start variables at ram location 0

gamestate	.rs 1  ; .rs 1 means reserve one byte of space
pointerLo	.rs 1  ; pointer variables are declared in RAM
pointerHi	.rs 1  ; low byte first, high byte immediately after
buttons1	.rs 1  ; player 1 gamepad buttons, one bit per button
bgTileLo	.rs 1
bgTileHi	.rs 1
lastPressed .rs 1
tiles       .rs 16
ramdomSeed .rs 1
soundTimer .rs 1

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
;STATETITLE     = $00  ; displaying title screen
;SOUND
BEEP_DURATION = $03
;;;;;;;;;;;;;;;;;;;

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

LoadBackground:
	LDA $2002             ; read PPU status to reset the high/low latch
	LDA #$20
	STA $2006             ; write the high byte of $2000 address
	LDA #$00
	STA $2006             ; write the low byte of $2000 address
	LDA #$00
	STA pointerLo         ; put the low byte of the address of background into pointer
	LDA #HIGH(background)
	STA pointerHi         ; put the high byte of the address into pointer
	LDX #$00              ; start at pointer + 0
	LDY #$00
OutsideLoop:

InsideLoop:
	LDA [pointerLo], y  ; copy one background byte from address in pointer plus Y
	STA $2007           ; this runs 256 * 4 times
	INY                 ; inside loop counter
	CPY #$00
	BNE InsideLoop      ; run the inside loop 256 times before continuing down
	INC pointerHi       ; low byte went 0 to 256, so high byte needs to be changed now
	INX
	CPX #$04
	BNE OutsideLoop     ; run the outside loop 256 times before continuing down

	LDA #%10010000      ; enable NMI, sprites from Pattern Table 0, background from Pattern Table 1
	STA $2000
	LDA #%00011110      ; enable sprites, enable background, no clipping on left side
	STA $2001


	LDA #STATEPLAYING
	STA gamestate

soundConfig:
	LDA #%00000001
	STA $4015 ;enable square 1
	LDA #$C9    ;0C9 is a C# in NTSC mode
	LDA #%10110000 ;Duty 10, Volume 0
	STA $4000
	STA $4002
	LDA #$00
	STA $4003

Forever:
	JMP Forever         ; jump back to Forever, infinite loop, waiting for NMI

NMI:
	LDA #$00
	STA $2003           ; set the low byte (00) of the RAM address
	LDA #$02
	STA $4014           ; set the high byte (02) of the RAM address, start the transfer


	;;This is the PPU clean up section, so rendering the next frame starts properly.
	LDA #%10010000      ; enable NMI, sprites from Pattern Table 0, background from Pattern Table 1
	STA $2000
	LDA #%00011110      ; enable sprites, enable background, no clipping on left side
	STA $2001
	LDA #$00            ; tell the ppu there is no background scrolling
	STA $2005
	STA $2005

	;;; all graphics updates done by here, run game engine
	JSR ReadController1 ; get the current button data for player 1

RandomSeed:
	LDA #$03
	STA ramdomSeed

GameEngine:
	JSR soundCheck

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

	;JSR UpdateSprites
	RTI             ; return from interrupt



;;;;;;;;

EngineTitle:
	;;if start button pressed
	;;  turn screen off
	;;  load game screen
	;;  set starting paddle/ball position
	;;  go to Playing State
	;;  turn screen on
	JMP GameEngineDone

;;;;;;;;;

EngineGameOver:
	;;if start button pressed
	;;  turn screen off
	;;  load title screen
	;;  go to Title State
	;;  turn screen on
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
	LDA #$01
	STA tiles
	LDA #$02
	LDX #$05
	STA tiles, x
	LDA #$08
	LDX #$0D
	STA tiles, x
	LDA #$0A
	LDX #$0E
	STA tiles, x
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
    LDA #$02
    STA tiles
    LDA #$0A
    LDX #$02
    STA tiles, x
    LDA #$08
    LDX #$0D
    STA tiles, x
    LDA #$0A
    LDX #$0E
    STA tiles, x
    JSR UpdateSprites
    LDA buttons1
    AND #GAMEPAD_DOWN
    STA lastPressed

    JSR playSound
MPD1Done:
	JMP GameEngineDone

UpdateSprites:
	LDX #$00
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
	CPX #$10 ;10 em hex eh 16 em dec
    BNE spriteLoop

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


random:
	LDA ramdomSeed
	ASL A
	ASL A
	CLC
	ADC ramdomSeed
	CLC
	ADC #$17
	STA ramdomSeed
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


moveRight:
	LDX #$0	;initialize indexes
	LDY #$0
loopMoveRight:
	CPX #$F	; check if looked at all tiles, no need to check last one
	BEQ DONEmoveRight
	TXA ; if not in the end
	AND #$03 ; mod 4 check if at the end of a tile row, then no need to check current tile
	CMP #$03
	BEQ SKIPMoveRight

	LDA tiles, x ; load the value of the current tile

	CMP #$00 ; if the tile is 0, no need to do anything
	BEQ SKIPMoveRight
	;else
	INX ; now we will check the next tile 
	LDA tiles, x
	DEX
	CMP #$00 ; if the next tile is zero we can make the move, else there is nothing to be done
	BNE SKIPMoveRight
	;else current not 0 and next 0 then swap
	LDA tiles, x ; load current tile again
	TAY  ; the current tile will be replaced with the value 0
	LDA #$0
	STA tiles, x ; save the value to zero
	TYA 
	INX			 ; now we'll make the swap
	STA tiles, x
	DEX
SKIPMoveRight:
	INX			; if there's no swap to be done we just increment the pointer
	JMP loopMoveRight
DONEmoveRight:
	RTS

;;;;;;;;;;;;;;
  .bank 1
  .org $E000    ;;align the background data so the lower address is $00

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

