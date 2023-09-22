package bmgt

type Cost StateChangers

var ZERO_COST = Cost{}

type Costs []Cost

type CostsMaker func(*Action) Costs

// https://www.db.yugioh-card.com/yugiohdb/card_search.action?ope=2&cid=5658&request_locale=ja
func NewRoyalMagicalLibraryCosts(action *Action) Costs {
	cost0 := ZERO_COST
	cost1 := Cost{
		StateChangerF.MonsterZoneSpellCounterRemoval(action.MonsterZoneIndices[0], -1),
	}
	return Costs{cost0, cost1} 
}

// https://www.db.yugioh-card.com/yugiohdb/card_search.action?ope=2&cid=4861&request_locale=ja
func NewSolemnJudgmentCosts(action *Action) Costs {
	cost0 := Cost{
		StateChangerF.PayHalfLifePoint,
	}
	return Costs{cost0}
}

// https://www.db.yugioh-card.com/yugiohdb/card_search.action?ope=2&cid=6400&request_locale=ja
func NewSummonerMonkCosts(action *Action) Costs {
	cost0 := ZERO_COST
	cost1 := ZERO_COST
	cost2 := Cost{
		StateChangerF.Discard(action.HandIndices),
	}
	return Costs{cost0, cost1, cost2}
}

// https://www.db.yugioh-card.com/yugiohdb/card_search.action?ope=2&cid=4431&request_locale=ja
func NewThunderDragonCosts(action *Action) Costs {
	cost0 := Cost{
		StateChangerF.Discard(action.HandIndices),
	}
	return Costs{cost0}
}

// https://www.db.yugioh-card.com/yugiohdb/card_search.action?ope=2&cid=5688&request_locale=ja
func NewMagicalStoneExcavationCosts(action *Action) Costs {
	cost0 := Cost{
		StateChangerF.Discard(action.HandIndices),
	}
	return Costs{cost0}
}

var COSTS = map[CardName]CostsMaker{
	ROYAL_MAGICAL_LIBRARY:NewRoyalMagicalLibraryCosts,
	SOLEMN_JUDGMENT:NewSolemnJudgmentCosts,
	SUMMONER_MONK:NewSummonerMonkCosts,
	THUNDER_DRAGON:NewThunderDragonCosts,
	MAGICAL_STONE_EXCAVATION:NewMagicalStoneExcavationCosts,
}
