package main
// 1.
import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"
)

// Что выведет код и почему?

func setLinkHome(link *string) {
	*link = "http://home"
}

link := "http://other"
setLinkHome(&link)
fmt.Println(link)


// 2.
// Будет ли напечатан “ok” ?

func main() {
	defer func() {
		recover()
	}()
	panic("test panic")
	fmt.Println("ok")
}


// 3.
// Функция должна напечатать:
// one
// two
// three
// (в любом порядке и в конце обязательно)
// done!
// Но это не так, исправь код

func printText(data []string) {

	wg := sync.WaitGroup{}
	for _, v := range data {
		wg.Add(1)
		v = v
		go func(v ) {
			defer wg.Done()
			fmt.Println(v)
		}()
	}

	wg.Wait()
	fmt.Println("done!")
}

data := []string{"one", "two", "three"}
printText(data)


// 4.
// Мы пытаемся подсчитать количество выполненных параллельно операций,
// что может пойти не так?

var callCounter uint

func main() {
	wg := sync.WaitGroup{}
	wg.Add(10000)
	for i := 0; i < 10000; i++ {
		go func() {
			// Ходим в базу, делаем долгую работу
			time.Sleep(time.Second)
			// Увеличиваем счетчик
			callCounter++
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("Call counter value = ", callCounter)
}


// 5.
// Есть функция processDataInternal, которая может выполняться неопределенно долго.
// Чтобы контролировать процесс, мы добавили таймаут выполнения ф-ии через context.
// Какие недостатки кода ниже?

func (s *Service) ProcessData(timeoutCtx context.Context, r io.Reader) error {
	errCh := make(chan error)

	go func() {
		errCh <- s.processDataInternal(r)
	}()

	select {
	case err := <-errCh:
		return err
	case <-timeoutCtx.Done():
		return timeoutCtx.Err()
	}
}


// 6.
// Опиши, что делает функция isCallAllowed?

var callCount = make(map[uint]uint)
var locker = &sync.Mutex{}

func isCallAllowed(allowedCount uint) bool {
	if allowedCount == 0 {
		return true
	}

	locker.Lock()
	defer locker.Unlock()

	curTimeIndex := uint(time.Now().Unix() / 30)

	prevIndexVal, _ := callCount[curTimeIndex-1]
	if prevIndexVal >= allowedCount {
		return false
	}

	if curIndexVal, ok := callCount[curTimeIndex]; !ok {
		callCount[curTimeIndex] = 1
		return true
	}

	if (curIndexVal + prevIndexVal) >= allowedCount {
		return false
	}

	callCount[curTimeIndex]++
	return true
}

func main() {
	fmt.Printf("%v\n", isCallAllowed(3)) // true
	fmt.Printf("%v\n", isCallAllowed(3)) // true
	fmt.Printf("%v\n", isCallAllowed(3)) // true
	// time.Sleep(time.Second*30)
	fmt.Printf("%v\n", isCallAllowed(3)) // false
	fmt.Printf("%v\n", isCallAllowed(3)) // false
}


// 7 Mysql. Есть две таблицы users, user_cars. У одного пользователя может быть неограниченное количество машин.
// Необходимо написать запрос, который вернет 10 пользователей, у которых нет авто с car_id = 1
+------------------+
|users             |
+------------------+
|id uint           | - PK
|name string       |
+------------------+

+------------------+
|user_cars         | - uniq index (user_id, car_id)
+------------------+
|user_id uint      |
|car_id uint       |
+------------------+

SELECT ...


// 8. Mysql. DevOps говорит, что в slowlog есть запрос, который выполняется дольше 10 секунд. Он отдал вам запрос и вы вызвали explain
// 8.1. О чем вам говорит вывод explain?

+--+------------------+-----+----------+------+-------------------------------------------+-------------------------------------------+-------+-----------------+------+--------+------------------------+
|id|select_type       |table|partitions|type  |possible_keys                              |key                                        |key_len|ref              |rows  |filtered|Extra                   |
+--+------------------+-----+----------+------+-------------------------------------------+-------------------------------------------+-------+-----------------+------+--------+------------------------+
|1 |PRIMARY           |mc   |NULL      |ref   |idx_manager_id_client_id_uindex            |idx_manager_id_client_id_uindex            |1023   |const            |1     |100     |Using where; Using index|
|1 |PRIMARY           |m    |NULL      |eq_ref|idx_user_id                                |idx_user_id                                |1022   |bind.mc.client_id|1     |100     |Using where             |
|2 |DEPENDENT SUBQUERY|cdp  |NULL      |index |idx_client_id                              |idx_client_id                              |1022   |NULL             |189480|20.61   |Using where             |
+--+------------------+-----+----------+------+-------------------------------------------+-------------------------------------------+-------+-----------------+------+--------+------------------------+

// 8.2 Это сам запрос, что можно сделать, чтобы он работал максимально быстро?
// Задача запроса выбрать клиентов определенного менеджера, у которых указано два этапа сделки
SELECT m.*
FROM members m
LEFT JOIN manager_clients mc on m.user_id = mc.client_id
WHERE mc.manager_id = '152734'
AND m.user_id IN (SELECT client_id
FROM client_deal_phases cdp
WHERE cdp.phase_id IN (45, 47)
GROUP BY client_id
HAVING count(client_id) = 2
);
