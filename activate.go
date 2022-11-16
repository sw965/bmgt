package bmgt

import (
  "math/rand"
)

//成金ゴブリン
func UpstartGoblinActivate(duel Duel, random *rand.Rand) (Duel, error) {
  duel, err := duel.Draw(1)
  duel.OpponentLifePoint += 1000
  return duel, err
}

//強欲で謙虚な壺
func PotOfDualityActivate(duel Duel, random *rand.Rand) (Duel, error) {
  n := 3
  duel.SelectCards = duel.SelfDeck[:n]
  duel.SelfDeck = duel.SelfDeck.DropStart(n)
  duel = duel.SelfPlayer(duel, &Playing{IsPotOfDuality:true})

  duel.SelectedCard.Show = true
  //ここでコピーされるので、SelfHandの Show = true は保たれたまま
  duel.SelfHand = duel.SelfHand.Add(Cards{duel.SelectedCard})
  duel.SelectedCard.Show = false

  returnCards := duel.SelectCards.Remove(duel.SelectedCard)
  duel.SelfDeck = duel.SelfDeck.Add(returnCards)
  duel.SelfDeck = duel.SelfDeck.Shuffle(random)

  duel.SelectCards = Cards{}
  duel.SelectedCard = Card{}
  return duel, nil
}

//闇の誘惑
func AllureOfDarknessActivate(duel Duel, random *rand.Rand) (Duel, error) {
  duel, err := duel.Draw(2)
  handDarkMonsters := duel.SelfHand.Filter(IsDarkMonster)

  //闇属性がなければ、手札を全て除外
  if len(handDarkMonsters) == 0 {
    selfHand := duel.SelfHand
    duel.SelfHand = Cards{}
    duel.SelfBanish = duel.SelfBanish.Add(selfHand)
    return duel, err
  }

  //手札の闇属性モンスターを一枚除外する
  duel.SelectCards = handDarkMonsters
  duel = duel.SelfPlayer(duel, &Playing{IsAllureOfDarkness:true})
  duel.SelectCards = Cards{}

  duel.SelfHand = duel.SelfHand.Remove(duel.SelectedCard)
  duel.SelfBanish = duel.SelfBanish.Add(Cards{duel.SelectedCard})
  duel.SelectedCard = Card{}
  return duel, err
}

var ACTIVATES = map[string]func(Duel, *rand.Rand) (Duel, error) {
  "成金ゴブリン":UpstartGoblinActivate,
  "強欲で謙虚な壺":PotOfDualityActivate,
  "闇の誘惑":AllureOfDarknessActivate,
}
