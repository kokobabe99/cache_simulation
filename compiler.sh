env GOOS="windows" GOARCH="amd64" CGO_ENABLED="1" CC="x86_64-w64-mingw32-gcc" go build


go build -o app main.go


env CGO_ENABLED="0" GOOS="linux" GOARCH="amd64" go build