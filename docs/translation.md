# Перевод: срез по игре

Всего строк интерфейса найдено: **183**. Из них переведено **115**, не переведено **68**.

Ниже — полный список. Формат: английская строка (как она стоит в коде) → русская. В последней колонке — файл, где строка встречается. Правки удобно писать прямо в таблицу или отдельным списком.

Правило, по которому переводятся не всё: сокращения на панели (`SPD`, `STN`, `CT`, `HOLD`), единицы (`T`, `LY`, `KG`), названия кораблей и систем, цифры и суммы остаются латиницей. Переводятся целые фразы и подписи экранов. Если строка попала в «не переведено», но по правилу должна быть по-русски — это и есть то, что надо поправить.

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
| `C CHART   J JUMP   ESC CONTROLS` | C КАРТА   J ПРЫЖОК   ESC ПАУЗА | ui.go |
| `CABIN TEMPERATURE CRITICAL` | ТЕМПЕРАТУРА КАБИНЫ КРИТИЧЕСКАЯ | update.go |
| `CAMPAIGN` | КАМПАНИЯ | screens.go |
| `CAMPAIGN COMPLETE` | КАМПАНИЯ ЗАВЕРШЕНА | screens.go |
| `CANNOT BUY` | КУПИТЬ НЕЛЬЗЯ | update.go |
| `CARGO CANISTER HIT HULL` | КОНТЕЙНЕР УДАРИЛ КОРПУС | update.go |
| `CARGO HOLD FULL` | ТРЮМ ПОЛОН | update.go |
| `COMMANDER RESTORED` | ПИЛОТ ЗАГРУЖЕН | save.go |
| `COMMANDER SAVED` | ПИЛОТ СОХРАНЁН | screens.go |
| `COMPLETE` | ЗАВЕРШЕНА | screens.go |
| `CONNECT ONE AND PRESS A BUTTON` | ПОДКЛЮЧИ И НАЖМИ КНОПКУ | screens.go |
| `CONTRACT ACCEPTED` | КОНТРАКТ ПРИНЯТ | contracts.go |
| `CONTRACT ACCEPTED - HUNT THE SHRIKE` | КОНТРАКТ ПРИНЯТ - ОХОТА НА ШРАЙКА | missions.go |
| `CONTRACT FAILED: CARGO MISSING` | КОНТРАКТ ПРОВАЛЕН: НЕТ ГРУЗА | campaign.go, contracts.go |
| `CONTRACT IN PROGRESS` | КОНТРАКТ ВЫПОЛНЯЕТСЯ | screens.go |
| `CONTRACT OFFER` | ПРЕДЛОЖЕНИЕ КОНТРАКТА | screens.go |
| `CONTROLS RESET` | УПРАВЛЕНИЕ СБРОШЕНО | screens.go |
| `DOCK AT THE DESTINATION TO DELIVER` | СТЫКУЙТЕСЬ С ЦЕЛЬЮ ДЛЯ ДОСТАВКИ | screens.go |
| `DOCKING COMPUTER NOT FITTED` | НЕТ СТЫКОВОЧНОГО КОМПЬЮТЕРА | game.go |
| `DOCKING COMPUTER ON` | СТЫКОВКА АВТО ВКЛ | game.go |
| `E.C.M. ACTIVE` | ECM ВКЛЮЧЕН | game.go |
| `E.C.M. ACTIVE - MISSILES DESTROYED` | ECM: РАКЕТЫ УНИЧТОЖЕНЫ | game.go |
| `EN` | ЭН | ui.go |
| `ENERGY BOMB DETONATED  (` | ЭНЕРГОБОМБА ( | game.go |
| `ENTER ACCEPT     ESC DECLINE` | ENTER ПРИНЯТЬ     ESC ОТКАЗ | screens.go |
| `ESCAPE CAPSULE JETTISONED` | КАПСУЛА СБРОШЕНА | game.go |
| `ESCAPE CAPSULE LAUNCHED - RESCUED` | КАПСУЛА: СПАСЕНИЕ | game.go |
| `FUEL SCOOPING` | ЗАБОР ТОПЛИВА | update.go |
| `GALACTIC HYPERSPACE - GALAXY ` | ГАЛАКТИЧЕСКИЙ ПРЫЖОК - ГАЛАКТИКА  | missions.go |
| `GAME OVER` | ИГРА ОКОНЧЕНА | ui.go |
| `H - HYPERSPACE` | H - ГИПЕРПРЫЖОК | screens.go |
| `HOSTILE CONTACT` | ВРАГ НА РАДАРЕ | game.go |
| `HYPERSPACE` | ГИПЕРПРЫЖОК | ui.go |
| `HYPERSPACE TO ` | ПРЫЖОК В  | game.go |
| `IN-SYSTEM JUMP` | ПРЫЖОК В СИСТЕМЕ | game.go |
| `INCOMING MISSILE` | РАКЕТА НА ПОДЛЁТЕ | game.go |
| `INSUFFICIENT FUEL` | НЕ ХВАТАЕТ ТОПЛИВА | game.go |
| `JUMP TO THE SYSTEM AND DESTROY IT` | ПРЫГНИТЕ В СИСТЕМУ И УНИЧТОЖЬТЕ ЕГО | screens.go |
| `KEY` | КЛАВИША | screens.go |
| `LAUNCHED FROM ` | СТАРТ СО СТАНЦИИ  | game.go |
| `MANUAL CONTROL` | РУЧНОЕ УПРАВЛЕНИЕ | game.go |
| `MISSILE AWAY` | РАКЕТА ПУЩЕНА | game.go |
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
| `OUT OF RANGE` | ВНЕ ДОСЯГАЕМОСТИ | screens.go |
| `PAD` | ГЕЙМПАД | screens.go |
| `PAD TEST` | ТЕСТ ГЕЙМПАДА | screens.go |
| `PARTIAL REFUEL` | ЧАСТИЧНАЯ ЗАПРАВКА | update.go |
| `POLICE SHIP DESTROYED  +50 CR BOUNTY ON YOU` | ПОЛИЦЕЙСКИЙ УНИЧТОЖЕН: +50 КР НАГРАДА | game.go |
| `POLICE VIPER ON INTERCEPT` | ПОЛИЦЕЙСКИЙ ВАЙПЕР | game.go |
| `PRESS 2ND` | ВТОРАЯ | screens.go |
| `PRESS ENTER TO CONTINUE` | НАЖМИ ENTER | ui.go |
| `PURCHASED: ` | КУПЛЕНО:  | update.go |
| `RAW AXES` | СЫРЫЕ ОСИ | screens.go |
| `REFUEL  [F]` | ЗАПРАВКА [F] | screens.go |
| `REFUELLED FOR ` | ЗАПРАВЛЕНО ЗА  | update.go |
| `RESET DEFAULTS` | СБРОС ПО УМОЛЧАНИЮ | screens.go |
| `RESTORED: ` | ВОССТАНОВЛЕНО:  | screens.go |
| `SAVE FAILED` | ОШИБКА СОХРАНЕНИЯ | screens.go |
| `SCOOPED: ` | ПОДОБРАНО:  | update.go |
| `SECOND KEY: ` | ВТОРАЯ КЛАВИША:  | screens.go |
| `SECOND SLOT TAKES A KEY` | ВО ВТОРОЙ СЛОТ ТОЛЬКО КЛАВИША | screens.go |
| `SELL CARGO FIRST` | СНАЧАЛА ПРОДАЙТЕ ГРУЗ | shipyard.go |
| `SHIFT+M - SOUND ON/OFF` | SHIFT+M - ЗВУК ВКЛ/ВЫКЛ | screens.go |
| `SHRIKE DETECTED` | ШРАЙК ОБНАРУЖЕН | missions.go |
| `SLOT EMPTY` | СЛОТ ПУСТ | screens.go |
| `SOLD ` | ПРОДАНО  | update.go |
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
| `BUY` | screens.go |
| `CARGO:` | screens.go |
| `CASH ` | ui.go |
| `COLLISION` | update.go |
| `COMMODITY` | screens.go |
| `CONTRACT ACCEPTED: ` | campaign.go |
| `CONTRACT EXPIRED  -` | contracts.go |
| `CT` | ui.go |
| `D-UP ` | screens.go |
| `DELIVERED  +` | campaign.go, contracts.go |
| `DIST ` | screens.go |
| `DOCKED` | game.go, screens.go |
| `DOCUMENTS ABOARD - DELIVER TO ` | missions.go |
| `DOCUMENTS DELIVERED  +` | missions.go |
| `EMPTY` | screens.go |
| `FEE: ` | screens.go |
| `FIT A MINING LASER FIRST` | update.go |
| `FITTED` | screens.go |
| `FLYING ` | screens.go |
| `FS` | ui.go |
| `FU` | ui.go |
| `FUEL ` | screens.go |
| `GALAXY ` | screens.go |
| `HELIOS COMBINE IS FINISHED` | campaign.go |
| `HOLD` | screens.go |
| `HOLD ` | screens.go, ui.go |
| `HT` | ui.go |
| `HYPERSPACE MISJUMP - WITCHSPACE` | game.go |
| `JETTISONED ` | mechanics.go |
| `KG` | screens.go |
| `LASER DOWNGRADED FOR THIS HULL` | shipyard.go |
| `LOAD FAILED` | screens.go |
| `LSR` | ui.go |
| `LT ` | screens.go |
| `MAIL DELIVERED  +` | contracts.go |
| `MISSILE RACKS FULL` | update.go |
| `MSL` | ui.go |
| `NEEDS A LARGE CARGO BAY FIRST` | update.go |
| `NEEDS A SHIELD BOOSTER FIRST` | update.go |
| `NO STATION IN THIS SYSTEM` | game.go |
| `NO STATION TO DOCK WITH` | game.go |
| `NO YAW IN CLASSIC FLIGHT` | update.go |
| `NOW FLYING ` | shipyard.go |
| `OK` | screens.go |
| `PASSENGERS DELIVERED  +` | contracts.go |
| `PAX ` | screens.go |
| `PAYS ` | screens.go |
| `POP ` | screens.go |
| `PRESS KEY` | screens.go |
| `PROD ` | screens.go |
| `QTY` | screens.go |
| `RADIUS ` | screens.go |
| `ROCK BREAKS UP` | mining.go |
| `SELL` | screens.go |
| `SHIELD ` | screens.go |
| `SPD ` | ui.go |
| `STAGE ` | screens.go |
| `STAGE COMPLETE` | campaign.go |
| `STAGE COMPLETE  +` | campaign.go |
| `STN ` | ui.go |
| `SYS ` | ui.go |
| `TECH LEVEL ` | screens.go |
| `THIS HULL CANNOT MOUNT IT` | update.go |
| `TO ` | ui.go |
| `TRADE IN ` | screens.go |

## В словаре есть, но в интерфейсе не найдено

Это ключи, которые остались после переименований и переделок экранов: их можно удалить или они используются кодом, который собирает строку из частей.

- ` IN ` →  В 
- ` TO ` →  В 
- `1-6 SELECT   UP/DOWN   F REFUEL   ENTER CONFIRM` → 1-6 ВЫБОР   ВВЕРХ/ВНИЗ   F ЗАПРАВКА   ENTER
- `1-8 SELECT   UP/DOWN   F REFUEL   ENTER CONFIRM` → 1-8 ВЫБОР   ВВЕРХ/ВНИЗ   F ЗАПРАВКА   ENTER
- `1984 RULES: TWO MISSIONS, NO CONTRACTS` → ПРАВИЛА 1984: ДВЕ МИССИИ, БЕЗ КОНТРАКТОВ
- `ACCELERATE` → УСКОРЕНИЕ
- `ACCELERATE / BRAKE` → ГАЗ / ТОРМОЗ
- `ADVANCED` → РАСШИРЕННЫЙ
- `ALIGN WITH THE SLOT` → СОВМЕСТИСЬ СО ЩЕЛЬЮ
- `AMBUSH` → ЗАСАДА
- `AND YOUR RECORD IS ALREADY FILTHY.` → ПО СЛЕДУ, А РЕПУТАЦИЯ УЖЕ ГРЯЗНАЯ.
- `ANSWERED IT.` → 
- `ARRIVED` → ПРИБЫТИЕ
- `ARROWS / MOUSE` → СТРЕЛКИ / МЫШЬ
- `ARROWS SELECT SYSTEM   H HYPERSPACE   ESC BACK` → СТРЕЛКИ ВЫБОР   H ПРЫЖОК   ESC
- `BACK` → НАЗАД
- `BLACK BOX` → ЧЁРНЫЙ ЯЩИК
- `BOMB` → БОМБА
- `BOUND` → НАЗНАЧЕНО
- `BOUNTY CONTRACT PAID` → КОНТРАКТ ОПЛАЧЕН
- `BOUNTY IN` → ЗАКАЗ В
- `BREAK THE CONVOY.` → РАЗБЕЙ КАРАВАН.
- `BURNED UP IN THE ATMOSPHERE` → СГОРЕЛ В АТМОСФЕРЕ
- `BUY / SELL CARGO` → ТОРГОВЛЯ
- `CABIN OVERHEATED` → КАБИНА ПЕРЕГРЕЛАСЬ
- `CARGO. IT WAS A CALL, AND SOMETHING` → А СИГНАЛОМ. И НА НЕГО ОТВЕТИЛИ.
- `CHART` → КАРТА
- `CLASSIC` → КЛАССИКА
- `COLLIDED WITH AN ASTEROID` → СТОЛКНОВЕНИЕ С АСТЕРОИДОМ
- `COMMANDER STATUS` → СТАТУС
- `CONSUMED BY THE STAR` → ПОГЛОЩЁН ЗВЕЗДОЙ
- `CONTRABAND ON BOARD` → КОНТРАБАНДА НА БОРТУ
- `CONTRACT EXPIRED` → КОНТРАКТ ПРОСРОЧЕН
- `CONTRACTS` → КОНТРАКТЫ
- `CONTROLS` → УПРАВЛЕНИЕ
- `COURIER RUN. SWARM DOCUMENTS MUST` → КУРЬЕРСКИЙ РЕЙС. ДОКУМЕНТЫ РОЯ
- `COURIER TO` → КУРЬЕР В
- `CRASHED INTO THE STATION` → ВРЕЗАЛСЯ В СТАНЦИЮ
- `DANGEROUS. FIND IT AND DESTROY IT.` → ВООРУЖЁН И ОПАСЕН. НАЙДИТЕ И УНИЧТОЖЬТЕ.
- `DATA OUT EVERY SIX DAYS.` → ВЫВОЗИТ ДАННЫЕ КАЖДЫЕ ШЕСТЬ ДНЕЙ.
- `DAYS` → ДН
- `DECELERATE` → ТОРМОЖЕНИЕ
- `DEGREES` → ГРАД
- `DELIVERED` → ДОСТАВЛЕНО
- `DESTINATION:` → НАЗНАЧЕНИЕ:
- `DESTROYED` → УНИЧТОЖЕН
- `DESTROYED BY A MISSILE` → УНИЧТОЖЕН РАКЕТОЙ
- `DESTROYED BY PIRATE` → УНИЧТОЖЕН ПИРАТОМ
- `DESTROYED BY THE SWARM` → УНИЧТОЖЕН РОЕМ
- `DESTROYED)` → УНИЧТОЖЕНО)
- `DETONATED` → ВЗОРВАНА
- `DOCKING CLEARANCE GRANTED` → ПОСАДКА РАЗРЕШЕНА
- `DOCKING COMPUTER` → СТЫКОВКА АВТО
- `DOCUMENTS DELIVERED` → ДОКУМЕНТЫ ДОСТАВЛЕНЫ
- `DOCUMENTS FOR` → ДОКУМЕНТЫ ДЛЯ
- `E.C.M.` → ECM
- `END IT.` → ЗАКОНЧИ ЭТО.
- `ENERGY BOMB` → ЭНЕРГОБОМБА
- `ENGLISH` → АНГЛИЙСКИЙ
- `ENTER ACCEPT   ESC DECLINE` → ENTER ПРИНЯТЬ   ESC ОТКАЗ
- `EQUIP SHIP` → ОСНАСТКА
- `EQUIPMENT` → ОБОРУДОВАНИЕ
- `ESC BACK` → ESC НАЗАД
- `EVERYTHING: CONTRACTS, SHIPYARD, EXTRAS` → ВСЁ: КОНТРАКТЫ, ВЕРФЬ, ДОПОЛНЕНИЯ
- `FEE:` → ГОНОРАР:
- `FINAL APPROACH - ` → ЗАХОД НА ПОСАДКУ - 
- `FIRE` → ОГОНЬ
- `FIRE LASERS` → ОГОНЬ
- `FLIGHT CONTROLS` → УПРАВЛЕНИЕ
- `FLYING` → ВЫ ЛЕТАЕТЕ НА
- `FOR WHOEVER CAME WITH THE SUPPLIES.` → КТО ПРИВЕЗЁТ ПРИПАСЫ.
- `FRAMED` → ПОДСТАВА
- `FROM` → С
- `GALACTIC CHART` → КАРТА
- `GALACTIC CHART - IN FLIGHT` → КАРТА - В ПОЛЁТЕ
- `GALACTIC HYPERSPACE` → ГАЛАКТИЧЕСКИЙ ПРЫЖОК
- `GALACTIC JUMP` → ГАЛАКТИЧЕСКИЙ ПРЫЖОК
- `GALACTIC MARKET` → ГАЛАКТИЧЕСКИЙ РЫНОК
- `HELIOS COMBINE IS RUNNING. THEIR FLAGSHIP` → КОМБИНАТ ХЕЛИОС БЕЖИТ. ИХ ФЛАГМАН
- `HYPERSPACE MISJUMP!` → СБОЙ ГИПЕРПРЫЖКА!
- `IS LOADED AND ALREADY CLEARED TO JUMP.` → ЗАГРУЖЕН И УЖЕ ГОТОВ К ПРЫЖКУ.
- `IS MOVING SOMETHING BIG. THREE HAULERS,` → ПЕРЕВОЗИТ ЧТО-ТО КРУПНОЕ. ТРИ ТЯЖАЧА,
- `IT IS A SIMPLE RUN. THAT IS THE POINT.` → ПРОСТОЙ РЕЙС. В ЭТОМ И СМЫСЛ.
- `JETTISON` → СБРОС ГРУЗА
- `JETTISONED` → СБРОШЕНО
- `JOURNAL` → ЖУРНАЛ
- `JUMP COORDINATES. HIS COURIER RUNS THE` → ТОРГУЮЩЕГО КООРДИНАТАМИ. ЕГО КУРЬЕР
- `JUMP IN SYSTEM` → ПРЫЖОК В СИСТЕМЕ
- `KILL THE MOTHERSHIP BEFORE IT WAKES.` → УБЕЙ МАТКОРАБЛЬ, ПОКА ОН НЕ ПРОСНУЛСЯ.
- `LANGUAGE` → ЯЗЫК
- `LAST SEEN:` → ПОСЛЕДНИЙ РАЙОН:
- `LAUNCH` → СТАРТ
- `LAUNCH MISSILE` → ПУСК РАКЕТЫ
- `LAUNCHED` → СТАРТ
- `LIGHT YEARS` → СВ. ЛЕТ
- `LOAD COMMANDER` → ЗАГРУЗИТЬ ПИЛОТА
- `MAIL` → ПОЧТА
- `MAIL DELIVERED` → ПОЧТА ДОСТАВЛЕНА
- `MASSACRE. NAVY VIPERS ARE ON YOUR TAIL` → КАРАВАНА. ФЛОТСКИЕ ВАЙПЕРЫ ИДУТ
- `MEDICAL RUN` → МЕДИЦИНСКИЙ РЕЙС
- `MISSILE` → РАКЕТА
- `MODE` → РЕЖИМ
- `MOTHERSHIP DESTROYED` → МАТКОРАБЛЬ УНИЧТОЖЕН
- `MOTHERSHIP IN` → МАТКОРАБЛЬ В
- `MOTHERSHIP. IT SITS DEEP IN THE SYSTEM` → РОЯ. ОН ГЛУБОКО В СИСТЕМЕ
- `NAVAL COMMAND HAS TRACKED A SWARM` → ШТАБ ЗАСЁК РОЙ
- `NAVAL INTELLIGENCE HAS LOST A PROTOTYPE` → ВОЕННАЯ РАЗВЕДКА ПОТЕРЯЛА ПРОТОТИП
- `NEW GAME` → НОВАЯ ИГРА
- `NO ESCORT WORTH THE NAME.` → ОХРАНА ТОЛЬКО ДЛЯ ВИДА.
- `NO GAMEPAD` → НЕТ ГЕЙМПАДА
- `NOW FLYING` → ТЕПЕРЬ ВЫ ЛЕТАЕТЕ НА
- `OFF` → ВЫКЛ
- `ON` → ВКЛ
- `PASSENGERS` → ПАССАЖИРЫ
- `PASSENGERS DELIVERED` → ПАССАЖИРЫ ДОСТАВЛЕНЫ
- `PAUSE` → ПАУЗА
- `PAUSED` → ПАУЗА
- `PAYS` → ОПЛАТА
- `PIRATES` → ПИРАТОВ
- `PITCH AND ROLL` → ТАНГАЖ И КРЕН
- `PITCH DOWN` → ТАНГАЖ ВНИЗ
- `PITCH UP` → ТАНГАЖ ВВЕРХ
- `PRESS KEY, BUTTON OR STICK   ESC CANCEL` → НАЖМИ КЛАВИШУ, КНОПКУ ИЛИ СТИК   ESC ОТМЕНА
- `PRESS THE SECOND KEY   ESC CANCEL` → НАЖМИ ВТОРУЮ КЛАВИШУ   ESC ОТМЕНА
- `PRESS...` → НАЖМИТЕ...
- `QUIT` → ВЫХОД
- `REACH NAVAL COMMAND IN THE NEXT GALAXY.` → ДОЛЖНЫ ПОПАСТЬ В ШТАБ СЛЕДУЮЩЕЙ ГАЛАКТИКИ.
- `RECORDER. FIND THE WRECK AND BURN IT -` → САМОПИСЕЦ УЦЕЛЕЛ. НАЙДИ ОБЛОМКИ
- `RESET ROW` → СБРОСИТЬ СТРОКУ
- `RESUME` → ПРОДОЛЖИТЬ
- `ROLL LEFT` → КРЕН ВЛЕВО
- `ROLL RIGHT` → КРЕН ВПРАВО
- `RUSSIAN` → РУССКИЙ
- `SAVE COMMANDER` → СОХРАНИТЬ
- `SET 2ND KEY` → ЗАДАТЬ ВТОРУЮ
- `SET KEY` → ЗАДАТЬ КЛАВИШУ
- `SETTINGS` → НАСТРОЙКИ
- `SHIP` → КОРАБЛЬ
- `SHIP LOST` → КОРАБЛЬ ПОТЕРЯН
- `SHIPYARD` → ВЕРФЬ
- `SHRIKE CLASS FIGHTER. ARMED AND` → ИСТРЕБИТЕЛЬ КЛАССА ШРАЙК.
- `SHRIKE DESTROYED` → ШРАЙК УНИЧТОЖЕН
- `SHRIKE IN` → ШРАЙК В
- `SOUND` → ЗВУК
- `SOUND OFF` → ЗВУК ВЫКЛ
- `SOUND ON` → ЗВУК ВКЛ
- `SPACE or LMB` → ПРОБЕЛ / ЛКМ
- `STATION SERVICES` → УСЛУГИ СТАНЦИИ
- `STATUS` → СТАТУС
- `STOP THE COURIER.` → ОСТАНОВИ КУРЬЕРА.
- `STORY: HUNT THE SHRIKE` → СЮЖЕТ: ОХОТА НА ШРАЙКА
- `STORY: SWARM DOCUMENTS` → СЮЖЕТ: ДОКУМЕНТЫ РОЯ
- `SURVIVE. TWO OF THEM.` → ВЫЖИВИ. ИХ ДВОЕ.
- `SYSTEM DATA` → ДАННЫЕ СИСТЕМЫ
- `TAKE FOUR TONNES OF FOOD OUT TO THEM.` → ОТВЕЗИ ИМ ЧЕТЫРЕ ТОННЫ ПРОДОВОЛЬСТВИЯ.
- `THE ARTEFACT IN THAT CONVOY WAS NOT` → АРТЕФАКТ ИЗ КАРАВАНА БЫЛ НЕ ГРУЗОМ,
- `THE BOX NAMED A NAVY LIAISON SELLING` → ЗАПИСЬ НАЗВАЛА ФЛОТСКОГО СВЯЗНОГО,
- `THE CLINIC WAS A DECOY. THEY WERE WAITING` → КЛИНИКА БЫЛА ПРИМАНКОЙ. ТЕБЯ ЖДАЛИ,
- `THE CONVOY` → КАРАВАН
- `THE COURIER` → КУРЬЕР
- `THE FLAGSHIP` → ФЛАГМАН
- `THE FRONTIER CLINIC IS OUT OF SUPPLIES.` → НА ПОГРАНИЧНОЙ СТАНЦИИ КОНЧИЛИСЬ ПРИПАСЫ.
- `THE HIVE` → УЛЕЙ
- `THEIR SHIP WENT DOWN WITH A FLIGHT` → ИХ КОРАБЛЬ УПАЛ, НО БОРТОВОЙ
- `THEY HAVE PUT YOUR NAME ON THE CONVOY` → ТЕБЯ СДЕЛАЛИ ВИНОВАТЫМ В ГИБЕЛИ
- `THEY WILL TRY TO STOP YOU.` → ВАС БУДУТ ПЫТАТЬСЯ ОСТАНОВИТЬ.
- `THREE OF THEM. CLEAR THE SKY AND LEAVE.` → ИХ ТРОЕ. ОЧИСТИ НЕБО И УХОДИ.
- `TRADE IN` → ЗАЧЁТ
- `UP/DN  B BUY  S SELL  SHIFT+ 10  ESC` → ВВЕРХ/ВНИЗ  B КУПИТЬ  S ПРОДАТЬ  SHIFT 10  ESC
- `UP/DOWN   ENTER LOAD   ESC BACK` → ВВЕРХ/ВНИЗ   ENTER ЗАГРУЗИТЬ   ESC
- `UP/DOWN   ENTER REBIND   ESC BACK` → ВВЕРХ/ВНИЗ   ENTER ИЗМЕНИТЬ   ESC
- `UP/DOWN   ENTER SAVE   ESC BACK` → ВВЕРХ/ВНИЗ   ENTER СОХРАНИТЬ   ESC
- `UP/DOWN   ENTER START   ESC BACK` → ВВЕРХ/ВНИЗ   ENTER СТАРТ   ESC
- `UP/DOWN  LEFT/RIGHT  ENTER  ESC` → ВВЕРХ/ВНИЗ  ВЛЕВО/ВПРАВО  ENTER  ESC
- `UP/DOWN SELECT   ENTER ACCEPT   ESC BACK` → ВВЕРХ/ВНИЗ ВЫБОР   ENTER ПРИНЯТЬ   ESC
- `UP/DOWN SELECT   ENTER BUY   ESC BACK` → ВВЕРХ/ВНИЗ ВЫБОР   ENTER КУПИТЬ   ESC
- `UP/DOWN SELECT   ESC BACK` → ВВЕРХ/ВНИЗ ВЫБОР   ESC НАЗАД
- `VOLUME` → ГРОМКОСТЬ
- `WHOEVER READS THAT BOX WILL COME LOOKING.` → И СОЖГИ. ЗАПИСЬ ПРИВЕДЁТ ИХ К ТЕБЕ.
- `WIGGLE STICKS AND PRESS BUTTONS   ESC BACK` → ШЕВЕЛИ СТИКИ И ЖМИ КНОПКИ   ESC НАЗАД
- `WITCHSPACE` → ВЕДЬМИНО ПРОСТРАНСТВО
- `WITH ITS DRONES AROUND IT.` → В ОКРУЖЕНИИ ДРОНОВ.
- `WITH THE LIAISON EXPOSED, HELIOS COMBINE` → СВЯЗНОЙ РАСКРЫТ, И КОМБИНАТ ХЕЛИОС
- `X / Z` → X / Z
- `YAW LEFT` → РЫСКАНЬЕ ВЛЕВО
- `YAW RIGHT` → РЫСКАНЬЕ ВПРАВО

