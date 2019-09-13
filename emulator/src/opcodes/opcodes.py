from more_itertools import flatten

from opcodes.arithmetic import ArithmeticAndLogicalOpCodes
from opcodes.flag import FlagOpCodes
from opcodes.jump import JumpOpCodes
from opcodes.move import MoveOpCodes


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
