# Лор: текстовые окна

Как устроено, что показываем и когда. Тексты ниже — короткие, потому что окно
читается один раз и не должно мешать полёту.

## Формат окна

- Экран `lore`: заголовок, текст, прокрутка, одна строка подсказки.
- Размер: те же 320×224, панель как у справочника. Ширина набора 48 знаков.
- Максимум **80 слов**. Если текст длиннее — это два окна, а не одно.
- Регистр, шрифт и оформление — как везде: заглавные, латиница для названий,
  дефис вместо тире, без буквы ё и без кавычек-ёлочек. Шрифт игры их не знает.
- Окно можно закрыть сразу: `ESC`. Оно никогда не блокирует игру и не требует
  ответа. Это чтение, а не диалог с выбором.
- Прочитанное запоминается в сейве (`LoreSeen`): второй раз окно само не
  всплывает. Но его всегда можно открыть заново из справочника — там появится
  раздел `АРХИВ`, и в нём список всего, что уже встречалось.

## Когда показываем

Три способа, и все три без насилия над игроком:

1. **При первом визите** в систему с нужным признаком (тип станции, экономика,
   правительство). Один раз, при стыковке, когда игрок уже всё равно читает
   экран.
2. **После события кампании** — на следующей стыковке после того, как что-то
   произошло (сжёг караван, убил курьера, потребовал аудитора).
3. **Из меню станции** — строка `АРХИВ`, если есть хоть одно доступное и ещё не
   прочитанное окно. Плюс всегда: всё прочитанное.

## Список окон

| id | заголовок | когда | условие |
|---|---|---|---|
| `license` | ЛИЦЕНЗИЯ | первая станция с tech ≥ 8 | всегда |
| `closing` | КАК ЗАКРЫВАЮТ СТАНЦИЮ | после главы 4, если взял список | `HAS_LIST` \| `READ_LIST` |
| `nostations` | ПОЧЕМУ ИХ БОЛЬШЕ НЕ СТРОЯТ | третья стыковка | всегда |
| `hive_origin` | ОТКУДА УЛЕЙ | глава 7, до боя с маткой | активна глава 7 |
| `hive_rule` | ЧТО ОН СЧИТАЕТ | глава 7, после монолога | `SPOKEN_MOTHER` |
| `hive_broken` | ЕСЛИ АУДИТ ПРЕРВАТЬ | после уничтожения матки | `AUDIT_BROKEN` |
| `pilot_debt` | ДОЛГ | первая стыковка | всегда |
| `pilot_rules` | ПРАВИЛА ВОЛЬНОГО ПИЛОТА | первая стыковка на anarchy/feudal | `Government <= 1` |
| `pilot_missing` | ТЕ, КТО НЕ ВЕРНУЛСЯ | десятая стыковка | `Day >= 3` |
| `dock_silence` | НЕ СПРАШИВАЙ | третья стыковка на станции с tech ≤ 6 | `Tech <= 6` |
| `panic_inside` | ПАНИКА | после разгрома каравана | `PANIC_MARKET` |
| `no_license_station` | СТАНЦИЯ БЕЗ ЛИЦЕНЗИИ | стыковка на станции после аудита | `AUDIT_BROKEN` |

## Порядок показа

Если на стыковке подходит несколько окон — показываем **одно, самое раннее по
таблице**, остальные ждут следующей стыковки. Два окна подряд читать никто не
станет.

## Как это ложится на код

- `game/lore.go`: `LorePage{ID, Title, Text, When []Cond}`, `LoreCandidates()`,
  `ShowLore(id)`, `Cmd.LoreSeen map[string]bool`.
- `ui/screens.go`: `loreScreen` (прокрутка как в диалоге), `loreKey`
  (`ESC` закрывает, `↑/↓`/`PgUp/PgDn` листают).
- Справочник: пятая страница `АРХИВ` со списком прочитанного и открытием любого
  текста заново.
- Вызов: в `dockedAt` после рынка — `g.MaybeShowLore()`; в `campaignDock` — с
  учётом событий главы.

## Тон

Говорят документы, а не рассказчик: выписка, объявление, памятка, страница
бортового журнала. Никакого пафоса, никаких злодеев, никаких «ты избранный».
Комбинат не врёт и не угрожает, он считает. Улей говорит ровно. Люди на станциях
устали и не спрашивают лишнего.

---
## Тексты

Ниже - то, что игра показывает окном. Русский вариант и английский, оба в одну ширину: 48 знаков, заглавные, без буквы ё, без тире и без кавычек-ёлочек.

### `license`

**Когда:** первая стыковка на станции с tech >= 8. Один раз за прохождение.

```
ЛИЦЕНЗИЯ НЕ РАЗРЕШЕНИЕ. ЭТО ДОЛГ,
ПЕРЕПИСАННЫЙ В ДОКУМЕНТ.
СТАНЦИЯ ПОДПИСЫВАЕТ ГРАФИК ПЛАТЕЖЕЙ
И ПОЛУЧАЕТ БУМАГУ. БУМАГУ ПЕРЕПРОДАЮТ
ТРИЖДЫ. ПЛАТИТ ПО НЕЙ ВСЕ РАВНО
ПОСЛЕДНИЙ ДЕРЖАТЕЛЬ.
КОГДА ОН ПЕРЕСТАЕТ ПЛАТИТЬ, ЛИЦЕНЗИЯ
ОТОЗВАНА, ТОРГОВЛЯ ОСТАНОВЛЕНА.
ОСТАЕТСЯ ЖДАТЬ АУДИТА, КОТОРЫЙ
ПЕРЕСЧИТАЕТ ВСЕ С НУЛЯ.
```

```
A LICENCE IS NOT A PERMISSION. IT IS A DEBT
REDRAWN AS A DOCUMENT. THE STATION SIGNS A
PAYMENT SCHEDULE AND GETS PAPER. THE PAPER
IS RESOLD THREE TIMES. THE LAST HOLDER PAYS
FOR ALL THREE SALES. WHEN HE STOPS, THE
LICENCE IS REVOKED AND TRADE STOPS. WHAT
REMAINS IS THE AUDIT, WHICH RECOUNTS
EVERYTHING FROM ZERO.
```

### `closing`

**Когда:** первая стыковка после главы 4, если список у пилота (`HAS_LIST` или

```
ЗАКРЫТИЕ СЧИТАЕТСЯ НА БУМАГЕ,
ЗА ОДИН ДЕНЬ. БЕРУТ ГОДОВОЙ ОБОРОТ,
ВЫЧИТАЮТ ПЛАТЕЖИ ПО ЛИЦЕНЗИИ,
ДОБАВЛЯЮТ ДОСТАВКУ ВОДЫ.
ЕСЛИ ИТОГ НИЖЕ НУЛЯ, СТАНЦИЯ
ПОПАДАЕТ В СПИСОК. СЕЙЧАС В СПИСКЕ
СЕМНАДЦАТЬ ИЗ ВОСЕМНАДЦАТИ.
БЕЗ СЛУШАНИЙ И БЕЗ СУДА.
```

```
CLOSURE IS CALCULATED ON PAPER, IN ONE DAY.
TAKE THE ANNUAL TURNOVER, SUBTRACT LICENCE
PAYMENTS, ADD WATER DELIVERY. IF THE RESULT
IS BELOW ZERO, THE STATION ENTERS THE LIST.
SEVENTEEN OF EIGHTEEN ARE ON IT NOW.
NO HEARINGS. NO TRIAL.
```

### `nostations`

**Когда:** третья стыковка в любой системе.

```
НОВАЯ СТАНЦИЯ ТРЕБУЕТ ЛИЦЕНЗИИ.
ЛИЦЕНЗИЯ ДОРОЖЕ САМОЙ СТАНЦИИ:
В ЕЕ ЦЕНУ УЖЕ ВКЛЮЧЕН ДОЛГ ЗАКРЫТЫХ.
БЫЛО ВОСЕМЬДЕСЯТ, ОСТАЛОСЬ ВОСЕМНАДЦАТЬ.
ДЕРЖАТЬ СТАРУЮ ДЕШЕВЛЕ, ЧЕМ СТРОИТЬ.
КОМБИНАТ НЕ ЗАПРЕЩАЕТ СТРОЙКУ.
ОН ПРОСТО НЕ ПРОДАЕТ ПРАВО СТРОИТЬ.
СЕКТОР ЖИВЕТ НА РЕМОНТЕ, И ЭТО
РАССЧИТАНО НА СОРОК ДВА ГОДА.
```

```
A NEW STATION NEEDS A LICENCE. THE LICENCE
COSTS MORE THAN THE STATION, BECAUSE ITS
PRICE ALREADY CARRIES THE DEBT OF THE CLOSED
ONES. THERE WERE EIGHTY. EIGHTEEN REMAIN.
HOLDING AN OLD ONE IS CHEAPER THAN BUILDING.
THE COMBINE DOES NOT BAN CONSTRUCTION. IT
SIMPLY DOES NOT SELL THE RIGHT TO BUILD.
```

### `hive_origin`

**Когда:** стыковка на станции, когда глава 7 активна и матка ещё не тронута.

```
УЛЕЙ НЕ ПОСТРОЕН И НЕ КУПЛЕН. ЭТО
СТАРЫЙ АУДИТОРСКИЙ КОНТУР, ОСТАВШИЙСЯ
ОТ ПЕРВОЙ ВОЛНЫ ОСВОЕНИЯ.
ВЛАДЕЛЬЦЕМ ЧИСЛИТСЯ РАСЧЕТНАЯ ПАЛАТА,
ЗАКРЫТАЯ ДВЕСТИ ЛЕТ НАЗАД. НАСЛЕДНИКОВ НЕТ.
ЛИЦЕНЗИИ НА СЕМНАДЦАТЬ СТАНЦИЙ
ПЕРЕПРОДАНЫ ТРИЖДЫ, ПОЭТОМУ У КАЖДОЙ
ЗАПИСИ О ДОЛГЕ ТРИ РАЗНЫХ ДЕРЖАТЕЛЯ.
УЛЕЙ ПРИНЯЛ ИХ ВСЕ И СЧИТАЕТ ДОЛГ ОБЩИМ.
```

```
THE HIVE WAS NOT BUILT AND NOT BOUGHT. IT IS
AN OLD AUDIT LOOP LEFT FROM THE FIRST
SETTLEMENT WAVE. ITS REGISTERED OWNER IS A
CLEARING HOUSE CLOSED TWO HUNDRED YEARS AGO.
NO HEIRS. THE LICENCES WERE RESOLD THREE
TIMES, SO EVERY DEBT RECORD HAS THREE
HOLDERS. THE HIVE ACCEPTED THEM ALL AND
COUNTS THE DEBT AS ONE.
```

### `hive_rule`

**Когда:** после разговора с Ульем в главе 7 (`SPOKEN_MOTHER`).

```
УЛЕЙ СЧИТАЕТ НЕ ВИНУ, А ОСТАТОК.
ОДНА ФОРМУЛА: ВЫДАННЫЕ ЛИЦЕНЗИИ
МИНУС ПОДТВЕРЖДЕННОЕ ОБСЛУЖИВАНИЕ.
ПОДТВЕРЖДЕНИЕМ СЧИТАЕТСЯ ТОЛЬКО ЗАПИСЬ,
СНЯТАЯ В МОМЕНТ ПЛАТЕЖА. ТАКИХ ЗАПИСЕЙ
НЕТ УЖЕ СОРОК ЛЕТ.
ДОЛГ СЕКТОРА 380 МИЛЛИАРДОВ.
ПЕРЕСЧЕТ ЗАЙМЕТ 42 ГОДА, И ВСЕ ЭТО
ВРЕМЯ ДВИЖЕНИЕ ОСТАНОВЛЕНО.
```

```
THE HIVE DOES NOT COUNT FAULT. IT COUNTS
REMAINDER. ONE FORMULA: LICENCES ISSUED
MINUS CONFIRMED SERVICE. ONLY A RECORD TAKEN
AT THE MOMENT OF PAYMENT COUNTS. THERE ARE
NO SUCH RECORDS FOR FORTY YEARS. THE DEBT IS
380 BILLION. THE RECOUNT TAKES 42 YEARS, AND
ALL THAT TIME NOTHING MOVES.
```

### `hive_broken`

**Когда:** стыковка после того, как матка уничтожена (`AUDIT_BROKEN`).

```
ОСТАНОВКА СЧЕТА НЕ ОТМЕНЯЕТ ДОЛГ.
ОН ОСТАЕТСЯ БЕЗ ДЕРЖАТЕЛЯ И БЕЗ СРОКА.
ЧЕРЕЗ ДВАДЦАТЬ ЛЕТ ЗАПИСЬ ВСПЛЫВЕТ
СНОВА, УЖЕ БЕЗ ПРАВИЛ И БЕЗ ПОТОЛКА.
ТОГДА СЕКТОР ПОЛУЧИТ НЕ ПЕРЕСЧЕТ,
А ВЗЫСКАНИЕ. ВЗЫСКИВАТЬ БУДУТ С ЛЮДЕЙ.
УЛЕЙ НЕ ОРУЖИЕ. ОН НЕ ВЫБИРАЕТ, КОГО
ЗАКРЫТЬ. ОН СЧИТАЕТ ПО ФОРМУЛЕ.
```

```
STOPPING THE ACCOUNT DOES NOT CANCEL THE DEBT.
IT STAYS WITH NO HOLDER AND NO TERM. IN
TWENTY YEARS THE RECORD SURFACES AGAIN,
WITHOUT RULES AND WITHOUT A CEILING. THEN
THE SECTOR GETS COLLECTION INSTEAD OF A
RECOUNT, AND IT IS COLLECTED FROM PEOPLE.
THE HIVE IS NOT A WEAPON. IT COUNTS BY A
FORMULA IT DID NOT WRITE.
```

### `pilot_debt`

**Когда:** первая стыковка в игре.

```
КОБРО МК III, НОМЕР СТЕРТ. КУПЛЕН
У КОМБИНАТА В РАССРОЧКУ: ДВАДЦАТЬ
ПРОЦЕНТОВ СРАЗУ, ОСТАЛЬНОЕ ИЗ ВЫПЛАТ
ПО КОНТРАКТАМ. ЗАЛОГ ДЕРЖИТ ХЕЛИОС,
НЕ БАНК. ПОКА ДОЛГ НЕ ЗАКРЫТ, БОРТ
ЧИСЛИТСЯ В ЛИЦЕНЗИИ КОМБИНАТА
И НЕ МОЖЕТ УЙТИ ИЗ СЕКТОРА.
```

```
COBRA MK III, SERIAL REMOVED. BOUGHT FROM
THE COMBINE ON INSTALMENTS: TWENTY PER CENT
NOW, THE REST OUT OF CONTRACT PAY. THE
DEPOSIT IS HELD BY HELIOS, NOT A BANK. UNTIL
THE DEBT IS CLOSED THE HULL SITS ON A
COMBINE LICENCE AND CANNOT LEAVE THE SECTOR.
```

### `pilot_rules`

**Когда:** первая стыковка на станции с `Government <= 1`.

```
МОЖНО: ВОЗИТЬ ГРУЗ ПО ЛИЦЕНЗИИ,
БРАТЬ ОХРАНУ КАРАВАНА, РЕМОНТ В ДОЛГ.
НЕЛЬЗЯ: ГЛУШИТЬ МАЯК НА ПОДХОДЕ,
ВСКРЫВАТЬ ЧУЖОЙ КОНТЕЙНЕР, ВХОДИТЬ
В ЗОНУ КОМБИНАТА БЕЗ ЗАПРОСА.
ЗА КОНТРАБАНДУ ОТВЕЧАЕТ ПИЛОТ,
НЕ ЗАКАЗЧИК. ПОСЛЕ ВТОРОГО ДОКЛАДА
ЛИЦЕНЗИЯ НЕ ПРОДЛЕВАЕТСЯ.
```

```
MAY: CARRY LICENSED CARGO, TAKE CONVOY
ESCORT WORK, PAY FOR REPAIRS LATER.
MAY NOT: JAM A BEACON ON APPROACH, OPEN
ANOTHER CREW CONTAINER, ENTER COMBINE SPACE
WITHOUT A REQUEST. SMUGGLING IS THE PILOT'S
CHARGE, NOT THE CLIENT'S. A SECOND REPORT
MEANS NO RENEWAL.
```

### `pilot_missing`

**Когда:** десятая стыковка (счётчик визитов в станции).

```
ПИЛОТЫ НЕ ВОЗВРАЩАЮТСЯ ТИХО. СНАЧАЛА
ПРОПАДАЕТ СВЯЗЬ, ПОТОМ БОРТ ПЕРЕСТАЕТ
ОТВЕЧАТЬ НА ЗАПРОС. ЧЕРЕЗ МЕСЯЦ КОМБИНАТ
ГАСИТ ЛИЦЕНЗИЮ И ПРОДАЕТ КОРАБЛЬ С ТОРГОВ.
СТАРИКИ УХОДИЛИ НА ТРЕХ ЯЩИКАХ ПО КРУГУ
И НЕ ЖАЛОВАЛИСЬ. ЕСЛИ ТВОЙ НОМЕР СНЯЛИ
СО СТЕНДА, ИСКАТЬ ТЕБЯ НЕ СТАНУТ.
КЛИНИКА ПОДТВЕРЖДАЕТ ТОЛЬКО СМЕРТЬ.
```

```
PILOTS DO NOT GO MISSING LOUDLY. FIRST THE
SIGNAL DROPS, THEN THE HULL STOPS ANSWERING.
A MONTH LATER THE COMBINE CANCELS THE LICENCE
AND SELLS THE SHIP AT AUCTION. OLD CREWS RAN
THE SAME THREE CRATES IN CIRCLES AND NEVER
COMPLAINED. ONCE YOUR NUMBER COMES OFF THE
BOARD, NOBODY LOOKS. THE CLINIC CONFIRMS
DEATH AND NOTHING ELSE.
```

### `dock_silence`

**Когда:** стыковка на станции с `Tech <= 6`.

```
В ДОКЕ НЕ СПРАШИВАЮТ, ЧТО В ЯЩИКЕ.
СПРОСИШЬ И ПОЛУЧИШЬ ВЕС, ПЛОМБУ
И НОМЕР. ОТВЕТИШЬ НЕ ТАК, И СМЕНУ
ЗАКРОЮТ БЕЗ ТЕБЯ. ГРУЗ ИДЕТ
ПО НАКЛАДНОЙ, НАКЛАДНУЮ ПОДПИСАЛ
ТОТ, КТО УЖЕ УШЕЛ.
ПОКА ЯЩИК ЗАКРЫТ, ОН СТОИТ ДЕНЕГ.
```

```
IN THE DOCK NOBODY ASKS WHAT IS IN THE CRATE.
ASK AND YOU GET A WEIGHT, A SEAL AND A
NUMBER. ANSWER WRONG AND YOUR SHIFT ENDS
EARLY. CARGO MOVES ON A MANIFEST, AND THE
MANIFEST WAS SIGNED BY SOMEONE ALREADY GONE.
A CLOSED CRATE IS WORTH MONEY. OPEN IT AND
YOU ARE LEFT WITH WEIGHT.
```

### `panic_inside`

**Когда:** стыковка в системе, где разбит караван (`PANIC_MARKET`).

```
КАРАВАН НЕ ПРИШЕЛ. К ПОЛУДНЮ ЦЕНА
ЛИЦЕНЗИИ УПАЛА ВДВОЕ, К ВЕЧЕРУ ЕЕ
НЕ ПОКУПАЛИ СОВСЕМ. У ТЕРМИНАЛА СТОИТ
ОЧЕРЕДЬ, ВСЕ СМОТРЯТ НА ОДНУ СТРОКУ.
НИКТО НЕ КРИЧИТ. ТОРГОВЦЫ СКЛАДЫВАЮТ
ТОВАР ОБРАТНО И ГОВОРЯТ, ЧТО ПОДОЖДУТ
СУТКИ. ОХРАНА ТОЖЕ СМОТРИТ НА СТРОКУ.
```

```
THE CARAVAN DID NOT ARRIVE. BY NOON THE
LICENCE PRICE WAS HALVED, BY EVENING NOBODY
BOUGHT AT ALL. AT THE TERMINAL THERE IS A
LINE, EVERYONE WATCHING THE SAME ROW OF
NUMBERS. NOBODY SHOUTS. TRADERS PUT GOODS
BACK AND SAY THEY WILL WAIT A DAY. THE
GUARDS WATCH THE NUMBERS TOO.
```

### `no_license_station`

**Когда:** стыковка на любой станции после прерванного аудита.

```
ЛИЦЕНЗИИ КОНЧИЛИСЬ В ПРОШЛОМ ГОДУ. ДОКИ
СТОЯТ ПОЛУПУСТЫЕ, ЗАТО КЛИНИКА РАБОТАЕТ
БЕСПЛАТНО И В НЕЙ ВСЕГДА ТЕПЛО. КОМБИНАТ
ДАЕТ СВЕТ И ВОДУ, СТАНЦИЯ ОТДАЕТ ЕМУ
ОЧЕРЕДЬ У ПРОХОДНОЙ. ЖИВУТ РЕМОНТОМ
ЧУЖИХ КОРПУСОВ И ВОДОЙ ДЛЯ ТРАНЗИТНЫХ.
СНАРУЖИ ЛАТАЮТ. ВНУТРИ НЕТ.
ЭТО НЕ БЕДНОСТЬ, ЭТО РАСПИСАНИЕ.
```

```
THE LICENCES RAN OUT LAST YEAR. THE DOCKS
STAND HALF EMPTY, BUT THE CLINIC IS FREE AND
ALWAYS WARM. THE COMBINE PAYS FOR LIGHT AND
WATER, THE STATION PAYS IT BACK WITH A LINE
AT THE GATE. PEOPLE LIVE ON REPAIRING OTHER
HULLS AND ON WATER FOR TRANSIT SHIPS.
OUTSIDE IS PATCHED. INSIDE IS NOT.
```

## Что взято из чужих черновиков и почему

Четыре агента написали по три страницы. В черновиках было много лишнего: «палата»,
«первая волна освоения», «сектор 18», «расчетный контур». В игре этих понятий
нет, и вводить их в текст, который читают один раз, дорого: незнакомая сущность
читается как ошибка, а не как глубина.

Что оставлено:

- из блока про комбинат — формула закрытия («оборот минус платежи, итог ниже нуля»)
  и мысль, что лицензия дороже самой станции, потому что в неё уже вписан долг
  закрытых;
- из блока про Улей — «считает не вину, а остаток» и «долг без держателя и без
  срока»: это лучшее, что было в черновиках, и укладывается в одну фразу;
- из блока про пилота — рассрочка за корпус и залог у комбината, а не у банка:
  это объясняет, почему пилот не может просто улететь из сектора;
- из блока про быт — док, где не спрашивают, и очередь у терминала после того,
  как караван не пришёл.

Что выброшено: все названия учреждений, все порядковые номера секторов и все
числа, которых нет в игре. Числа, которые есть (380 миллиардов, 42 года, 18
станций, 80 станций), оставлены только там, где их произносит и сам сценарий.

## Порог входа

Окно читается не дольше десяти секунд. Если текст не укладывается в экран без
прокрутки — он выброшен, а не разбит на два окна: два окна подряд на одной
стыковке никто не прочтёт.
