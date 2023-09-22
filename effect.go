package bmgt

type Effect StateChangers

var ZERO_EFFECT = Effect{}

type Effects []Effect

// https://www.db.yugioh-card.com/yugiohdb/card_search.action?ope=2&cid=9778&request_locale=ja
func NewOneDayOfPeaceEffects(action *Action) Effects {
	effect0 := Effect{
		StateChangerF.Draw(1),
		StateChangerF.ReversePlayer1AndPlayer2,
		StateChangerF.Draw(1),
		StateChangerF.ReversePlayer1AndPlayer2,
		StateChangerF.OneDayOfPeace,
	}
	return Effects{effect0}
}

// https://www.db.yugioh-card.com/yugiohdb/card_search.action?ope=2&cid=6541&request_locale=ja
func NewMagicalMalletEffects(action *Action) Effects {
	effect0 := Effect{
		StateChangerF.HandToDeck(action.HandIndices),
		StateChangerF.Draw(len(action.HandIndices)),
	}
	return Effects{effect0}
}

// https://www.db.yugioh-card.com/yugiohdb/card_search.action?ope=2&cid=5658&request_locale=ja
func NewRoyalMagicalLibraryEffects(action *Action) Effects {
	effect0 := ZERO_EFFECT
	effect1 := Effect{
		StateChangerF.Draw(1),
	}
	return Effects{effect0, effect1}
}

// https://www.db.yugioh-card.com/yugiohdb/card_search.action?ope=2&cid=4861&request_locale=ja
func NewSolemnJudgmentEffects(action *Action) Effects {
	effect0 := Effect{

	}
}