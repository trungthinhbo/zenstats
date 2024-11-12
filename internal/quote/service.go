package quote

import (
	"math/rand/v2"
)

type Quote struct {
	Author string
	Quote  string
}

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (m *Service) GetQuote() Quote {
	quotes := []Quote{
		{
			Author: "Dreams of Code",
			Quote:  "Good things come to those who write Go",
		},
		{
			Author: "Wei Wu Wei",
			Quote:  "Wise men don’t judge – they seek to understand.",
		},
		{
			Quote:  "The noble-minded are calm and steady. Little people are forever fussing and fretting.",
			Author: "Confucius",
		},
		{
			Quote:  "The successful warrior is the average man, with laser-like focus.",
			Author: "Bruce Lee",
		},
		{
			Quote:  "The first and the best victory is to conquer self.",
			Author: "Plato",
		},
		{
			Quote:  "To be truly ignorant, be content with your own knowledge.",
			Author: "Zhuangz",
		},
		{
			Quote:  "One must be deeply aware of the impermanence of the world.",
			Author: "Dogen",
		},
		{
			Quote:  "In the midst of chaos, there is also opportunity.",
			Author: "Sun Tzu",
		},
	}

	return quotes[rand.IntN(len(quotes))]
}
