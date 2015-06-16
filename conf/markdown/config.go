package markdown_config

type ConfigValue string
type Config map[string][]ConfigValue

type Configs map[string]Config

func (s Configs) String() string {
	return ToJson(s)
}
