from more_itertools import flatten

from emulator.opcodes.arithmetic import ArithmeticAndLogicalOpCodes
from emulator.opcodes.flag import FlagOpCodes
from emulator.opcodes.jump import JumpOpCodes
from emulator.opcodes.move import MoveOpCodes
from emulator.opcodes.unofficial import UnofficialOpcodes


class OpCodes:
    types = [
        ArithmeticAndLogicalOpCodes,
        MoveOpCodes,
        JumpOpCodes,
        FlagOpCodes,
        UnofficialOpcodes
    ]

    all = dict(flatten(
        map(lambda x: list(x.all_commands()), types)
    ))
