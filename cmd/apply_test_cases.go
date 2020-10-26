package cmd

var applyWithServerConfigsSetupCases = []struct {
	desc         string
	flags        map[string]string
	route        string
	responseCode int
	response     map[string]string
	expected     bool
}{
	{
		desc: "basic resource apply with server configs pre-configured",
		flags: map[string]string{
			"port":  "8070",
			"host":  "localhost",
			"auth":  "false",
			"token": "",
		},
		expected:     true,
		route:        "/alpha1/upload",
		responseCode: 200,
		response: map[string]string{
			"status": "uploaded",
		},
	},
}
