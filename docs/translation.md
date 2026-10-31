# Перевод: срез по игре

Всего строк интерфейса: **183**. Переведено **134**, не переведено **49**.

Правило: сокращения панели (`SPD`, `STN`, `CT`, `HOLD`), единицы (`T`, `LY`, `KG`), названия кораблей и систем остаются латиницей. Переводятся фразы и подписи экранов. Если строка попала в «не переведено», но по правилу должна быть по-русски — это и есть то, что надо поправить.

## Переведено

| EN | RU | где |
|---|---|---|
| `ACCEPTED` | ПРИНЯТЫ | screens.go |
| `ACCEPTED JOBS` | ПРИНЯТЫЕ ЗАДАНИЯ | screens.go |
| `ACTION` | ДЕЙСТВИЕ | screens.go |
| `ACTIVE: ` | АКТИВНО:  | screens.go |
| `ALREADY AT THE STATION` | УЖЕ У СТАНЦИИ | game.go |
| `ALREADY FITTED` | УЖЕ УСТАНОВЛЕНО | update.go |
| `ALREADY IN SYSTEM` | УЖЕ В ЭТОЙ СИСТЕМЕ | game.go |
| `ARRIVED: ` | ПРИБЫТИЕ:  | game.go |
| `ARROWS / MOUSE STEER   SPACE FIRE   X/Z SPEED` | СТРЕЛКИ / МЫШЬ   ПРОБЕЛ ОГОНЬ   X/Z СКОРОСТЬ | ui.go |
| `AVAILABLE` | ДОСТУПНЫ | screens.go |
| `BOUGHT ` | КУПЛЕНО  | update.go |
| `BOUND: ` | ПРИВЯЗАНО:  | screens.go |
| `BUTTONS 0-15` | КНОПКИ 0-15 | screens.go |
| `BUY` | КУПИТЬ | screens.go |
| `C CHART   J JUMP   ESC CONTROLS` | C КАРТА   J ПРЫЖОК   ESC ПАУЗА | ui.go |
| `CABIN TEMPERATURE CRITICAL` | ТЕМПЕРАТУРА КАБИНЫ КРИТИЧЕСКАЯ | update.go |
| `CAMPAIGN` | КАМПАНИЯ | screens.go |
| `CAMPAIGN COMPLETE` | КАМПАНИЯ ЗАВЕРШЕНА | screens.go |
| `CANNOT BUY` | КУПИТЬ НЕЛЬЗЯ | update.go |
| `CARGO CANISTER HIT HULL` | КОНТЕЙНЕР УДАРИЛ КОРПУС | update.go |
| `CARGO HOLD FULL` | ТРЮМ ПОЛОН | update.go |
| `COLLISION` | СТОЛКНОВЕНИЕ | update.go |
| `COMMANDER RESTORED` | ПИЛОТ ЗАГРУЖЕН | save.go |
| `COMMANDER SAVED` | ПИЛОТ СОХРАНЁН | screens.go |
| `COMMODITY` | ТОВАР | screens.go |
| `COMPLETE` | ЗАВЕРШЕНА | screens.go |
| `CONNECT ONE AND PRESS A BUTTON` | ПОДКЛЮЧИ И НАЖМИ КНОПКУ | screens.go |
| `CONTRACT ACCEPTED` | КОНТРАКТ ПРИНЯТ | contracts.go |
| `CONTRACT ACCEPTED - HUNT THE SHRIKE` | КОНТРАКТ ПРИНЯТ - ОХОТА НА ШРАЙКА | missions.go |
| `CONTRACT FAILED: CARGO MISSING` | КОНТРАКТ ПРОВАЛЕН: НЕТ ГРУЗА | campaign.go, contracts.go |
| `CONTRACT IN PROGRESS` | КОНТРАКТ ВЫПОЛНЯЕТСЯ | screens.go |
| `CONTRACT OFFER` | ПРЕДЛОЖЕНИЕ КОНТРАКТА | screens.go |
| `CONTROLS RESET` | УПРАВЛЕНИЕ СБРОШЕНО | screens.go |
| `DOCK AT THE DESTINATION TO DELIVER` | СТЫКУЙТЕСЬ С ЦЕЛЬЮ ДЛЯ ДОСТАВКИ | screens.go |
| `DOCKED` | СТЫКОВКА | game.go, screens.go |
| `DOCKING COMPUTER NOT FITTED` | НЕТ СТЫКОВОЧНОГО КОМПЬЮТЕРА | game.go |
| `DOCKING COMPUTER ON` | СТЫКОВКА АВТО ВКЛ | game.go |
| `E.C.M. ACTIVE` | ECM ВКЛЮЧЕН | game.go |
| `E.C.M. ACTIVE - MISSILES DESTROYED` | ECM: РАКЕТЫ УНИЧТОЖЕНЫ | game.go |
| `EMPTY` | ПУСТО | screens.go |
| `EN` | ЭН | ui.go |
| `ENERGY BOMB DETONATED  (` | ЭНЕРГОБОМБА ( | game.go |
| `ENTER ACCEPT     ESC DECLINE` | ENTER ПРИНЯТЬ     ESC ОТКАЗ | screens.go |
| `ESCAPE CAPSULE JETTISONED` | КАПСУЛА СБРОШЕНА | game.go |
| `ESCAPE CAPSULE LAUNCHED - RESCUED` | КАПСУЛА: СПАСЕНИЕ | game.go |
| `FIT A MINING LASER FIRST` | СНАЧАЛА БУРИЛЬНЫЙ ЛАЗЕР | update.go |
| `FITTED` | УСТАНОВЛЕНО | screens.go |
| `FUEL SCOOPING` | ЗАБОР ТОПЛИВА | update.go |
| `GALACTIC HYPERSPACE - GALAXY ` | ГАЛАКТИЧЕСКИЙ ПРЫЖОК - ГАЛАКТИКА  | missions.go |
| `GAME OVER` | ИГРА ОКОНЧЕНА | ui.go |
| `H - HYPERSPACE` | H - ГИПЕРПРЫЖОК | screens.go |
| `HELIOS COMBINE IS FINISHED` | КОНСОРЦИУМ ХЕЛИОС РАЗБИТ | campaign.go |
| `HOSTILE CONTACT` | ВРАГ НА РАДАРЕ | game.go |
| `HYPERSPACE` | ГИПЕРПРЫЖОК | ui.go |
| `HYPERSPACE MISJUMP - WITCHSPACE` | СБОЙ ПРЫЖКА - ВЕДЬМИНО ПРОСТРАНСТВО | game.go |
| `HYPERSPACE TO ` | ПРЫЖОК В  | game.go |
| `IN-SYSTEM JUMP` | ПРЫЖОК В СИСТЕМЕ | game.go |
| `INCOMING MISSILE` | РАКЕТА НА ПОДЛЁТЕ | game.go |
| `INSUFFICIENT FUEL` | НЕ ХВАТАЕТ ТОПЛИВА | game.go |
| `JUMP TO THE SYSTEM AND DESTROY IT` | ПРЫГНИТЕ В СИСТЕМУ И УНИЧТОЖЬТЕ ЕГО | screens.go |
| `KEY` | КЛАВИША | screens.go |
| `LASER DOWNGRADED FOR THIS HULL` | ЛАЗЕР СНЯТ: СЛАБЫЙ КОРПУС | shipyard.go |
| `LAUNCHED FROM ` | СТАРТ СО СТАНЦИИ  | game.go |
| `LOAD FAILED` | ЗАГРУЗКА НЕ УДАЛАСЬ | screens.go |
| `MANUAL CONTROL` | РУЧНОЕ УПРАВЛЕНИЕ | game.go |
| `MISSILE AWAY` | РАКЕТА ПУЩЕНА | game.go |
| `MISSILE RACKS FULL` | РАКЕТЫ НЕ ВЛЕЗАЮТ | update.go |
| `MISSION` | ЗАДАНИЕ | screens.go |
| `NO CONTRACT ON OFFER` | КОНТРАКТОВ НЕТ | screens.go |
| `NO CONTRACTS ON THE BOARD` | НА ДОСКЕ НЕТ КОНТРАКТОВ | screens.go |
| `NO E.C.M. FITTED` | НЕТ СИСТЕМЫ ECM | game.go |
| `NO ENERGY BOMB` | НЕТ ЭНЕРГОБОМБЫ | game.go |
| `NO GALACTIC HYPERDRIVE` | НЕТ ГАЛАКТИЧЕСКОГО ГИПЕРДРАЙВА | missions.go |
| `NO GAMEPAD CONNECTED` | ГЕЙМПАД НЕ ПОДКЛЮЧЁН | screens.go |
| `NO MISSILES` | РАКЕТ НЕТ | game.go |
| `NO SAVED COMMANDER` | НЕТ СОХРАНЁННОГО ПИЛОТА | screens.go |
| `NO STANDARD MAPPING` | НЕТ СТАНДАРТНОЙ РАСКЛАДКИ | screens.go |
| `NO TARGET LOCKED` | НЕТ ЗАХВАТА | game.go |
| `NONE` | НЕТ | screens.go |
| `NONE ACTIVE` | НЕТ АКТИВНОГО | screens.go |
| `NONE RIGHT NOW` | ПОКА НЕТ | screens.go |
| `NOT AVAILABLE HERE` | ЗДЕСЬ НЕ ПРОДАЁТСЯ | update.go |
| `NOT AVAILABLE IN CLASSIC MODE` | НЕДОСТУПНО В РЕЖИМЕ КЛАССИКИ | contracts.go, game.go, mechanics.go, shipyard.go |
| `NOT ENOUGH CARGO SPACE` | МАЛО МЕСТА В ТРЮМЕ | campaign.go, contracts.go |
| `NOT ENOUGH CREDITS` | НЕ ХВАТАЕТ КРЕДИТОВ | shipyard.go, update.go |
| `NOT ENOUGH FUEL FOR JUMP` | МАЛО ТОПЛИВА ДЛЯ ПРЫЖКА | game.go |
| `NOTHING TO JETTISON` | НЕЧЕГО СБРАСЫВАТЬ | mechanics.go |
| `NOTHING TO REPORT` | ДОКЛАДЫВАТЬ НЕЧЕГО | screens.go |
| `NOTHING TO SELL` | НЕЧЕГО ПРОДАВАТЬ | update.go |
| `OK` | ГОТОВО | screens.go |
| `OUT OF RANGE` | ВНЕ ДОСЯГАЕМОСТИ | screens.go |
| `PAD` | ГЕЙМПАД | screens.go |
| `PAD TEST` | ТЕСТ ГЕЙМПАДА | screens.go |
| `PARTIAL REFUEL` | ЧАСТИЧНАЯ ЗАПРАВКА | update.go |
| `POLICE SHIP DESTROYED  +50 CR BOUNTY ON YOU` | ПОЛИЦЕЙСКИЙ УНИЧТОЖЕН: +50 КР НАГРАДА | game.go |
| `POLICE VIPER ON INTERCEPT` | ПОЛИЦЕЙСКИЙ ВАЙПЕР | game.go |
| `PRESS 2ND` | ВТОРАЯ | screens.go |
| `PRESS ENTER TO CONTINUE` | НАЖМИ ENTER | ui.go |
| `PRESS KEY` | НАЖМИ КЛАВИШУ | screens.go |
| `PURCHASED: ` | КУПЛЕНО:  | update.go |
| `QTY` | КОЛ-ВО | screens.go |
| `RAW AXES` | СЫРЫЕ ОСИ | screens.go |
| `REFUEL  [F]` | ЗАПРАВКА [F] | screens.go |
| `REFUELLED FOR ` | ЗАПРАВЛЕНО ЗА  | update.go |
| `RESET DEFAULTS` | СБРОС ПО УМОЛЧАНИЮ | screens.go |
| `RESTORED: ` | ВОССТАНОВЛЕНО:  | screens.go |
| `ROCK BREAKS UP` | ГЛЫБА РАЗБИТА | mining.go |
| `SAVE FAILED` | ОШИБКА СОХРАНЕНИЯ | screens.go |
| `SCOOPED: ` | ПОДОБРАНО:  | update.go |
| `SECOND KEY: ` | ВТОРАЯ КЛАВИША:  | screens.go |
| `SECOND SLOT TAKES A KEY` | ВО ВТОРОЙ СЛОТ ТОЛЬКО КЛАВИША | screens.go |
| `SELL` | ПРОДАТЬ | screens.go |
| `SELL CARGO FIRST` | СНАЧАЛА ПРОДАЙТЕ ГРУЗ | shipyard.go |
| `SHIFT+M - SOUND ON/OFF` | SHIFT+M - ЗВУК ВКЛ/ВЫКЛ | screens.go |
| `SHRIKE DETECTED` | ШРАЙК ОБНАРУЖЕН | missions.go |
| `SLOT EMPTY` | СЛОТ ПУСТ | screens.go |
| `SOLD ` | ПРОДАНО  | update.go |
| `STAGE COMPLETE` | ЭТАП ПРОЙДЕН | campaign.go |
| `STD AXES` | СТАНДАРТНЫЕ | screens.go |
| `STORY CONTRACT` | СЮЖЕТНЫЙ КОНТРАКТ | screens.go |
| `STRIKE ORDER ACCEPTED` | ПРИКАЗ НА УДАР ПРИНЯТ | missions.go |
| `SWARM ATTACK` | АТАКА РОЯ | game.go |
| `SWARM MOTHERSHIP DETECTED` | ОБНАРУЖЕН МАТКОРАБЛЬ РОЯ | missions.go |
| `TANKS ALREADY FULL` | БАКИ ПОЛНЫ | update.go |
| `TARGET OUT OF RANGE (MAX 7.0 LY)` | ЦЕЛЬ ВНЕ ДОСЯГАЕМОСТИ (7 СГ) | game.go |
| `TECH ` | ТЕХ  | screens.go |
| `THE GALAXY IS YOURS, COMMANDER` | ГАЛАКТИКА ТВОЯ, ПИЛОТ | ui.go |
| `THE MODE IS FIXED FOR THIS COMMANDER` | РЕЖИМ ФИКСИРУЕТСЯ ДЛЯ ЭТОГО ПИЛОТА | screens.go |
| `THIS HULL CANNOT MOUNT IT` | ЭТОТ КОРПУС НЕ ПОТЯНЕТ | update.go |
| `TOO MANY CONTRACTS ALREADY` | СЛИШКОМ МНОГО КОНТРАКТОВ | contracts.go |
| `TRADER DESTROYED  LEGAL RECORD WORSENED` | ТОРГОВЕЦ УНИЧТОЖЕН: СТАТУС ХУЖЕ | game.go |
| `WARNING: ATMOSPHERE` | ВНИМАНИЕ: АТМОСФЕРА | update.go |
| `YOU ARE ALREADY FLYING ONE` | У ВАС УЖЕ ТАКОЙ КОРАБЛЬ | shipyard.go |
| `YOU NEED A GALACTIC HYPERDRIVE FIRST` | СНАЧАЛА НУЖЕН ГАЛАКТИЧЕСКИЙ ГИПЕРДРАЙВ | missions.go |

## Не переведено

| EN | где |
|---|---|
| `AS` | ui.go |
| `BOUNTY ` | game.go |
| `BOUNTY CONTRACT PAID  +` | contracts.go |
| `CARGO:` | screens.go |
| `CASH ` | ui.go |
| `CONTRACT ACCEPTED: ` | campaign.go |
| `CONTRACT EXPIRED  -` | contracts.go |
| `CT` | ui.go |
| `D-UP ` | screens.go |
| `DELIVERED  +` | campaign.go, contracts.go |
| `DIST ` | screens.go |
| `DOCUMENTS ABOARD - DELIVER TO ` | missions.go |
| `DOCUMENTS DELIVERED  +` | missions.go |
| `FEE: ` | screens.go |
| `FLYING ` | screens.go |
| `FS` | ui.go |
| `FU` | ui.go |
| `FUEL ` | screens.go |
| `GALAXY ` | screens.go |
| `HOLD` | screens.go |
| `HOLD ` | screens.go, ui.go |
| `HT` | ui.go |
| `JETTISONED ` | mechanics.go |
| `KG` | screens.go |
| `LSR` | ui.go |
| `LT ` | screens.go |
| `MAIL DELIVERED  +` | contracts.go |
| `MSL` | ui.go |
| `NEEDS A LARGE CARGO BAY FIRST` | update.go |
| `NEEDS A SHIELD BOOSTER FIRST` | update.go |
| `NO STATION IN THIS SYSTEM` | game.go |
| `NO STATION TO DOCK WITH` | game.go |
| `NO YAW IN CLASSIC FLIGHT` | update.go |
| `NOW FLYING ` | shipyard.go |
| `PASSENGERS DELIVERED  +` | contracts.go |
| `PAX ` | screens.go |
| `PAYS ` | screens.go |
| `POP ` | screens.go |
| `PROD ` | screens.go |
| `RADIUS ` | screens.go |
| `SHIELD ` | screens.go |
| `SPD ` | ui.go |
| `STAGE ` | screens.go |
| `STAGE COMPLETE  +` | campaign.go |
| `STN ` | ui.go |
| `SYS ` | ui.go |
| `TECH LEVEL ` | screens.go |
| `TO ` | ui.go |
| `TRADE IN ` | screens.go |

