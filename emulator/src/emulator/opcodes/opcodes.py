from more_itertools import flatten

from emulator.opcodes.arithmetic import ArithmeticAndLogicalOpCodes
from emulator.opcodes.flag import FlagOpCodes
from emulator.opcodes.jump import JumpOpCodes
from emulator.opcodes.move import MoveOpCodes


class OpCodes:
    types = [
        ArithmeticAndLogicalOpCodes,
        MoveOpCodes,
        JumpOpCodes,
        FlagOpCodes
    ]

    all = dict(flatten(
        map(lambda x: list(x.all_commands()), types)
    ))
