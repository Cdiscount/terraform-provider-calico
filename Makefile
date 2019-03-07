.PHONY: all build binaries test

GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

all: run

run:
	echo "Doing nothing"

build:
	mkdir -p ./bin
#	glide up
    # remove vendor from terraform (FIX for https://github.com/coreos/etcd/issues/9357)
	@rm -rf ./vendor/github.com/hashicorp/terraform/vendor/golang.org/x/net
	GOOS=linux go build -o ./bin/terraform-provider-calico .

test: fmt
	docker run -d \
 		-p 4001:4001 \
 		-p 2380:2380 \
 		-p 2379:2379 \
 		--name etcd_test \
 		 quay.io/coreos/etcd:v3.2.24 \
 		 etcd \
 			--name etcd_test \
 			--advertise-client-urls http://localhost:2379,http://localhost:4001 \
 			--listen-client-urls http://0.0.0.0:2379,http://0.0.0.0:4001 \
 			--initial-advertise-peer-urls http://localhost:2380 \
 			--listen-peer-urls http://0.0.0.0:2380 \
 			--initial-cluster-token etcd_test \
 			--initial-cluster etcd_test=http://172.17.0.2:2380 \
 			--initial-cluster-state new

	TF_ACC=1 CALICO_BACKEND_TYPE="etcdv3" CALICO_ETCD_ENDPOINTS="http://127.0.0.1:2379" go test -v ./calico -run="TestAcc"
	docker rm -f etcd_test

fmt:
	gofmt -w $(GOFMT_FILES)

clean:
	rm -rf ./bin