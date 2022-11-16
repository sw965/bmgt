package bmgt

const MONSTER_ZONE_LENGTH = 5
type MonsterZone [MONSTER_ZONE_LENGTH]Card

const SPELL_AND_TRAP_ZONE_LENGTH = 5
type SpellAndTrapZone [SPELL_AND_TRAP_ZONE_LENGTH]Card

type Phase string

const (
  DRAW_PHASE = Phase("ドローフェイズ")
  STANDBY_PHASE = Phase("スタンバイフェイズ")
  MAIN_PHASE1 = Phase("メインフェイズ1")
  BATTLE_PHASE = Phase("バトルフェイズ")
  MAIN_PHASE2 = Phase("メインフェイズ2")
  END_PHASE = Phase("エンドフェイズ")
)

const (
  MIN_MAIN_DECK_NUM = 40
  MAX_MAIN_DECK_NUM = 60
  MAX_EX_DECK_NUM = 15
)

type Playing struct {
  IsPotOfDuality bool
  IsAllureOfDarkness bool
}

type Duel struct {
  SelfLifePoint int
  OpponentLifePoint int

  SelfDeck Cards
  OpponentDeck Cards

  SelfMonsterZone MonsterZone
  SelfSpellAndTrapZone SpellAndTrapZone

  OpponentMonsterZone MonsterZone
  OpponentSpellAndTrapZone SpellAndTrapZone

  SelfHand Cards
  SelfGraveyard Cards
  SelfBanish Cards

  OpponentHand Cards
  OpponentGraveyard Cards
  OpponentBanish Cards

  TurnNum int
  Phase string

  SelectCards Cards
  SelectedCard Card

  SelfPlayer func(Duel, *Playing) Duel
  OpponentPlayer func(Duel, *Playing) Duel
}

func (duel Duel) Draw(num int) (Duel, error) {
  selfDeck, drawCards, err := duel.SelfDeck.Draw(num)
  duel.SelfDeck = selfDeck
  duel.SelfHand = duel.SelfHand.Add(drawCards)
  return duel, err
}
