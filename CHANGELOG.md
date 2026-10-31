## Unreleased

### Added
- Every hull has its own cockpit. The Kestrel keeps the wide all-round view; the
  Adder is cramped with thick pillars and a strut across your sight; the Python
  has panoramic glass, a low dashboard and a centre console; the Asp is armoured
  with narrow slits and a gunsight column; the Fer-de-Lance sits low behind a
  V-brace. Same seat, different glass.
- `COCKPIT: ON / OFF` in the settings, remembered between runs. Off leaves the
  view clean, HUD and sight only.
- The title screen offers `HOW TO PLAY`, and pause opens the handbook with `H`.

### Fixed
- The pause list grew to fourteen rows and the last one fell off the plate.
- Lower-case letters rendered as question marks; the font now falls back to the
  capital.
- The Russian handbook was English: the translations were never added.
- Headless modes (`-shots`, `-bot`, `-sim`, `-probe`, ...) no longer rewrite the
  player's config. A shots run had switched the cockpit off and left the game
  in English because a scripted key press landed on a settings row.

## v1.14 — 2026-10-31

### Added
- Plasma accelerator, shield booster and heavy shielding, a second cargo bay and
  a heavy mining rig. The mining laser no longer takes the weapon slot, and the
  rig cuts more out of a rock.
- Contracts pay enough to matter: a few runs fund a real gun, not thirty.
- The last stage is a boss worth preparing for: 1400 hull, 180 shield and two
  escorts. Measured, a pulse laser needs 176 hits, a military one 61.
- `docs/progression.md`: what you earn, what it costs, and why the finale wants
  a better ship.
- `-ship` picks the hull for a bot run, so both paths can be measured.

### Changed
- Each hull now flies differently. The yard lists turning, shields and the
  heaviest gun the hull can carry, and those numbers are real: a Fer-de-Lance
  turns 2.4x faster than a Python, an Adder accelerates 1.8x harder than a
  Python and carries 60 instead of 140 of shield, and a hauler cannot mount a
  military laser. Buying a ship is a trade-off now, not just a bigger hold.
- `-bench-ships` flies every hull through the same fights and prints the table.

## v1.13.2 — 2026-10-30

### Changed
- README rewritten as one page: what the game is, which file to download, how
  to fly, the two modes. The JavaScript reference build is gone.
- The pre-rename save and config fallbacks are removed: nothing ever shipped
  under the old name.

### Fixed
- `PACKAGE-README.txt` inside the archives said `Shift+M` mutes; it is `F10`.

## v1.13.1 — 2026-10-30

### Added
- The world seed is picked at random per launch and printed in the corner of
  the title screen; `-seed N` plays that exact world again.

### Fixed
- Every launch played the same world: the seed defaulted to 1 after the
  generator became seeded. Reference screenshots and probes now pin seed 1 so
  they stay stable.
- README named the pre-rename save files.

## v1.13 — 2026-10-30

### Added
- Bot player that flies, trades, fights, jumps and docks headless, records its
  input and replays it exactly (`-bot`, `-record`, `-replay`, `-report`).
- Seeded world generator, `StateHash`, per-frame progress hash: a seed plus an
  input stream always produces the same run.
- Invariant checks each frame, stuck detection, stranded check.
- `-god` and `-ace` harness switches for measuring story completion and combat.
- `docs/bot-testing.md`: how to run it and what it checks.

### Fixed
- Crash when toggling the docking computer or jumping in-system with no station
  (witchspace): `autopilotState` and `JumpInSystem` dereferenced a nil station.
- A missed contract deadline could push cash below zero, leaving the commander
  unable to buy fuel.
- Replays drifted because recorded input was float32 against a float64 game.
- `-mode classic` had no effect: the game always started in Advanced.

# Изменения

Релиз выпускается не после каждой правки, а когда набралась связная порция
проверенных изменений. До этого всё живёт в разделе «Не выпущено».

## v1.12 — 2026-10-29

Управление геймпадом, рысканье и порядок в меню.

- **Русский интерфейс**: переведены названия действий в управлении (включая
  новое рысканье) и подписи экрана теста геймпада; в сообщении о привязке
  действие тоже называется по-русски. Данные, единицы и обозначения
  клавиш/кнопок остаются латиницей.

- **Выход из паузы** ведёт в главное меню (`Q`), а не закрывает приложение:
  закрытие игры осталось только на титульном экране.
- **«Продолжить полёт»** на титульном экране: возвращает в прерванную сессию,
  поэтому уход в меню больше не теряет прогресс. Пункт появляется только при
  активной сессии и исчезает после новой игры, загрузки сохранения или
  гибели корабля.

- **Знак рысканья**: первая версия крутила в обратную сторону — стик вправо
  уводил нос влево. Исправлено, направление закреплено тестами
  `TestYawTurnsRight` и `TestYawTurnsLeft`.
- **Рысканье** добавлено для Advanced: `Q`/`R` по умолчанию, назначается как
  обычное действие, работает с геймпада. В Classic ось отключена (в 1984 её
  не было) — игра один раз сообщает об этом.

- **Тест геймпада** переехал из настроек в конец экрана управления, оттуда
  видно и список привязок, и живые показания; `Esc` возвращает обратно.
- **В главном меню появился «Выход»** (четвёртый пункт и клавиша `4`).

- **Геймпад**: кнопки и оси переведены на стандартную раскладку Ebiten. Из-за
  этого D-pad в меню не работал вниз: номера кнопок были сырыми индексами
  устройства, а не позициями стандартной раскладки. Теперь D-pad в меню — это
  стрелки, а в полёте — четыре обычные кнопки (`D-UP`, `D-DOWN`, `D-LEFT`,
  `D-RIGHT`), которые можно повесить на любое действие.
- **Оси назначаются**: в привязках можно отдать действие стику или курку —
  в колонке `PAD` видно `LSX-`, `RSY+`, `LT`, `RT`. Стики остаются
  аналоговыми (мёртвая зона 0.15), поэтому скорость можно перенести на правый
  стик.
- **Оси в стандартной раскладке**: правый стик показывался как `LT`, потому
  что оси читались сырыми индексами SDL, а у части геймпадов порядок
  `LX, LY, LT, RX, RY, RT`. Теперь стики берутся через стандартные оси
  Ebiten, а курки стали обычными кнопками `LT`/`RT`.
- **Один слот на паде**: назначение оси или кнопки теперь снимает прежнюю
  привязку того же действия. Раньше `ACCELERATE` с осью правого стика
  продолжал ускоряться и с курка, а в колонке `PAD` показывалась кнопка,
  из-за чего казалось, что значение не поменялось.
- **Захват оси** переделан на простой и предсказуемый: клавиша и кнопка
  ловятся по свежему нажатию, а ось — по удержанию за порогом, без логики
  «пересечения». Раньше привязка то не срабатывала, то срабатывала позже,
  когда игрок уже ушёл на другую строку. Сообщение теперь называет действие:
  `BOUND: RSY- FOR ACCELERATE`.
- **Миграция привязок**: файл старого формата сохраняет клавиши и получает
  кнопки по умолчанию; в нём есть номер версии, так что дальше миграции не
  понадобятся.
- Документ «Изменения» и правило выпуска релизов.

## v1.11 — 2026-10-29

- **Экран управления**: вторая клавиша назначается по `Shift+Enter`,
  `Backspace` возвращает выбранной строке заводскую привязку.
- **Клавиша A**: в Ebiten `KeyA` — это ноль, а перебор клавиш начинался с
  `Key(1)`, плюс проверки `k > 0` считали ноль «пустым значением». Из-за
  этого `A` нельзя было ни назначить, ни увидеть в строке, ни нажать в
  полёте. Перебор идёт с нуля, «пусто» теперь `-1`.
- **Привязки**: назначенная клавиша забирается у прежнего действия, иначе
  старое поведение оставалось и срабатывали обе команды. `SOUND` переехал с
  `M` (делила клавишу с ракетой) на `F10`, старые файлы привязок чинятся при
  загрузке.
- **Версия в сборке**: подставляется из `git describe` или имени тега, видна
  в углу титульного экрана и печатается по `-version`.
- **macOS**: поставляется настоящим `HeliosDrift.app`; голый бинарник
  объясняет, почему окно не открывается, и называет правильный файл.
- **Русский интерфейс**: подписи, единицы, товары, снаряжение, корабли и
  ранги остаются латиницей, а меню, подсказки, сообщения и брифинги — на
  русском. Перевод теперь идёт по целым строкам и ведущим фразам, а не по
  словам: `HOLD 12/20 T` больше не превращается в `ТРЮМ 12/20 T`.
- **CI**: тесты и headless-прогон на Linux, Windows и macOS на каждый push;
  релиз собирается по тегу и прикладывает сборки для всех платформ.

## v1.0 — 2026-10-29

- Первая сборка: два режима (Classic и Advanced), кампания из восьми этапов,
  контракты, верфь, шахты, три типа станций, двадцать корпусов, русский и
  английский интерфейс, геймпад, три слота сохранений.
