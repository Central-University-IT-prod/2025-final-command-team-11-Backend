package main

import "REDACTED/team-11/backend/admin/internal/app"

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	a := app.New()

	a.Run()
}
