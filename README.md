# ContestOzonGolangSchool2020
Решение задач отбора на школу Go разработки Ozon

ТЗ лежит в корне в pdf

В итоге не решил последнюю задачу E на go.
Проблема в том, что есть ситуация, когда пишет ошибку WA и в логах пусто.
Пробовал разные подходы.

Получилось 11 вариантов, хотя первый слишком медленный, но он был нужен для тестов :)
1) asyncF - просто в goroutine запускаем перебор, выдает IL стабильно, что ожидаемо.
2) asyncF2 - расчет f по каждому входному каналу в отдельной анонимной goroutine, передача через переменные, ожидание через WaitGroup для сложения. Хорошая скорость(IL не выдает).
3) asyncChannels - по сути тоже самое, что и async2, но реализованно отдельными функциями и результаты вычислений передаются через промежуточные каналы, которые читает и складывает еще одна goroutine. Хорошая скорость(IL не выдает).
4) asyncFAndSum - одна goroutine на вычисление, другая на сложение, результаты из первой во вторую передаются через каналы. Похожа на asyncF, просто как промежуточный этап развития.
5) sync1WG - в goroutine вычисляем сумму и ждем через WaitGroup. Повторение asyncFAndSum на основе WaitGroup.
6) sync2WG - на каждую итерацию запускаем goroutine, которая вычисляет сумму и ожидает завершения предыдущей goroutine перед тем как писать в Out, чтобы соблюсти порядок. Снхронизация через WaitGroup. Быстрая. В результате цепочка ожидающих WaitGroup длинной в кол-во выполнений, при устремлении n в бесконечность мне кажется что закончится какойнибудь ресурс WaitGroup или стек, не знаю как все это под капотом..
7) sync4WG - тоже самое, что и sync2WG, только чтение из каналов и вычисление так же разнесено на разные goroutine и синхронизируется через WaitGroup, но судя из того, что запись в каналы выполнена перед запуском первого теста - толку нет, только лишние накладные расходы.
8) mapAndChan - для каждого входного канала goroutine читает значение и на каждую итерацию запускает еще одну goroutine для вычисления результата и сохранения в syncMap, где ключ - это номер итерации, а значение - это результат выполнения f, так же после записи в syncMap через канал посылается сигнал, что в syncMap есть, что почитать. Еще одна goroutine читает syncMap по сигналу и пишет в промежуточный канал. Последняя goroutine читает промежуточные каналы и складывает.
9) map2out - тоже что и mapAndChan, только goroutine, вычисляющая  f для конкретной итерации чтения канала, читает сразу syncMap каждого канала и складывает сразу в out. Думаю, что чтение каналов в mapAndChan, эффективнее, чем опрос syncMap, особонно если вычисление f абсолютно разное время занимает, а порядок в Out сохранять нужно..
10) mapAndCache - аналогично mapAndChan, только расчеты f кэшируются и если есть одинаковые значения во входных каналах, то заметно прибавляет в скорости..
11) syncCounter - похоже на mapAndChan, только вместо отдельной goroutine syncMap читает та же goroutine, что и вычисляет f для конкретной итерации, а синхронизацию, для порядка в out обеспечивает счетчик, который через мьютекс считается. Просто попробовал ради интереса, а так тоже мьютекс думаю дороже, чем лишняя goroutine и передача через каналы.

Самые быстрые получились это map2out, sync2WG и asyncF2(asyncChannels).
Суть у перых двух одна, это вычислять f для каждого значения асинхронно, но при этом запоминая порядок для выдачи суммы в out. Первая делает это через мьютекс, вторая через WaitGroup.
Последняя просто параллельно вычисляет f для каждого канала, но чтение из каналов синхронное.

Мне больше всего понравилась mapAndChan.

Результаты тестов на i7-9750H
30 итераций
[asyncF] max: 30 complexity: 40000000. elapsed 573.1365ms
[asyncF2] max: 30 complexity: 40000000. elapsed 295.7153ms
[asyncChannels] max: 30 complexity: 40000000. elapsed 294.8214ms
[asyncFAndSum] max: 30 complexity: 40000000. elapsed 612.3635ms
[sync1WG] max: 30 complexity: 40000000. elapsed 574.0347ms
[sync2WG] max: 30 complexity: 40000000. elapsed 116.6828ms
[sync4WG] max: 30 complexity: 40000000. elapsed 107.685ms
[mapAndChan] max: 30 complexity: 40000000. elapsed 104.7199ms
[map2out] max: 30 complexity: 40000000. elapsed 105.7176ms
[mapAndCache] max: 30 complexity: 40000000. elapsed 24.933ms
[syncCounter] max: 30 complexity: 40000000. elapsed 109.7061ms

300
[asyncF] max: 300 complexity: 40000000. elapsed 5.7353861s
[asyncF2] max: 300 complexity: 40000000. elapsed 2.9136942s
[asyncChannels] max: 300 complexity: 40000000. elapsed 2.9103989s
[asyncFAndSum] max: 300 complexity: 40000000. elapsed 5.678115s
[sync1WG] max: 300 complexity: 40000000. elapsed 5.6859444s
[sync2WG] max: 300 complexity: 40000000. elapsed 1.0262657s
[sync4WG] max: 300 complexity: 40000000. elapsed 1.0202696s
[mapAndChan] max: 300 complexity: 40000000. elapsed 1.0177831s
[map2out] max: 300 complexity: 40000000. elapsed 1.0212678s
[mapAndCache] max: 300 complexity: 40000000. elapsed 25.9307ms
[syncCounter] max: 300 complexity: 40000000. elapsed 1.0227899s

3000
[asyncF] max: 3000 complexity: 40000000. elapsed 58.9259064s
[asyncF2] max: 3000 complexity: 40000000. elapsed 29.9316544s
[asyncChannels] max: 3000 complexity: 40000000. elapsed 29.8778229s
[asyncFAndSum] max: 3000 complexity: 40000000. elapsed 58.9555956s
[sync1WG] max: 3000 complexity: 40000000. elapsed 59.1204983s
[sync2WG] max: 3000 complexity: 40000000. elapsed 10.2157454s
[sync4WG] max: 3000 complexity: 40000000. elapsed 10.225214s
[mapAndChan] max: 3000 complexity: 40000000. elapsed 10.2162285s
[map2out] max: 3000 complexity: 40000000. elapsed 10.5373989s
[mapAndCache] max: 3000 complexity: 40000000. elapsed 26.9275ms
[syncCounter] max: 3000 complexity: 40000000. elapsed 10.2970105s

По остальным задачам нет проблем, если не тупить, ничего особенного они не представляют, хотя я бы попробовал решить F на nodejs еще..
У меня они решены на C# 
