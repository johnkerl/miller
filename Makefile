all: c
c: .always
	make -C c top
doc: .always
	cd doc && poki
install: .always
	make -C c install
.always:
	@true
