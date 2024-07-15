#!/bin/bash
CURRENT_DIR=$1
rm -rf ${CURRENT_DIR}/genproto
for x in $(find ${CURRENT_DIR}/ArtisanConnect_protos1/* -type d); do
  protoc -I=${x} -I=${CURRENT_DIR}/ArtisanConnect_protos1 -I /usr/local/go --go_out=${CURRENT_DIR} \
   --go-grpc_out=${CURRENT_DIR} ${x}/*.proto
done
