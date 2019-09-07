import argparse

from constants import OpsCodes
from states import BaseCPUState


def main(args):
    emulate(args.file)


def emulate(file_path):
    cartridge_content = read_file(file_path)
    cpu_state = BaseCPUState(0xBBBB, 0xCCCC, 0xFF, 0xEE, 0xDD, 0xAA)
    print(cpu_state)
    cpu_state = BaseCPUState(0xBBBB, 0xCCCC, 0xFF, 0xEE, 0xDD, 0xAA, 0xFFFF, 0x99)
    print(cpu_state)

    while (True):
        instruction = read_mem(cpu_state.pc)
        decoded = decode_instruction(instruction)
        print_debug_line(cpu_state)
        break


def read_file(file_path):
    pass


def read_mem(pc):
    pass


def decode_instruction(instruction):
    pass


def print_debug_line(cpu_state):
    pass


if __name__ == "__main__":
    """ This is executed when run from the command line """
    parser = argparse.ArgumentParser()

    # Required positional argument
    parser.add_argument("file", help="NES Cartridge file path")

    args = parser.parse_args()
    main(args)
