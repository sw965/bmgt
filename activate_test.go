package bmgt

import (
  "testing"
  "fmt"
  "math/rand"
  "time"
  "github.com/seehuhn/mt19937"
)

func TestPotOfDuality(t *testing.T) {
  mtRandom := rand.New(mt19937.New())
  mtRandom.Seed(time.Now().UnixNano())

  player := func(duel Duel, playing *Playing) Duel {
    if playing.IsPotOfDuality {
      duel.SelectedCard = duel.SelectCards[1]
      return duel
    }
    return duel
  }

  duel := Duel{SelfPlayer:player}
  duel.SelfHand = Cards{*CARDS["封印されしエクゾディア"], *CARDS["封印されし者の右腕"], *CARDS["封印されし者の右足"]}
  duel.SelfDeck = Cards{*CARDS["成金ゴブリン"], *CARDS["封印されし者の左腕"], *CARDS["封印されし者の左足"]}
  duel, err := PotOfDualityActivate(duel, mtRandom)
  if err != nil {
    panic(err)
  }
  fmt.Println("selfHand")
  fmt.Println(duel.SelfHand)
  fmt.Println("selfDeck")
  fmt.Println(duel.SelfDeck)
}
