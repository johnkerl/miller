#!/usr/bin/python

import os
import sys
import getopt
import re
import collections

# ================================================================
#  o inclflds a,x,b
#  o newflds '{$y:$x*$x, $z:$x/2, $n:-$z}'
#  o greprecs '$x <= 2 && $y eq "zebra"'
#
#  o tabular pretty-print
#  o mean
#  o sort

# absolute essentials:
# * RECORD-LEVEL:
#   k include/exclude fields
#   o new field as function of old
#   o vertical pretty-print
# * STREAM-LEVEL:
#   o include/exclude records
#   o sort
#   o summarizations: min, max, mean, count, sum, first, last
#   o tabular pretty-print


# ================================================================
def usage():
    print(
        "Usage: %s [options] {modulator-spec} {zero or more filenames}"
        % os.path.basename(sys.argv[0]),
        file=sys.stderr,
    )
    msg = """
Options:
  -R {rs}   Input/output record separator
  -F {fs}   Input/output field separator
  -P {ps}   Input/output key-value-pair separator
  -v {name=value} xxx needs more doc

  --idkvp  Input  format is delimited by IRS,IFS,IPS
  --odkvp  Output format is delimited by IRS,IFS,IPS
  --icsv   Input  format is delimited by IRS,IFS,IPS, with header line followed by data lines (e.g. CSV)
  --ocsv   Output format is delimited by IRS,IFS,IPS, with header line followed by data lines (e.g. CSV)
  --inidx  Input  format is implicitly integer-indexed (awk-style)
  --onidx  Output format is implicitly integer-indexed (awk-style)
  --ixtab  Input  format is transposed-tabular-pretty-print
  --oxtab  Output format is transposed-tabular-pretty-print
Modulator specs:
--cat
--tac
--cut
--cutx
--sortfields
--sortfieldsup
--sortfieldsdown
"""
    print(msg, file=sys.stderr)
    sys.exit(1)


# ----------------------------------------------------------------
def parse_command_line():
    namespace = set_up_namespace()
    rreader = None
    rwriter = None
    rmodulator = None

    try:
        optargs, non_option_args = getopt.getopt(
            sys.argv[1:],
            "R:F:P:v:h",
            [
                "help",
                "idkvp",
                "odkvp",
                "icsv",
                "ocsv",
                "inidx",
                "onidx",
                "ixtab",
                "oxtab",
                "cat",
                "tac",
                "cut=",
                "cutx=",
                "sortfields",
                "sortfieldsup",
                "sortfieldsdown",
            ],
        )

    except getopt.GetoptError as e:
        print(str(e))
        usage()
        sys.exit(1)

    for opt, arg in optargs:
        if opt == "-R":
            rs = arg
            namespace.put("ORS", namespace.put("IRS", rs))
        elif opt == "-F":
            fs = arg
            namespace.put("OFS", namespace.put("IFS", fs))
        elif opt == "-P":
            ps = arg
            namespace.put("OPS", namespace.put("IPS", ps))
        elif opt == "-v":
            kv = arg.split("=", 1)
            namespace.put(kv[0], kv[1])

        elif opt == "--idkvp":
            rreader = RecordReaderDefault(
                istream=sys.stdin,
                namespace=namespace,
                irs=namespace.get("IRS"),
                ifs=namespace.get("IFS"),
                ips=namespace.get("IPS"),
            )
        elif opt == "--odkvp":
            rwriter = RecordWriterDefault(
                ostream=sys.stdout,
                ors=namespace.get("ORS"),
                ofs=namespace.get("OFS"),
                ops=namespace.get("OPS"),
            )

        elif opt == "--icsv":
            rreader = RecordReaderHeaderFirst(
                istream=sys.stdin,
                namespace=namespace,
                irs=namespace.get("IRS"),
                ifs=namespace.get("IFS"),
            )
        elif opt == "--ocsv":
            rwriter = RecordWriterHeaderFirst(
                ostream=sys.stdout,
                ors=namespace.get("ORS"),
                ofs=namespace.get("OFS"),
            )

        elif opt == "--inidx":
            rreader = RecordReaderIntegerIndexed(
                istream=sys.stdin,
                namespace=namespace,
                irs=namespace.get("IRS"),
                ifs=namespace.get("IFS"),
            )
        elif opt == "--onidx":
            rwriter = RecordWriterIntegerIndexed(
                ostream=sys.stdout,
                ors=namespace.get("ORS"),
                ofs=namespace.get("OFS"),
            )

        # elif opt == '--ixtab':
        #   pass
        elif opt == "--oxtab":
            rwriter = RecordWriterVerticallyTabulated(
                ostream=sys.stdout
            )  # xxx args w/r/t/ RS/FS/PS?!?

        elif opt == "--cat":
            rmodulator = CatModulator()
        elif opt == "--tac":
            rmodulator = TacModulator()
        elif opt == "--cut":
            rmodulator = SelectFieldsModulator(arg.split(namespace.get("IFS")))
        elif opt == "--cutx":
            rmodulator = DeselectFieldsModulator(arg.split(namespace.get("IFS")))
        elif opt == "--cutx":
            rmodulator = DeselectFieldsModulator(arg.split(namespace.get("IFS")))
        elif opt == "--sortfields":
            rmodulator = SortFieldsInRecordModulator(True)
        elif opt == "--sortfieldsup":
            rmodulator = SortFieldsInRecordModulator(True)
        elif opt == "--sortfieldsdown":
            rmodulator = SortFieldsInRecordModulator(False)

        elif opt == "--help":
            usage()
        else:
            print('Unhandled option "%s".' % opt, file=sys.stderr)
            sys.exit(1)

    # xxx non_option_arg_count = len(non_option_args)

    if rreader is None:
        rreader = RecordReaderDefault(
            istream=sys.stdin,
            namespace=namespace,
            irs=namespace.get("IRS"),
            ifs=namespace.get("IFS"),
            ips=namespace.get("IPS"),
        )
    if rwriter is None:
        rwriter = RecordWriterDefault(
            ostream=sys.stdout,
            ors=namespace.get("ORS"),
            ofs=namespace.get("OFS"),
            ops=namespace.get("OPS"),
        )
    if rmodulator is None:
        rmodulator = CatModulator()

    return {
        "namespace": namespace,
        "rreader": rreader,
        "rwriter": rwriter,
        "rmodulator": rmodulator,
    }


def main():
    options = parse_command_line()

    # parse ARGV:
    # * --ifmt: dkvp,hdr1st,iidxed,align,xposealign
    # * --ofmt: dkvp,hdr1st,iidxed,align,xposealign
    # * which-control-language spec?!?
    # * modulators/script ... this is the key decision area for language(s) design.
    # * filenames

    rreader = options["rreader"]
    rmodulator = options["rmodulator"]
    rwriter = options["rwriter"]

    smodulator = StreamModulator()
    smodulator.modulate(rreader, rmodulator, rwriter)


# ================================================================
class MillerNamespace:
    def __init__(self):
        self.mapping = {}
        self.imapping = {}

    def get(self, name):
        return self.mapping[name]

    def iget(self, name):
        return self.imapping[name]

    def put(self, name, value):
        self.mapping[name] = value
        return value

    def iput(self, name, ivalue):
        self.imapping[name] = ivalue
        return ivalue


# ================================================================
class Record:
    # kvs is list of pair-lists. (xxx: do tuples work too?)
    def __init__(self, kvs=[]):
        self.fields = collections.OrderedDict()
        self.mput(kvs)

    def put(self, k, v):
        self.fields[k] = v

    def mput(self, kvs):
        for [k, v] in kvs:
            self.fields[k] = v

    def get(self, k):
        return self.fields[k]

    def has_key(self, k):
        return k in self.fields.keys()

    def get_field_names(self):
        return self.fields.keys()

    def get_pairs(self):
        return self.fields.items()

    def num_pairs(self):
        return len(self.fields.items())

    # xxx xref to record-formatter classes
    def __str__(self):
        return self.fields.__repr__

    def __repr__(self):
        return self.fields.__repr__


# ================================================================
# Each record is a sequence of fields delimited by FS, each of which is a
# key-value pair separated by PS.


class RecordReader:
    def __init__(self, istream, namespace, irs, ifs, ips):
        self.istream = istream
        self.namespace = namespace
        self.irs = irs
        self.ifs = ifs
        self.ips = ips


class RecordReaderDefault(RecordReader):
    def __init__(self, istream, namespace, irs, ifs, ips):
        RecordReader.__init__(self, istream, namespace, irs, ifs, ips)

    def read(self):
        line = self.istream.readline()  # xxx use self.irs
        if line == "":
            return None

        line = (
            line.strip()
        )  # Remove leading/trailing whitespace including carriage return from readline().
        fields = line.split(self.ifs)
        kvs = [field.split(self.ips, 1) for field in fields]
        record = Record(kvs)

        self.namespace.iput("NF", record.num_pairs)
        self.namespace.iput("NR", self.namespace.iget("NR") + 1)

        # xxx stub
        self.namespace.put("FILENAME", None)
        self.namespace.iput("FNR", self.namespace.iget("FNR") + 1)

        return record


# ----------------------------------------------------------------
# awk-style
class RecordReaderIntegerIndexed(RecordReader):
    # xxx ctor with istream context?!? or independent of that?!? for cskv, no matter.
    # csv reader of course needs context.
    def __init__(self, istream, namespace, irs, ifs):
        RecordReader.__init__(self, istream, namespace, irs, ifs, None)

    def read(self):
        # xxx use self.irs
        line = self.istream.readline()
        if line == "":
            return None
        line = (
            line.strip()
        )  # Remove leading/trailing whitespace including carriage return from readline().
        fields = re.split(self.ifs, line)
        kvs = []
        i = 1
        for field in fields:
            kvs.append([i, field])
            i += 1
        return Record(kvs)


# ----------------------------------------------------------------
# csv-style
class RecordReaderHeaderFirst(RecordReader):
    def __init__(self, istream, namespace, irs, ifs):
        RecordReader.__init__(self, istream, namespace, irs, ifs, None)
        self.field_names = None
        self.header_line = None

    def read(self):
        if not self.field_names:
            header_line = self.istream.readline()
            if header_line == "":
                return None
            # Remove leading/trailing whitespace including carriage return from readline().
            header_line = header_line.strip()
            self.field_names = header_line.split(self.ifs, -1)
            self.header_line = header_line

        data_line = self.istream.readline()
        if data_line == "":
            return None
        # Remove leading/trailing whitespace including carriage return from readline().
        data_line = data_line.strip()
        field_values = data_line.split(self.ifs, -1)
        if len(self.field_names) != len(field_values):
            raise Exception(
                'Header/data length mismatch: %d != %d in "%s" and "%s"'
                % (
                    len(self.field_names),
                    len(field_values),
                    self.header_line,
                    data_line,
                )
            )

        return Record(zip(self.field_names, field_values))


# ================================================================
# xxx ostream at ctor??  needs drain-at-end logic for prettyprint.


class RecordWriter:
    def __init__(self, ostream, ors, ofs, ops):
        self.ostream = ostream
        self.ors = ors
        self.ofs = ofs
        self.ops = ops


class RecordWriterDefault(RecordWriter):
    def __init__(self, ostream, ors, ofs, ops):
        RecordWriter.__init__(self, ostream, ors, ofs, ops)

    def write(self, record):
        self.ostream.write(
            self.ofs.join([str(k) + self.ops + str(v) for [k, v] in record.get_pairs()])
        )
        self.ostream.write("\n")


# ----------------------------------------------------------------
class RecordWriterHeaderFirst(RecordWriter):
    def __init__(self, ostream, ors, ofs):
        RecordWriter.__init__(self, ostream, ors, ofs, None)
        self.field_names = None

    def write(self, record):
        data_string = self.ofs.join([str(v) for [k, v] in record.get_pairs()])
        if self.field_names is None:
            self.field_names = record.get_field_names()
            header_string = self.ofs.join([str(k) for [k, v] in record.get_pairs()])
            self.ostream.write(header_string)
            self.ostream.write("\n")
        self.ostream.write(data_string)
        self.ostream.write("\n")


# ----------------------------------------------------------------
# xxx rename


class RecordWriterVerticallyTabulated(RecordWriter):
    def __init__(self, ostream):
        RecordWriter.__init__(self, ostream, None, None, None)

    def write(self, record):
        max_field_name_width = 1
        field_names = record.get_field_names()
        for field_name in field_names:
            field_name_width = len(field_name)
            if field_name_width > max_field_name_width:
                max_field_name_width = field_name_width
        lines = []
        for field_name in field_names:
            lines.append(
                "%-*s %s" % (max_field_name_width, field_name, record.get(field_name))
            )
        self.ostream.write("\n".join(lines))
        self.ostream.write("\n\n")


# ----------------------------------------------------------------
class RecordWriterIntegerIndexed:
    def __init__(self, ostream, ors, ofs):
        self.ostream = ostream
        self.ors = ors
        self.ofs = ofs

    def write(self, record):
        self.ostream.write(self.ofs.join([str(v) for [k, v] in record.get_pairs()]))
        self.ostream.write("\n")


# ================================================================
class CatModulator:
    def __init__(self):
        pass

    def modulate(self, record):
        if record is None:  # drain at end
            return []
        return [record]


class TacModulator:
    def __init__(self):
        self.records = []

    def modulate(self, record):
        if record is None:  # drain at end
            self.records.reverse()
            rv = self.records
            self.records = []
            return rv
        else:
            self.records.append(record)
            return []


class SelectFieldsModulator:
    def __init__(self, field_names):
        self.field_names = field_names

    def modulate(self, record):
        if record is None:  # drain at end
            return []
        kvs = []
        for field_name in self.field_names:
            if record.has_key(field_name):
                kvs.append((field_name, record.get(field_name)))
        new_record = Record()
        new_record.mput(kvs)
        return [new_record]


# The field_names argument may be a list or hash-set -- as long as it supports
# the "in" operator as in "name in field_names".
# xxx to do: use a hash-set internally.
class DeselectFieldsModulator:
    def __init__(self, field_names):
        self.field_names = field_names

    def modulate(self, record):
        if record is None:  # drain at end
            return []
        kvs = []
        for field_name in record.get_field_names():
            if field_name not in self.field_names:
                kvs.append((field_name, record.get(field_name)))
        new_record = Record()
        new_record.mput(kvs)
        return [new_record]


class SortFieldsInRecordModulator:
    def __init__(self, do_ascending_sort=True):
        self.do_ascending_sort = do_ascending_sort

    def modulate(self, record):
        if record is None:  # drain at end
            return []
        kvs = []
        sorted_field_names = sorted(record.get_field_names())
        if not self.do_ascending_sort:
            sorted_field_names.reverse()  # xxx optimize
        for field_name in sorted_field_names:
            kvs.append((field_name, record.get(field_name)))
        new_record = Record()
        new_record.mput(kvs)
        return [new_record]


class MeanKeeper:
    def __init__(self):
        self.sum = 0.0
        self.count = 0

    def put(self, x):
        self.sum += x
        self.count += 1

    def get_sum(self):
        return self.sum

    def get_count(self):
        return self.count

    def get_mean(self):
        # In IEEE-standard floating-point this would give NaN in the empty case.
        # But Python throws an exception on divide by zero instead.
        if self.count == 0:
            return None
        else:
            return self.sum / self.count


class MeanModulator:
    def __init__(self, collate_field_names, key_field_names=[]):
        self.collate_field_names = collate_field_names
        self.key_field_names = key_field_names
        # map from key-field values to (map from collate-field names to MSCKeeper objects).
        self.collate_outputs = {}

    def modulate(self, record):
        if record is not None:  # drain at end
            # xxx optimize
            for value_field_name in self.collate_field_names:
                if not record.has_key(value_field_name):
                    return []
            for key_field_name in self.key_field_names:
                if not record.has_key(key_field_name):
                    return []

            collate_field_values = [
                float(record.get(k)) for k in self.collate_field_names
            ]
            key_string = ",".join([record.get(k) for k in self.key_field_names])

            # xxx wip
            return []
        else:
            # xxx stub
            output_record = Record()
            output_record.put("foo", "bar")
            return [output_record]


# ================================================================
class StreamModulator:
    def __init__(self):
        pass

    def modulate(self, rreader, rmodulator, rwriter):
        while True:
            in_record = rreader.read()

            out_records = rmodulator.modulate(in_record)

            for out_record in out_records:
                rwriter.write(out_record)

            if in_record is None:
                break


# ================================================================
def set_up_namespace():
    namespace = MillerNamespace()
    namespace.put("ORS", namespace.put("IRS", "\n"))
    namespace.put("OFS", namespace.put("IFS", ","))
    namespace.put("OPS", namespace.put("IPS", "="))

    # xxx CONVFMT

    namespace.put("FILENAME", None)
    namespace.iput("NF", None)
    namespace.iput("NR", 0)
    namespace.iput("FNR", 0)

    return namespace


# ================================================================
main()
