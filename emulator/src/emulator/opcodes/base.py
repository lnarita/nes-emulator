class OpCode(object):
    def __init__(self, code, addressing_mode, cycles):
        self.id = code
        self.addressing_mode = addressing_mode
        self.cycles = cycles

    def exec(self, cpu, memory):
        raise NotImplementedError("Class {} not implemented yet! CPU: {}, MEM: {}".format(self, cpu, memory))

    @classmethod
    def create_variations(cls):
        pass

    @classmethod
    def create_dict_entry(cls, x):
        return tuple((x[0], cls(*x)))

    def __str__(self):
        def __str_addr():
            if self.addressing_mode is not None:
                if self.addressing_mode.low is not None and self.addressing_mode.high is not None:
                    return "{:02X} {:02X}".format(self.addressing_mode.low, (self.addressing_mode.high >> 8))
                elif self.addressing_mode.low is not None:
                    return "{:02X}".format(self.addressing_mode.low)
            return ""

        def __str_addr_2():
            if self.addressing_mode is not None:
                if self.addressing_mode.addr is not None and self.addressing_mode.data is not None:
                    return "{} {} {}".format(type(self).__name__, self.addressing_mode.addr, self.addressing_mode.data)
                elif self.addressing_mode.addr is not None:
                    return "{} {}".format(type(self).__name__, self.addressing_mode.addr)
            return "{}".format(type(self).__name__)

        return "{:02X} {:<6} {:<30}".format(self.id, __str_addr(), __str_addr_2())
