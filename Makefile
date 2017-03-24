CXX=g++
CXFLAGS= -DSTANDALONE_BUILD=1 -std=c++11 # -g -O3 -Wall

CBDEPS_DIR=/Users/$(USER)/.cbdepscache/
DYLD_LIBRARY_PATH=/Users/$(USER)/.cbdepscache/lib
CMD_DIR=cmd/producer/

LDFLAGS=-luv -L$(CBDEPS_DIR)lib/ -ljemalloc -L$(CBDEPS_DIR)lib/ -lv8 -lcouchbase


SOURCES=v8_consumer/src/client.cc v8_consumer/src/commands.cc \
				v8_consumer/src/message.cc v8_consumer/src/v8worker.cc \
				v8_consumer/src/n1ql.cc v8_consumer/src/bucket.cc \
				v8_consumer/src/parse_deployment.cc \
				v8_consumer/src/log.cc

INCLUDE_DIRS=-I$(CBDEPS_DIR) -I$(CBDEPS_DIR)include -I v8_consumer/include/

OUT=$(CMD_DIR)client.bin

build:
	$(CBDEPS_DIR)bin/flatc -o flatbuf/include/ -c flatbuf/schema/*.fbs
	$(CBDEPS_DIR)bin/flatc -g flatbuf/schema/*.fbs
	$(CXX) $(CXFLAGS) $(SOURCES) $(INCLUDE_DIRS) $(LDFLAGS) -o $(OUT)
	cd $(CMD_DIR); go build -race; bash fix_rpath.sh

allopt:
	$(CXX) $(CXFLAGS) $(SOURCES) $(INCLUDE_DIRS) $(LDFLAGS) -O3 -o $(OUT)
	cd $(CMD_DIR); go build; bash fix_rpath.sh
	cd $(CMD_DIR); ./producer

run: build
	cd $(CMD_DIR); ./producer

clean:
	rm -rf $(OUT)
