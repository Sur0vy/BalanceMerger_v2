# BalanceMerger_v2
Списание материалов. Обновленная версия на Go.

## Сборка из под MacOS
- MacOS: fyne package -os darwin -icon Icon.png
- Windows 32: env GOOS="windows" GOARCH="386"   CGO_ENABLED="1" CC="i686-w64-mingw32-gcc"   fyne package -os windows -icon Icon.png
(предварительно установить нужные пакеты:  brew install mingw-w64)
