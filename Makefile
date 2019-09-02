
build:
	go get -d ./...
	go get k8s.io/klog && cd $(GOPATH)/src/k8s.io/klog && git checkout v0.4.0
	go build main.go
	mv main kubectl-evict