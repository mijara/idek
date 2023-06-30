package app

type Headers struct {
	XAPIToken     string `header:"x-api-token"`
	XCustomHeader string `header:"x-custom-header"`
}
