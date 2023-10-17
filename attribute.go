package bmgt

import (
	"github.com/sw965/omw/fn"
	omwmaps "github.com/sw965/omw/maps"
)

type Attribute int

const (
	DARK Attribute = iota
	LIGHT
	EARTH
	WATER
	FIRE
	WIND
)

type Attributes []Attribute

var ATTRIBUTES = Attributes{DARK, LIGHT, EARTH, WATER, FIRE, WIND}

func AttributeToString(attribute Attribute) string {
	switch attribute {
		case DARK:
			return "闇"
		case LIGHT:
			return "光"
		case EARTH:
			return "地"
		case WATER:
			return "水"
		case FIRE:
			return "炎"
		case WIND:
			return "風"
		default:
			return ""
	}
}

var ATTRIBUTE_TO_STRING = fn.Memo[map[Attribute]string](ATTRIBUTES, AttributeToString)
var STRING_TO_ATTRIBUTE = omwmaps.Reverse[map[string]Attribute](ATTRIBUTE_TO_STRING)