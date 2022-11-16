package bmgt

import (
  "fmt"
  "math/rand"
)

type Card struct {
  CardName string
  IsMonster bool
  IsSpell bool
  MonsterLevel int
  MonsterAttribute Attribute
  MonsterType Type
  ATK int
  DEF int
  Show bool
}

func (card1 Card) Equal(card2 Card) bool {
  if card1.CardName != card2.CardName {
    return false
  }

  if card1.IsMonster != card2.IsMonster {
    return false
  }

  if card1.IsSpell != card2.IsSpell {
    return false
  }

  if card1.MonsterLevel != card2.MonsterLevel {
    return false
  }

  if card1.MonsterAttribute != card2.MonsterAttribute {
    return false
  }

  if card1.MonsterType != card2.MonsterType {
    return false
  }

  if card1.ATK != card2.ATK {
    return false
  }

  if card1.DEF != card2.DEF {
    return false
  }

  if card1.Show != card2.Show {
    return false
  }
  
  return true
}

func IsDarkMonster(card *Card) bool {
  return card.MonsterAttribute == DARK
}

var CARDS = func() map[string]*Card {
  result := map[string]*Card{}

  for cardName, monster := range MONSTERS {
    card := &Card{
      CardName:cardName, IsMonster:true, IsSpell:false,
      MonsterLevel:monster.Level, MonsterAttribute:monster.Attribute, MonsterType:monster.Type,
      ATK:monster.ATK, DEF:monster.DEF, Show:false,
    }
    result[cardName] = card
  }

  for cardName, _ := range SPELLS {
    card := &Card{
      CardName:cardName, IsMonster:false, IsSpell:true,
      MonsterAttribute:"", MonsterType:"",
      ATK:-1, DEF:-1,
    }
    result[cardName] = card
  }
  return result
}()

type Cards []Card

func (cards Cards) Copy() Cards {
  result := make(Cards, len(cards))
  for i, card := range cards {
    result[i] = card
  }
  return result
}

func (cards1 Cards) Add(cards2 Cards) Cards {
  result := make(Cards, 0, len(cards1) + len(cards2))
  for _, card := range cards1 {
    result = append(result, card)
  }

  for _, card := range cards2 {
    result = append(result, card)
  }
  return result
}

func (deck Cards) Draw(num int) (Cards, Cards, error) {
  if len(deck) < num {
    return Cards{}, Cards{}, fmt.Errorf("ドロー枚数 > カード枚数 になっている")
  }

  newDeck := make(Cards, len(deck) - num)
  for i := 0; i < len(newDeck); i++ {
    newDeck[i] = deck[num + i]
  }

  drawCards := make(Cards, 0, num)
  for i := 0; i < num; i++ {
    drawCards = append(drawCards, deck[i])
  }
  return newDeck, drawCards, nil
}

func (cards Cards) Shuffle(random *rand.Rand) Cards {
  result := cards.Copy()
	for i := len(result); i > 1; i-- {
		j := random.Intn(i)
		result[i - 1], result[j] = result[j], result[i - 1]
	}
  return result
}

func (cards Cards) Filter(f func(*Card) bool) Cards {
  result := make(Cards, 0)
  for _, card := range cards {
    if f(&card) {
      result = append(result, card)
    }
  }
  return result
}

func (cards Cards) Remove(card Card) Cards {
  result := make(Cards, 0, len(cards) - 1)
  removeOk := false

  for _, iCard := range cards {
    if !iCard.Equal(card) {
      result = append(result, iCard)
    } else if removeOk {
      result = append(result, iCard)
      removeOk = true
    }
  }
  return result
}

func (cards Cards) DropStart(num int) Cards {
  resultLength := len(cards) - num
  result := make(Cards, 0, resultLength)
  for i := 0; i < resultLength; i++ {
    result[i] = cards[num + i]
  }
  return result
}
