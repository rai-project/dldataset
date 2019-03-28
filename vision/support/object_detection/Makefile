all: generate

fmt:
	go fmt ./...

install-deps:
	go get github.com/jteeuwen/go-bindata/...
	go get github.com/elazarl/go-bindata-assetfs/...

generate: clean generate-models

generate-proto:
	protoc --gogofaster_out=. -Iproto -I$(GOPATH)/src proto/label_map.proto

generate-assets:
	go-bindata -nomemcopy -prefix assets/ -pkg object_detection -o object_detection_static.go -ignore=.DS_Store -ignore=README.md assets/...

clean-assets:
	rm -fr object_detection_static.go

clean-proto:
	rm -fr *pb.go

clean: clean-models

travis: install-deps generate
	echo "building..."
	go build
