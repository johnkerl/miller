EXTRA_DIST=	mlr.1
man1_MANS=	mlr.1

# This target is not for use at package install; only at package-creation time.
# The mlr.1 file should be committed to source control. It uses the mlr
# executable from the PATH variable.
#
# See also
# https://miller.readthedocs.io/en/latest/build.html#creating-a-new-release-for-developers
manpage:
	ruby mkman.rb > mlr.1
