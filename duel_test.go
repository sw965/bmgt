package bmgt

import (
  "testing"
)

func TestDraw(t *testing.T) {
  duel := Duel{}
  duel.SelfHand = Cards{*CARDS["封印されしエクゾディア"], *CARDS["封印されし者の左腕"], *CARDS["封印されし者の左足"]}
  duel.SelfDeck = Cards{*CARDS["封印されし者の右腕"], *CARDS["封印されし者の右足"]}

  duel, err := duel.Draw(2)
  if err != nil{
    panic(err)
  }
  //
  // fmt.Println("selfHand")
  // fmt.Println(duel.SelfHand)
  // fmt.Println("selfDeck")
  // fmt.Println(duel.SelfDeck)
}
