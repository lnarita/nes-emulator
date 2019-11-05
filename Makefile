#PY=./emulator/venv/bin/python
PY=python3

TST=./tst
RES=./res
BIN=./bin
LOG=./log
EXT=./emulator/ext
NES=./emulator/src/main.py
ASSMBLR_DIR=${EXT}/asm6/
CROSS_AS=${EXT}/asm6/asm6

all: ${BIN} ${LOG}

assemble: ${CROSS_AS}
	{	cd ${TST}; \
		for i in *.s; \
		do \
			bin=$${i%.s}; \
			../${CROSS_AS} $$i ../${BIN}/$$bin; \
		done \
	}

${BIN}:
	@mkdir -p ${BIN}

${LOG}:
	@mkdir -p ${LOG}

${CROSS_AS}:
	@echo "compiling asm6f..."
	{ cd ${ASSMBLR_DIR}; make all; }

test: ${BIN} ${LOG} assemble
	@{  echo "************************* Tests ******************************"; \
		test_failed=0; \
		test_passed=0; \
		for test in ${BIN}/*; do \
			result="${LOG}/$$(basename $$test).log"; \
			expected="${RES}/$$(basename $$test).r"; \
			printf "Running $$test: "; \
			${PY} ${NES} $$test > $$result 2>&1; \
			errors=`diff -y --suppress-common-lines $$expected $$result | grep '^' | wc -l`; \
			if [ "$$errors" -eq 0 ]; then \
				printf "\033[0;32mPASSED\033[0m\n"; \
				test_passed=$$((test_passed+1)); \
			else \
				printf "\033[0;31mFAILED [$$errors errors]\033[0m\n"; \
				test_failed=$$((test_failed+1)); \
			fi; \
		done; \
		echo "*********************** Summary ******************************"; \
		echo "- $$test_passed tests passed"; \
		echo "- $$test_failed tests failed"; \
		echo "**************************************************************"; \
	}

clean:
	rm -rf ${BIN}/* ${LOG}/* ${CROSS_AS}

.PHONY: res
res: ${BIN} ${LOG} assemble
	@{	for test in ${BIN}/*; do \
			expected="${RES}/$$(basename $$test).r"; \
			if [ ! -f "$$expected" ]; then \
			    touch "$$expected"; \
			fi; \
			printf "Running $$test\n"; \
			${PY} ${NES} $$test 2>&1 > $$expected; \
		done; \
	}
