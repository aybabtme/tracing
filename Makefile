root_dir = $(shell git rev-parse --show-toplevel)
pkg_dir = $(root_dir:$(GOPATH)/src/%=%)
gen_dir = gen
flat_idl = idl/span.fbs

all: codegen

codegen: $(gen_dir)/span/Span.go

$(gen_dir)/span:
	mkdir -p $(gen_dir)/span

$(gen_dir)/span/Span.go: $(gen_dir)/span $(flat_idl)
	docker run -v $(root_dir):$(root_dir) -w $(root_dir) neomantra/flatbuffers flatc --grpc --go -o $(gen_dir)/span/ $(flat_idl)
