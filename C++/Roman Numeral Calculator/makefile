COMP := g++
CFLAGS := -std=c++11 -Wall -g
TARGET := out.o

SRC := $(wildcard *.cpp)
OBJ := $(patsubst %.cpp,%.o,$(SRC))

all: $(TARGET)
rebuild: clean all

$(TARGET): $(OBJ)
	$(COMP) -o $@ $^

%.o: %.cpp
	$(COMP) $(CFLAGS) -c $<

clean:
	rm -rf $(TARGET) *.o


.PHONY: all	clean