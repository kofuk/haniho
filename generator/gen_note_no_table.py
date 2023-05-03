import math

# Generate table to convert MIDI note number to frequency.
# It also prints MIDI note number and Yamaha key name for reference.

def main():
    print(r'package generator')
    print(r'var noteNo = []float64{')

    for d in range(128):
        freq = 2 ** ((d - 69)/12) * 440
        name = ['C', 'C#', 'D', 'D#', 'E', 'F', 'F#', 'G', 'G#', 'A', 'A#', 'B'][d % 12]
        octave = d // 12 - 2
        print('\t{:.4f},  // {}: {}{}'.format(freq, d, name, octave))

    print(r'}')

if __name__ == '__main__':
    main()
