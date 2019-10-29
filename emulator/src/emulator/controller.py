# from emulator.window import Window

class Controller:
    def __init__(self, ppu, read_count=0, strobe=False):
        # self.screen = ppu.screen
        self.keys = [False] * 8
        self.read_count = read_count
        self.strobe = strobe

    def read(self):
        # print(self.read_count)
        if self.read_count >= 8:
            return 1

        if not self.strobe:
            key = self.keys[self.read_count]
            # self.read_count += 1 % 8
            return key
        else:
            pressed = pygame.key.get_pressed()
            print("pressed", pressed)

            self.keys[0] = pressed[pygame.K_3]      # A
            self.keys[1] = pressed[pygame.K_2]      # B
            self.keys[2] = pressed[pygame.K_TAB]    # Select
            self.keys[3] = pressed[pygame.K_RETURN] # Start
            self.keys[4] = pressed[pygame.K_w]      # Up
            self.keys[5] = pressed[pygame.K_s]      # Down
            self.keys[6] = pressed[pygame.K_a]      # Left
            self.keys[7] = pressed[pygame.K_d]      # Right

            key = self.keys[self.read_count]
            # self.read_count += 1 % 8
            return key



    def latch(self, value):
        self.strobe = value
        # self.read_count = 0
        pressed = pygame.key.get_pressed()
        print("pressed", pressed)
        self.keys[0] = pressed[pygame.K_3]      # A
        self.keys[1] = pressed[pygame.K_2]      # B
        self.keys[2] = pressed[pygame.K_TAB]    # Select
        self.keys[3] = pressed[pygame.K_RETURN] # Start
        self.keys[4] = pressed[pygame.K_w]      # Up
        self.keys[5] = pressed[pygame.K_s]      # Down
        self.keys[6] = pressed[pygame.K_a]      # Left
        self.keys[7] = pressed[pygame.K_d]      # Right