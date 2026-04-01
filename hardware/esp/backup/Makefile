CXX      := g++
CXXFLAGS := -std=c++17 -Wall -Wextra -O2
LDFLAGS  := -lmosquitto -lpthread

# Adapter object files
OBJS := can_bus.o mqtt_client.o

# Static library
LIB := liblt.a

.PHONY: all clean example

all: $(LIB)

# Build static library from adapters
$(LIB): $(OBJS)
	ar rcs $@ $^

can_bus.o: can_bus.cpp can_bus.h can_protocol.h
	$(CXX) $(CXXFLAGS) -c -o $@ can_bus.cpp

mqtt_client.o: mqtt_client.cpp mqtt_client.h
	$(CXX) $(CXXFLAGS) -c -o $@ mqtt_client.cpp

# Build example (optional)
example: main $(LIB)

main: main.cpp $(LIB)
	$(CXX) $(CXXFLAGS) -o $@ main.cpp -L. -llt $(LDFLAGS)

clean:
	rm -f $(OBJS) $(LIB) main
