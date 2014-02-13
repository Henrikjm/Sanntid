SOURCES = $(wildcard *.c)
TARGETS = $(SOURCES:.c=)

all: $(TARGETS)

%: %.c
	gcc -g -o $@ $<

clean:
	rm -f $(SOURCES:.c=)

rebuild: clean all

.PHONY: rebuild clean all

