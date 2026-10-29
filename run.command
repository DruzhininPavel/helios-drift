#!/bin/bash
# Двойной клик — собрать (если нужно) и запустить игру.
cd "$(dirname "$0")" || exit 1

APP="HeliosDrift.app"
BIN="$APP/Contents/MacOS/heliosdrift"

if [ ! -x "$BIN" ]; then
  if ! command -v go >/dev/null 2>&1; then
    echo "Go не найден. Установи: brew install go"
    read -n 1 -s -r -p "Нажми любую клавишу..."
    exit 1
  fi
fi

if command -v go >/dev/null 2>&1; then
  NEWER=$(find go -name '*.go' -newer "$BIN" -print -quit 2>/dev/null)
  if [ ! -x "$BIN" ] || [ -n "$NEWER" ]; then
    echo "Сборка..."
    ( cd go && go build -ldflags "-X heliosdrift/game.Version=dev $(date '+%d.%m %H:%M')" -o "../$BIN" . ) || {
      echo "Ошибка сборки."
      read -n 1 -s -r -p "Enter..."
      exit 1
    }
    codesign --force --sign - "$APP" >/dev/null 2>&1
  fi
fi

open -W "$APP"
