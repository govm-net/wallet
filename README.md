# wallet

govm wallet(cmd,web assembly,html)

1. command:
    * cd bin/cmd
    * go build
    * [more info](bin/cmd/README.md)
2. web assembly:
    * cd bin/wasm
    * ./build.sh
    * [how to use](https://github.com/golang/go/wiki/WebAssembly)
3. local web wallet:
    * cd bin/html
    * ./build.sh
    * ./html
4. gui:
    * cd bin/gui
    * fyne.exe package -os windows -name govm.net -icon ./assets/govm.png
5. android library:
    * cd trans
    * gomobile bind -target=android
