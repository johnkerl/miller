EXTRA_DIST=	mlr.1.txt mlr.1
man1_MANS=	mlr.1

# Temporary for work on https://github.com/johnkerl/miller/issues/9:
# In Travis build:
# $ a2x -d manpage -f manpage mlr.1.txt
# a2x: ERROR: "xmllint" --nonet --noout --valid "/path/to/pub_http_internet/miller/doc/mlr.1.xml" returned non-zero exit status 4

MAINTAINERCLEANFILES=	mlr.1
mlr.1: mlr.1.txt
	cp mlr.1.premade mlr.1
