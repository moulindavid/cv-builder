package resume

type Localized map[string]string

func (l Localized) Get(lang string) string { return l[lang] }

type Resume struct {
	Profile    Profile      `yaml:"profile"`
	Contact    Contact      `yaml:"contact"`
	Summaries  Localized    `yaml:"summaries"`
	Skills     []Skill      `yaml:"skills"`
	Experience []Experience `yaml:"experience"`
	Education  []Education  `yaml:"education"`
	Languages  []Language   `yaml:"languages"`
	Links      []Link       `yaml:"links"`
	Tailoring  Tailoring    `yaml:"tailoring"`
}

type Profile struct {
	Name   string    `yaml:"name"`
	Titles Localized `yaml:"titles"`
}
type Contact struct {
	Email    string `yaml:"email"`
	Phone    string `yaml:"phone"`
	Location string `yaml:"location"`
	Address  string `yaml:"address"`
}
type Skill struct {
	ID       string   `yaml:"id"`
	Name     string   `yaml:"name"`
	Category string   `yaml:"category"`
	Tags     []string `yaml:"tags"`
	Aliases  []string `yaml:"aliases"`
}
type Bullet struct {
	ID   string    `yaml:"id"`
	Text Localized `yaml:"text"`
	Tags []string  `yaml:"tags"`
}
type Experience struct {
	ID           string    `yaml:"id"`
	Company      string    `yaml:"company"`
	Location     string    `yaml:"location"`
	Roles        Localized `yaml:"roles"`
	Start        string    `yaml:"start"`
	End          string    `yaml:"end"`
	Descriptions Localized `yaml:"descriptions"`
	Technologies []string  `yaml:"technologies"`
	Bullets      []Bullet  `yaml:"bullets"`
}
type Education struct {
	ID          string    `yaml:"id"`
	Institution string    `yaml:"institution"`
	Degrees     Localized `yaml:"degrees"`
	Start       string    `yaml:"start"`
	End         string    `yaml:"end"`
}
type Language struct {
	Name  Localized `yaml:"name"`
	Level Localized `yaml:"level"`
}
type Link struct {
	Label string `yaml:"label"`
	URL   string `yaml:"url"`
}
type Tailoring struct {
	ExternalKeywords []string `yaml:"external_keywords"`
}

type Target struct {
	Name       string    `yaml:"name"`
	Titles     Localized `yaml:"titles"`
	Priorities []string  `yaml:"priorities"`
	MaxBullets int       `yaml:"max_bullets"`
}
