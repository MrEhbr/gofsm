package fsm

import (
	"math/rand"
)

var availableColors = []string{
	"aliceblue", "antiquewhite", "aqua", "aquamarine", "azure", "bisque", "blanchedalmond", "blue", "brown", "burlywood", "cadetblue", "chartreuse", "chocolate", "coral", "cornflowerblue", "cornsilk", "crimson", "cyan", "deeppink", "deepskyblue", "dodgerblue", "firebrick", "floralwhite", "forestgreen", "fuchsia", "gainsboro", "ghostwhite", "gold", "goldenrod", "greenyellow", "honeydew", "hotpink", "indianred", "ivory", "khaki", "lavender", "lawngreen", "lemonchiffon", "lime", "limegreen", "linen", "magenta", "mediumaquamarine", "mediumblue", "mediumorchid", "mediumpurple", "mediumseagreen", "mediumslateblue", "mediumspringgreen", "mediumturquoise", "mediumvioletred", "mintcream", "mistyrose", "moccasin", "navajowhite", "oldlace", "olivedrab", "orange", "orangered", "orchid", "palegoldenrod", "palegreen", "paleturquoise", "palevioletred", "papayawhip", "peachpuff", "peru", "pink", "plum", "powderblue", "red", "rosybrown", "royalblue", "saddlebrown", "salmon", "sandybrown", "seagreen", "sienna", "silver", "skyblue", "slateblue", "slategrey", "snow", "springgreen", "steelblue", "tan", "thistle", "tomato", "turquoise", "violet", "wheat", "white", "yellowgreen",
}

type Colors struct {
	values []string
}

func NewColors() *Colors {
	values := make([]string, len(availableColors))
	copy(values, availableColors)
	rand.Shuffle(len(values), func(i, j int) {
		values[i], values[j] = values[j], values[i]
	})

	return &Colors{
		values: values,
	}
}

func (c *Colors) Pick() string {
	if len(c.values) == 0 {
		c.values = make([]string, len(availableColors))
		copy(c.values, availableColors)
	}

	if len(c.values) == 0 {
		return ""
	}

	i := rand.Intn(len(c.values))
	val := c.values[i]

	c.values[i] = c.values[len(c.values)-1]
	c.values = c.values[:len(c.values)-1]
	return val
}
