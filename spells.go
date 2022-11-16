package bmgt

type Spell struct {
  IsQuickPlay bool
}

var SPELLS = func() map[string]*Spell {
  result := map[string]*Spell{}

  //あ行
  //か行
  result["強欲で謙虚な壺"] = &Spell{
    IsQuickPlay:false,
  }

  //さ行
  //た行
  result["トゥーンのもくじ"] = &Spell{
    IsQuickPlay:false,
  }

  //な行
  result["成金ゴブリン"] = &Spell{
    IsQuickPlay:false,
  }
  //は行

  //ま行
  result["魔法石の採掘"] = &Spell{
    IsQuickPlay:false,
  }

  //や行

  // result[""] = Spell{
  //   IsQuickPlay:false,
  // }
  return result
}()
