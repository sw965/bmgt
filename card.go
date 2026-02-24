package bmgt

type CardName int

type Type int

const (
	Spellcaster  Type = iota // 魔法使い族
	Dragon                   // ドラゴン族
	Zombie                   // アンデット族
	Warrior                  // 戦士族
	BeastWarrior             // 獣戦士族
	Beast                    // 獣族
	WingedBeast              // 鳥獣族
	Fiend                    // 悪魔族
	Fairy                    // 天使族
	Insect                   // 昆虫族
	Dinosaur                 // 恐竜族
	Reptile                  // 爬虫類族
	Fish                     // 魚族
	SeaSerpent               // 海竜族
	Aqua                     // 水族
	Pyro                     // 炎族
	Thunder                  // 雷族
	Rock                     // 岩石族
	Plant                    // 植物族
	Machine                  // 機械族
	Psychic                  // サイキック族
	Wyrm                     // 幻竜族
	Cyberse                  // サイバース族
	Illusion                 // 幻想魔族 (第12期で追加)
	DivineBeast              // 幻神獣族 (三幻神など)
	CreatorGod               // 創造神族
)

type Attribute int

const (
	Dark   Attribute = iota // 闇属性
	Light                   // 光属性
	Earth                   // 地属性
	Water                   // 水属性
	Fire                    // 炎属性
	Wind                    // 風属性
	Divine                  // 神属性
)

type CardCategory int

const (
	MonsterCard CardCategory = iota
	SpellCard
	TrapCard
)

type Card struct {
	Name CardName
	Atk  int
	Def  int
	Id   int
}

type Cards []Card