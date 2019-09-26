from emulator.adressing import Immediate
from emulator.cpu import CPU
from emulator.memory import Memory, MemoryPositions


def test_Immediate():
     address_mode = Immediate
     cpu = CPU()
     memory = Memory(rom=[0x00])

     cpu.inc_cycle_by(-cpu.cycle)
     address = address_mode.fetch_address(cpu, memory)
     value = address_mode.read_from(cpu, memory, address)

     assert address == 0x00
     assert value == 0x00
     assert cpu.cycle == 1

     cpu.inc_cycle_by(-cpu.cycle)
     cpu.inc_pc_by(-cpu.pc)
     cpu.inc_pc_by(MemoryPositions.PRG_ROM_START.start)
     memory.rom[0] = 0x10
     address = address_mode.fetch_address(cpu, memory)
     value = address_mode.read_from(cpu, memory, address)

     assert address == 0x10
     assert value == 0x10
     assert cpu.cycle == 1

     cpu.inc_cycle_by(-cpu.cycle)
     cpu.inc_pc_by(-cpu.pc)
     cpu.inc_pc_by(MemoryPositions.PRG_ROM_START.start)
     memory.rom[0] = 0xFF
     address = address_mode.fetch_address(cpu, memory)
     value = address_mode.read_from(cpu, memory, address)

     assert address == 0xFF
     assert value == 0xFF
     assert cpu.cycle == 1
