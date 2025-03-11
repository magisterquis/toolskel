# Makefile
# Build txtar testdata
# By J. Stuart McMurray
# Created 20250310
# Last Modified 20250310

# Files which need to be built
MACROFILES != find . -name '*.m4'
SHMORE      = ../../shmore.subr

.m4: # This is silly

all::
.PHONY: all

# Apparently one can't add a dependency for an inference rule with only one
# suffix, so we use a loop.  This is silly.
.for FN in ${MACROFILES}
${FN:R}: ${SHMORE} ${FN}
	m4 -PEE -Dm4_shmore=${SHMORE:Q} ${>:M*.m4} >$@.tmp
	mv $@.tmp $@
all:: ${FN:R}
.endfor

