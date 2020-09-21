package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/PrideSt/otus-golang/hw03_frequency_analysis/internal/pipe"
)

// runSygnalListener subscribes on chSygnal channel and close chTerminate channel on any income sygnal.
func runSygnalListener(wg *sync.WaitGroup, chTerminate chan struct{}) chan os.Signal {
	chSygnal := make(chan os.Signal, 1)
	signal.Notify(chSygnal, syscall.SIGTERM, syscall.SIGINT)

	wg.Add(1)

	go func() {
		defer wg.Done()

		// close terminate channel can only this listener, but it can't write to them
		// main thread can write, but not close
		defer close(chTerminate)

		_, ok := <-chSygnal
		if ok {
			log.Println("terminate sygnal received")
		} else {
			log.Println("terminate sygnal-listener")
		}
	}()

	return chSygnal
}

// Command get some text and find top 10 the most frequent words in it.
// Try to handle signals (it works).
// It wrong way to add signal handling into hw03_frequency_analysis, we add just terminate chan,
// we don't know in TopN call which reasons went to termination.
func main() {
	log.Println("my pid is:", os.Getpid())

	wg := sync.WaitGroup{}
	defer wg.Wait()

	chTerminate := make(chan struct{})
	chSygnal := runSygnalListener(&wg, chTerminate)
	// defer close(chSygnal)
	defer signal.Stop(chSygnal)

	// we can get text from stdin of from file @todo
	text := `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`

	top := pipe.Top10(text, chTerminate)

	fmt.Println(top)
}
