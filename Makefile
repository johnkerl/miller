all: c
devall: c install doc
# TODO: the install target exists to put most-recent mlr executable in the
# path to be picked up by the mlr-execs in the docs dir. better would be to
# export PATH here with ./c at its head.
c: .always
	make -C c top
doc: .always
	cd doc && poki
install: .always
	make -C c install
clean: .always
	make -C c clean
.always:
	@true
