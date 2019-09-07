from enum import unique, auto, Enum


# TODO: this is a WIP, it should probably not be designed like this

@unique
class MemoryPositions(Enum):
    INITIAL = (auto(), 0x0000, 0x0000)
    ZERO_PAGE = (auto(), 0x0000, 0x00FF)
    STACK = (auto(), 0x0100, 0x01FF)
    RAM = (auto(), 0x0200, 0x07FF)

    def __init__(self, value, start, end):
        self._value_ = value
        self.start = start
        self.end = end


class OpsCodes(Enum):
    # Logical and arithmetic
    ORA = (auto(),)  # A:=A or {adr}
    AND = (auto(),)  # A:=A&{adr}
    EOR = (auto(),)  # A:=A
    ADC = (auto(),)  # A:=A+{adr}
    SBC = (auto(),)  # A:=A-{adr}
    CMP = (auto(),)  # A-{adr}
    CPX = (auto(),)  # X-{adr}
    CPY = (auto(),)  # Y-{adr}
    DEC = (auto(),)  # {adr}:={adr}-1
    DEX = (auto(),)  # X:=X-1
    DEY = (auto(),)  # Y:=Y-1
    INC = (auto(),)  # {adr}:={adr}+1
    INX = (auto(),)  # X:=X+1
    INY = (auto(),)  # Y:=Y+1
    ASL = (auto(),)  # {adr}:={adr}*2
    ROL = (auto(),)  # {adr}:={adr}*2+C
    LSR = (auto(),)  # {adr}:={adr}/2
    ROR = (auto(),)  # {adr}:={adr}/2+C*128
    # Move
    LDA = (auto(),)  # A:={adr}
    STA = (auto(),)  # {adr}:=A
    LDX = (auto(),)  # X:={adr}
    STX = (auto(),)  # {adr}:=X
    LDY = (auto(),)  # Y:={adr}
    STY = (auto(),)  # {adr}:=Y
    TAX = (auto(),)  # X:=A
    TXA = (auto(),)  # A:=X
    TAY = (auto(),)  # Y:=A
    TYA = (auto(),)  # A:=Y
    TSX = (auto(),)  # X:=S
    TXS = (auto(),)  # S:=X
    PLA = (auto(),)  # A:=+(S)
    PHA = (auto(),)  # (S)-:=A
    PLP = (auto(),)  # P:=+(S)
    PHP = (auto(),)  # (S)-:=P
    # Jump/Flag
    BPL = (auto(),)  # branch on N=0
    BMI = (auto(),)  # branch on N=1
    BVC = (auto(),)  # branch on V=0
    BVS = (auto(),)  # branch on V=1
    BCC = (auto(),)  # branch on C=0
    BCS = (auto(),)  # branch on C=1
    BNE = (auto(),)  # branch on Z=0
    BEQ = (auto(),)  # branch on Z=1
    BRK = (auto(),)  # (S)-:=PC,P PC:=($FFFE)
    RTI = (auto(),)  # P,PC:=+(S)
    JSR = (auto(),)  # (S)-:=PC PC:={adr}
    RTS = (auto(),)  # PC:=+(S)
    JMP = (auto(),)  # PC:={adr}
    BIT = (auto(),)  # N:=b7 V:=b6 Z:=A&{adr}
    CLC = (auto(),)  # C:=0
    SEC = (auto(),)  # C:=1
    CLD = (auto(),)  # D:=0
    SED = (auto(),)  # D:=1
    CLI = (auto(),)  # I:=0
    SEI = (auto(),)  # I:=1
    CLV = (auto(),)  # V:=0
    NOP = (auto(),)
