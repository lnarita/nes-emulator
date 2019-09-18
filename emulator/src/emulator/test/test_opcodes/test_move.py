from emulator.cpu import CPU
from emulator.memory import Memory
import emulator.opcodes.move as move
from emulator.opcodes.move import MoveOpCodes
from emulator.constants import AddressingMode
from emulator.opcodes.opcodes import OpCodes

def test_TXA():
     cpu = CPU()
     memory = Memory()
     cpu.a = 1
     cpu.x = 2
     TXA = OpCodes.all[138]
     TXA.exec(cpu, memory)

     assert cpu.a == 2
     assert cpu.zero == False
     assert cpu.negative == False

     cpu.x = 0
     cpu.a = 1
     TXA.exec(cpu, memory)

     assert cpu.a == 0
     assert cpu.zero == True
     assert cpu.negative == False

     cpu.x = 0b10000001
     cpu.a = 1
     TXA.exec(cpu, memory)

     assert cpu.a == 0b10000001
     assert cpu.zero == False
     print(cpu.negative)
     assert cpu.negative == True

def test_TAX():
     cpu = CPU()
     memory = Memory()
     cpu.x = 1
     cpu.a = 2
     TAX = OpCodes.all[170]
     TAX.exec(cpu, memory)

     assert cpu.x == 2
     assert cpu.zero == False
     assert cpu.negative == False

     cpu.a = 0
     cpu.x = 1
     TAX.exec(cpu, memory)

     assert cpu.x == 0
     assert cpu.zero == True
     assert cpu.negative == False

     cpu.a = 0b10000001
     cpu.x = 1
     TAX.exec(cpu, memory)

     assert cpu.x == 0b10000001
     assert cpu.zero == False
     print(cpu.negative)
     assert cpu.negative == True