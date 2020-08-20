package hw03_frequency_analysis //nolint:golint

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/PrideSt/otus-golang/hw03_frequency_analysis/internal/topbuffer"
)

func TestMain(m *testing.M) {
	logger = log.New(ioutil.Discard, "", log.LstdFlags)
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

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		require.Len(t, Top10("", nil), 0)
	})

	t.Run("positive test", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			expected := []string{"он", "а", "и", "что", "ты", "не", "если", "то", "его", "кристофер", "робин", "в"}
			actual := Top10(text, nil)
			require.Subset(t, expected, actual)
		} else {
			expected := []string{"он", "и", "а", "что", "ты", "не", "если", "-", "то", "Кристофер"}
			require.ElementsMatch(t, expected, Top10(text, nil))
		}
	})
}

func TestNormalizeWord(t *testing.T) {
	normalizer := getWordNormalizer()
	for _, tt := range []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     `empty`,
			input:    ``,
			expected: []string(nil),
		},
		{
			name:     `one symbol`,
			input:    `a`,
			expected: []string{`a`},
		},
		{
			name:     `simple word`,
			input:    `word`,
			expected: []string{`word`},
		},
		{
			name:     `simple слово`,
			input:    `слово`,
			expected: []string{`слово`},
		},
		{
			name:     `dash`,
			input:    `-`,
			expected: []string(nil),
		},
		{
			name:     `dash first`,
			input:    `-one`,
			expected: []string{`one`},
		},
		{
			name:     `dash last`,
			input:    `one-`,
			expected: []string{`one`},
		},
		{
			name:     `dash in the middle`,
			input:    `какой-то`,
			expected: []string{`какой-то`},
		},
		{
			name:     `with punctuation last`,
			input:    `hello!`,
			expected: []string{`hello`},
		},
		{
			name:     `with punctuation in the middle`,
			input:    `hello,Masha`,
			expected: []string{`hello`, `masha`},
		},
		{
			name:     `case insensetive`,
			input:    `hElLo`,
			expected: []string{`hello`},
		},
		{
			name:     `with numbers`,
			input:    `i18n`,
			expected: []string{`i18n`},
		},
		{
			name:     `with special characters`,
			input:    "with\ttab",
			expected: []string{`with`, `tab`},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, normalizer(tt.input))
		})
	}
}

func TestTopN(t *testing.T) {
	for _, tt := range []struct {
		name     string
		input    string
		topLen   int
		expected []topbuffer.FreqEntry
	}{
		{
			name:     `empty`,
			input:    ``,
			topLen:   10,
			expected: []topbuffer.FreqEntry{},
		},
		{
			name:     `one word`,
			input:    `one one one one one one`,
			topLen:   10,
			expected: []topbuffer.FreqEntry{{`one`, 6}},
		},
		{
			name:   `top overflow`,
			input:  `one two two three three three four four four four`,
			topLen: 3,
			expected: []topbuffer.FreqEntry{
				{`four`, 4},
				{`three`, 3},
				{`two`, 2},
			},
		},
		{
			name:   `top case sensetive`,
			input:  `one two tWo Three tHree thRee Four fOur foUr fouR`,
			topLen: 3,
			expected: []topbuffer.FreqEntry{
				{`four`, 4},
				{`three`, 3},
				{`two`, 2},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, TopN(tt.input, tt.topLen, nil))
		})
	}
}
