	.inesprg 1   ; 1x 16KB PRG code
	.ineschr 1   ; 1x  8KB CHR data
	.inesmap 0   ; mapper 0 = NROM, no bank swapping
	.inesmir 1   ; background mirroring

	;; DECLARE SOME VARIABLES HERE
	.rsset $0000  ;;start variables at ram location 0
gamestate	.rs 1  ; .rs 1 means reserve one byte of space
pointerLo	.rs 1   ; pointer variables are declared in RAM
pointerHi	.rs 1   ; low byte first, high byte immediately after
buttons1	.rs 1  ; player 1 gamepad buttons, one bit per button
buttons2	.rs 1  ; player 2 gamepad buttons, one bit per button
bgTileLo	.rs 1
bgTileHi	.rs 1
lastPressed .rs 1
tiles	.rs 16

	;; DECLARE SOME CONSTANTS HERE
STATETITLE     = $00  ; displaying title screen
STATEPLAYING   = $01  ; move paddles/ball, check for collisions
STATEGAMEOVER  = $02  ; displaying game over screen


GUL = $0e
GUR = $0f
GDL = $10
GDR = $11
GLB = $0a
GRB = $0b
GUB = $0c
GDB = $0d
GBL = $26
GBG = $24
GLG = $12
GRG = $13
GUG = $14
GDG = $15
    ;STATETITLE     = $00  ; displaying title screen

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
	STA pointerLo       ; put the low byte of the address of background into pointer
	LDA #HIGH(background)
	STA pointerHi       ; put the high byte of the address into pointer
	LDX #$00            ; start at pointer + 0
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
	LDA #%10010000   ; enable NMI, sprites from Pattern Table 0, background from Pattern Table 1
	STA $2000
	LDA #%00011110   ; enable sprites, enable background, no clipping on left side
	STA $2001


	LDA #STATEPLAYING
	STA gamestate
Forever:
	JMP Forever     ;jump back to Forever, infinite loop, waiting for NMI

NMI:
	LDA #$00
	STA $2003       ; set the low byte (00) of the RAM address
	LDA #$02
	STA $4014       ; set the high byte (02) of the RAM address, start the transfer

	;;This is the PPU clean up section, so rendering the next frame starts properly.
	LDA #%10010000   ; enable NMI, sprites from Pattern Table 0, background from Pattern Table 1
	STA $2000
	LDA #%00011110   ; enable sprites, enable background, no clipping on left side
	STA $2001
	LDA #$00        ;;tell the ppu there is no background scrolling
	STA $2005
	STA $2005

	;;;all graphics updates done by here, run game engine
	JSR ReadController1  ;;get the current button data for player 1
	JSR ReadController2  ;;get the current button data for player 2

GameEngine:  
	LDA gamestate
	CMP #STATETITLE
	BEQ EngineTitle    ;;game is displaying title screen

	LDA gamestate
	CMP #STATEGAMEOVER
	BEQ EngineGameOver  ;;game is displaying ending screen

	LDA gamestate
	CMP #STATEPLAYING
	BEQ EnginePlaying   ;;game is playing
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
	AND #%00001000
	BEQ MPU1Done
	LDA lastPressed
	CMP #%00001000
	BEQ MPU1Done


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
MPU1Done:


	LDA buttons1
	AND #%00001000
	STA lastPressed

	JMP GameEngineDone


UpdateSprites:
	LDX #$00
spriteLoop:
	LDA #$00
	STA bgTileLo     
	LDA #$20
	STA bgTileHi       ; draws the background from memory pos 2000

    ; calcula posicao de memoria do background da x-esima tile
    ; cada tile2048 -> 7x7 tiles do NES
    TXA 
    AND #%00000011 ; A%4
    TAY
horizLoop:
	CPY #$00
	BEQ horizLoopDone

	LDA bgTileLo      ; load low 8 bits of 16 bit value
	CLC              ; clear carry
	ADC #$07         ; add 7, as one tile2048 is 7 tiles wide
	STA bgTileLo      ; done with low bits, save back
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
	ADC #$E0 ; add 7*32, as one tile2048 is 7 tiles tall, and one row has 32 tiles
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

ReadController2:
	LDA #$01
	STA $4016
	LDA #$00
	STA $4016
	LDX #$08
ReadController2Loop:
	LDA $4017
	LSR A            ; bit0 -> Carry
	ROL buttons2     ; bit0 <- Carry
	DEX
	BNE ReadController2Loop
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
	LDA #$00
	STA $2007
	LDA #$FF ; um valor aleatorio pra n cair nas outras condicionais
	RTS
tile2:
	LDA #$02
	STA $2007
	LDA #$FF ; um valor aleatorio pra n cair nas outras condicionais
	RTS
tile4:
	LDA #$04
	STA $2007
	LDA #$FF ; um valor aleatorio pra n cair nas outras condicionais
	RTS
tile8:
	LDA #$08
	STA $2007
	LDA #$FF ; um valor aleatorio pra n cair nas outras condicionais
	RTS
tile16:
	LDA #$01
	STA $2007
	LDA #$06
	STA $2007
	LDA #$FF ; um valor aleatorio pra n cair nas outras condicionais
	RTS
tile32:
	LDA #$03
	STA $2007
	LDA #$02
	STA $2007
	LDA #$FF ; um valor aleatorio pra n cair nas outras condicionais
	RTS
tile64:
	LDA #$06
	STA $2007
	LDA #$04
	STA $2007
	LDA #$FF ; um valor aleatorio pra n cair nas outras condicionais
	RTS
tile128:
	LDA #$01
	STA $2007
	LDA #$02
	STA $2007
	LDA #$08
	STA $2007
	LDA #$FF ; um valor aleatorio pra n cair nas outras condicionais
	RTS
tile256:
	LDA #$02
	STA $2007
	LDA #$05
	STA $2007
	LDA #$06
	STA $2007
	LDA #$FF ; um valor aleatorio pra n cair nas outras condicionais
	RTS
tile512:
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
