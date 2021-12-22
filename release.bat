@echo off
go build -ldflags -H=windowsgui
ResourceHacker -open karuba-tile-picker.exe -save KarubaTileRandomizer.exe -action addskip -res assets/images/icon.ico -mask ICONGROUP,MAIN,
xcopy /s /y assets\* D:\Programos\custom\karuba-tile-randomizer\assets\
xcopy /y KarubaTileRandomizer.exe D:\Programos\custom\karuba-tile-randomizer\