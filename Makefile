MAIN_PKG = src/cmd/main.go

BUILD_DIR = dist

serve:
	go build -o ${BUILD_DIR}/main.exe ${MAIN_PKG}
	${BUILD_DIR}/main.exe