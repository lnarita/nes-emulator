; Author: tokumaru
; http://forums.nesdev.com/viewtopic.php?%20p=58138#58138
;----------------------------------------------------------------
; constants
;----------------------------------------------------------------
PRG_COUNT = 1 ;1 = 16KB, 2 = 32KB
MIRRORING = %0001 ;%0000 = horizontal, %0001 = vertical, %1000 = four-screen

;----------------------------------------------------------------
; variables
;----------------------------------------------------------------

   .enum $0000

   ;NOTE: declare variables using the DSB and DSW directives, like this:

   ;MyVariable0 .dsb 1
   ;MyVariable1 .dsb 3

   .ende

   ;NOTE: you can also split the variable declarations into individual pages, like this:

   ;.enum $0100
   ;.ende

   ;.enum $0200
   ;.ende

;----------------------------------------------------------------
; iNES header
;----------------------------------------------------------------

   .db "NES", $1a ;identification of the iNES header
   .db PRG_COUNT ;number of 16KB PRG-ROM pages
   .db $01 ;number of 8KB CHR-ROM pages
   .db $00|MIRRORING ;mapper 0 and mirroring
   .dsb 9, $00 ;clear the remaining bytes

;----------------------------------------------------------------
; program bank(s)
;----------------------------------------------------------------

   .base $10000-(PRG_COUNT*$4000)


NMI:

   ;NOTE: NMI code goes here

IRQ:

   ;NOTE: IRQ code goes here

Reset:

LoadPalettes:
   LDA $2002             ; read PPU status to reset the high/low latch
   LDA #$3F
   STA $2006             ; write the high byte of $3F00 address
   LDA #$00
   STA $2006             ; write the low byte of $3F00 address
   LDX #$00              ; start out at 0
LoadPalettesLoop:
   LDA palette, x        ; load data from address (palette + the value in x)
                         ; 1st time through loop it will load palette+0
                         ; 2nd time through loop it will load palette+1
                         ; 3rd time through loop it will load palette+2
                         ; etc
   STA $2007             ; write to PPU
   INX                   ; X = X + 1
   CPX #$20              ; Compare X to hex $10, decimal 16 - copying 16 bytes = 4 sprites
   BNE LoadPalettesLoop  ; Branch to LoadPalettesLoop if compare was Not Equal to zero
                         ; if compare was equal to 32, keep going down

LoadSprites:
   LDX #$00              ; start at 0
LoadSpritesLoop:
   LDA sprites, x        ; load data from address (sprites +  x)
   STA $0200, x          ; store into RAM address ($0200 + x)
   INX                   ; X = X + 1
   CPX #$10              ; Compare X to hex $10, decimal 16
   BNE LoadSpritesLoop   ; Branch to LoadSpritesLoop if compare was Not Equal to zero
                         ; if compare was equal to 16, keep going down

   LDA #$00
   STA $2003             ; set the low byte (00) of the RAM address
   LDA #$02
   STA $4014             ; set the high byte (02) of the RAM address, start the transfer

LoadBackground:
   LDA $2002             ; read PPU status to reset the high/low latch
   LDA #$2C
   STA $2006             ; write the high byte of $2000 address
   LDA #$00
   STA $2006             ; write the low byte of $2000 address
   LDX #$00              ; start out at 0
LoadBackgroundLoop:
   LDA background, x     ; load data from address (background + the value in x)
   STA $2007             ; write to PPU
   INX                   ; X = X + 1
   CPX #$80              ; Compare X to hex $80, decimal 128 - copying 128 bytes
   BNE LoadBackgroundLoop  ; Branch to LoadBackgroundLoop if compare was Not Equal to zero
                         ; if compare was equal to 128, keep going down

LoadAttribute:
   LDA $2002             ; read PPU status to reset the high/low latch
   LDA #$23
   STA $2006             ; write the high byte of $23C0 address
   LDA #$C0
   STA $2006             ; write the low byte of $23C0 address
   LDX #$00              ; start out at 0
LoadAttributeLoop:
   LDA attribute, x      ; load data from address (attribute + the value in x)
   STA $2007             ; write to PPU
   INX                   ; X = X + 1
   CPX #$08              ; Compare X to hex $08, decimal 8 - copying 8 bytes
   BNE LoadAttributeLoop ; Branch to LoadAttributeLoop if compare was Not Equal to zero
                         ; if compare was equal to 128, keep going down

palette:
   .db $22,$29,$1A,$0F,  $22,$36,$17,$0F,  $22,$30,$21,$0F,  $22,$27,$17,$0F   ;;background palette
   .db $22,$1C,$15,$14,  $22,$02,$38,$3C,  $22,$1C,$15,$14,  $22,$02,$38,$3C   ;;sprite palette

sprites:
       ;vert tile attr horiz
   .db $80, $32, $00, $80   ;sprite 0
   .db $80, $33, $00, $88   ;sprite 1
   .db $88, $34, $00, $80   ;sprite 2
   .db $88, $35, $00, $88   ;sprite 3


background:
   .db $24,$24,$24,$24,$24,$24,$24,$24,$24,$24,$24,$24,$24,$24,$24,$24  ;;row 1
   .db $24,$24,$24,$24,$24,$24,$24,$24,$24,$24,$24,$24,$24,$24,$24,$24  ;;all sky

   .db $24,$24,$24,$24,$24,$24,$24,$24,$24,$24,$24,$24,$24,$24,$24,$24  ;;row 2
   .db $24,$24,$24,$24,$24,$24,$24,$24,$24,$24,$24,$24,$24,$24,$24,$24  ;;all sky

   .db $24,$24,$24,$24,$45,$45,$24,$24,$45,$45,$45,$45,$45,$45,$24,$24  ;;row 3
   .db $24,$24,$24,$24,$24,$24,$24,$24,$24,$24,$24,$24,$53,$54,$24,$24  ;;some brick tops

   .db $24,$24,$24,$24,$47,$47,$24,$24,$47,$47,$47,$47,$47,$47,$24,$24  ;;row 4
   .db $24,$24,$24,$24,$24,$24,$24,$24,$24,$24,$24,$24,$55,$56,$24,$24  ;;brick bottoms

attribute:
   .db %00000000, %00010000, %01010000, %00010000, %00000000, %00000000, %00000000, %00110000


;----------------------------------------------------------------
; interrupt vectors
;----------------------------------------------------------------

   .org $fffa

   .dw NMI
   .dw Reset
   .dw IRQ

;----------------------------------------------------------------
; CHR-ROM bank
;----------------------------------------------------------------

   .incbin "tiles.chr"
