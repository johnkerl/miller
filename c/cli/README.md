# Miller command-line interface

mlrcli.h/c bags up all command-line options

I use argparse.h/c in place of getopt in order not to depend on GNU-isms, yet,
elsewhere in Miller I readily depend on GNU-isms such as getdelim().
