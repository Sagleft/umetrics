package frontend

const (
	maxTopChannels          = 10
	maxTopUsers             = 10
	maxTopOwners            = 10
	activeUsersDaysInterval = 7
)

var (
	frontFilesToHash = []string{
		"assets/css/custom.css",
		"assets/js/globe.js",
		"assets/img/favicon.png",
	}
)
