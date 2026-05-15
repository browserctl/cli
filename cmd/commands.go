// commands imports all subcommand packages to trigger their init() functions,
// which register the commands with RootCmd.
package cmd

import (
	_ "browserctl/cli/cmd/attach"
	_ "browserctl/cli/cmd/back"
	_ "browserctl/cli/cmd/click"
	_ "browserctl/cli/cmd/close"
	_ "browserctl/cli/cmd/cookies"
	_ "browserctl/cli/cmd/eval"
	_ "browserctl/cli/cmd/fill"
	_ "browserctl/cli/cmd/find"
	_ "browserctl/cli/cmd/forward"
	_ "browserctl/cli/cmd/hover"
	_ "browserctl/cli/cmd/html"
	_ "browserctl/cli/cmd/navigate"
	_ "browserctl/cli/cmd/new"
	_ "browserctl/cli/cmd/reload"
	_ "browserctl/cli/cmd/screenshot"
	_ "browserctl/cli/cmd/scroll"
	_ "browserctl/cli/cmd/switch"
	_ "browserctl/cli/cmd/tabs"
	_ "browserctl/cli/cmd/typeinput"
	_ "browserctl/cli/cmd/url"
	_ "browserctl/cli/cmd/version"
)