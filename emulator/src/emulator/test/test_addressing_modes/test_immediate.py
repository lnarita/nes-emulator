from emulator.adressing import Immediate
from emulator.cpu import CPU
from emulator.memory import Memory


def test_Immediate():
     address_mode = Immediate
     cpu = CPU()
     memory = Memory(rom=[])
