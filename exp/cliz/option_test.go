package cliz

type unknownOptionType struct {
	Name        string
	Aliases     []string
	Environment string
	Default     string
	Required    bool
	Description string
	value       *string
}

func (o *unknownOptionType) GetName() string         { return o.Name }
func (o *unknownOptionType) GetAliases() []string    { return o.Aliases }
func (o *unknownOptionType) GetEnvironment() string  { return o.Environment }
func (o *unknownOptionType) GetDefault() interface{} { return o.Default }
func (o *unknownOptionType) IsRequired() bool        { return o.Required }
func (o *unknownOptionType) IsZero() bool            { return o.value == nil || *o.value == "" }
func (o *unknownOptionType) GetDescription() string {
	if o.Description != "" {
		return o.Description
	}
	return "string value of " + o.Name
}
