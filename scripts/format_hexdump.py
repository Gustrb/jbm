def main():
    with open("hexdump.txt") as f:
        # The file is in the following format:
        # aaaa bbbb cccc dddd
        # And we want to convert it to:
        # 0xAA, 0xAA, 0xBB, 0xBB, 0xCC, 0xCC, 0xDD, 0xDD

        col = 0
        for line in f:
            line = line.strip()
            bytes = line.split(" ")

            # split the bytes into pairs
            for byte in bytes:
                if col == 8:
                    print()
                    col = 0
                # the first two characters are the first byte
                first_byte = byte[0:2]
                # the last two characters are the second byte
                second_byte = byte[2:4]
                print(f"0x{second_byte.upper()}, 0x{first_byte.upper()}", end=", ")
                col += 2

        print()

if __name__ == "__main__":
    main()
