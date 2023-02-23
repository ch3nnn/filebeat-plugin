package actions

type config struct {
	FileHasSuffix string `config:"file_has_suffix" validate:"required"`
}
