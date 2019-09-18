from emulator.cpu import CPU
from emulator.memory import Memory
import emulator.opcodes.move as move
from emulator.opcodes.move import MoveOpCodes
from emulator.constants import AddressingMode
from emulator.opcodes.opcodes import OpCodes

def test_TXA():
     print("lalala")
     cpu = CPU()
     memory = Memory()
     print(OpCodes.all[138].exec(cpu, memory))


def test_TSA():
     assert sum([1, 2, 2]) == 5
