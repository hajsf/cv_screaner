rsrc -manifest manifest.xml -ico=icon.ico -o rsrc.syso
go build -ldflags="-H windowsgui"