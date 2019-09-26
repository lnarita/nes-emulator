from emulator.adressing import Indirect
from emulator.cpu import CPU
from emulator.memory import Memory, MemoryPositions


def test_Indirect():
    address_mode = Indirect
    cpu = CPU()
    memory = Memory(rom=[0x32, 0x02, 0x7F, 0x01, 0xFF, 0x07], ram=list(map(lambda x: x % 256, range(Memory.ram_size()))))

    cpu.inc_cycle_by(-cpu.cycle)
    cpu.inc_pc_by(-cpu.pc + MemoryPositions.PRG_ROM_START.start)
    address = address_mode.fetch_address(cpu, memory)
    assert address == (((0x0232 + 1) % 256) << 8 | (0x0232 % 256))
    assert cpu.cycle == 4

    cpu.inc_cycle_by(-cpu.cycle)
    address = address_mode.fetch_address(cpu, memory)

    assert address == (((0x017F + 1) % 256) << 8 | (0x017F % 256))
    assert cpu.cycle == 4

    cpu.inc_cycle_by(-cpu.cycle)
    address = address_mode.fetch_address(cpu, memory)

    assert address == (((0x07FF + 1) % 256) << 8 | (0x07FF % 256))
    assert cpu.cycle == 4
