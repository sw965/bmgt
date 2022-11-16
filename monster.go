package bmgt

type Type string

const (
  SPELLCASTER = Type("魔法使い族")
  MACHINE = Type("機械")
)

type Attribute string

const (
  DARK = "闇属性"
)

type Monster struct {
  Level int
  Attribute Attribute
  Type Type
  ATK int
  DEF int
}

var MONSTERS = func() map[string]*Monster {
  result := map[string]*Monster{}

  //は行
  result["封印されしエクゾディア"] = &Monster{
    Level:3,
    Attribute:DARK,
    Type:SPELLCASTER,
    ATK:1000,
    DEF:1000,
  }

  result["封印されし者の左足"] = &Monster{
    Level:1,
    Attribute:DARK,
    Type:SPELLCASTER,
    ATK:200,
    DEF:200,
  }

  result["封印されし者の左腕"] = &Monster{
    Level:1,
    Attribute:DARK,
    Type:SPELLCASTER,
    ATK:200,
    DEF:200,
  }

  result["封印されし者の右足"] = &Monster{
    Level:1,
    Attribute:DARK,
    Type:SPELLCASTER,
    ATK:200,
    DEF:200,
  }

  result["封印されし者の右腕"] = &Monster{
    Level:1,
    Attribute:DARK,
    Type:SPELLCASTER,
    ATK:200,
    DEF:200,
  }

  //ま行

  // result[""] = Monster{
  //   Level:,
  //   Attribute:,
  //   Type:,
  //   ATK:,
  //   DEF:,
  // }
  return result
}()
