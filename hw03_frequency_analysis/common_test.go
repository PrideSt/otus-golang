package hw03_frequency_analysis //nolint:golint

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/PrideSt/otus-golang/hw03_frequency_analysis/internal/pool"
	"github.com/PrideSt/otus-golang/hw03_frequency_analysis/internal/top"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	goleak.VerifyTestMain(m)

	os.Exit(m.Run())
}

// Change to true if needed
var taskWithAsteriskIsCompleted = true

var text = `Как видите, он  спускается  по  лестнице  вслед  за  своим
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

func TestTop10Single(t *testing.T) {
	testTop10("single", t, Top10)
}

func TestTop10Multy(t *testing.T) {
	testTop10("multy", t, pool.Top10)
}

func testTop10(preffix string, t *testing.T, f func(s string, chTerm <-chan struct{}) []string) {
	t.Run(fmt.Sprintf("[%s] no words in empty string", preffix), func(t *testing.T) {
		require.Len(t, f("", nil), 0)
	})

	t.Run(fmt.Sprintf("[%s] positive test", preffix), func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			expected := []string{"он", "а", "и", "что", "ты", "не", "если", "то", "его", "кристофер", "робин", "в"}
			actual := f(text, nil)
			require.Subset(t, expected, actual)
		} else {
			expected := []string{"он", "и", "а", "что", "ты", "не", "если", "-", "то", "Кристофер"}
			require.ElementsMatch(t, expected, f(text, nil))
		}
	})
}

func TestTopNSingle(t *testing.T) {
	testTopN("single", t, TopN)
}

func TestTopNMulty(t *testing.T) {
	testTopN("multy", t, pool.TopN)
}

func testTopN(preffix string, t *testing.T, f func(s string, topLen int, chTerm <-chan struct{}) []top.FreqEntry) {
	for _, tt := range [...]struct {
		name     string
		input    string
		topLen   int
		expected []top.FreqEntry
	}{
		{
			name:     `empty`,
			input:    ``,
			topLen:   10,
			expected: []top.FreqEntry{},
		},
		{
			name:     `one word`,
			input:    `one one one one one one`,
			topLen:   10,
			expected: []top.FreqEntry{{`one`, 6}},
		},
		{
			name:   `top overflow`,
			input:  `one two two three three three four four four four`,
			topLen: 3,
			expected: []top.FreqEntry{
				{`four`, 4},
				{`three`, 3},
				{`two`, 2},
			},
		},
		{
			name:   `top case sensetive`,
			input:  `one two tWo Three tHree thRee Four fOur foUr fouR`,
			topLen: 3,
			expected: []top.FreqEntry{
				{`four`, 4},
				{`three`, 3},
				{`two`, 2},
			},
		},
	} {
		t.Run(fmt.Sprintf("[%s] %s", preffix, tt.name), func(t *testing.T) {
			require.Equal(t, tt.expected, f(tt.input, tt.topLen, nil))
		})
	}
}
