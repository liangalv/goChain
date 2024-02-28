#!/bin/zsh

# Set the directory where your .proto files are located
PROTO_DIR="./proto"
# Set the output directory for your Go files
GO_OUT_DIR="./core/types"

# Ensure the output directory exists
mkdir -p ${GO_OUT_DIR}

# Loop through all .proto files in the directory and compile them
for PROTO_FILE in ${PROTO_DIR}/*.proto; do
    echo "Compiling ${PROTO_FILE}..."
    protoc -I${PROTO_DIR} \
          --go_out=${GO_OUT_DIR} \
          --go_opt=paths=source_relative \
          --go-grpc_out=${GO_OUT_DIR} \
          --go-grpc_opt=paths=source_relative \
          ${PROTO_FILE}
    if [ $? -ne 0 ]; then
        echo "Failed to compile ${PROTO_FILE}"
        exit 1
    fi
done

echo "Compilation successful."
