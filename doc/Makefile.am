EXTRA_DIST=	miller.1.txt miller.1
man1_MANS=	miller.1

# Temporary for work on https://github.com/johnkerl/miller/issues/9:
# In Travis build:
# $ a2x -d manpage -f manpage miller.1.txt
# a2x: ERROR: "xmllint" --nonet --noout --valid "/path/to/pub_http_internet/miller/doc/miller.1.xml" returned non-zero exit status 4

MAINTAINERCLEANFILES=	miller.1
miller.1: miller.1.txt
	cp miller.1.premade miller.1
