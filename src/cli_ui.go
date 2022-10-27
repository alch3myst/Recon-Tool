package recon

func UiBanner() {
	banner := `
    ____  ____________
   / __ \/ ____/ ____/
  / /_/ / /   / /_    
 / _, _/ /___/ __/    
/_/ |_|\____/_/       `

	print("RCF Toolkit v1.0")
	println(banner)
}

func UiHelp() {
	println("\n-- Usage")
	println("-t	Target url")
	println("-m	Mode (recon, ...)")
	println("-http	Set the request to http")
	print("\n")
}
