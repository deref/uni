package internal

type Config struct {
	// TODO: Do not load these in as maps, since duplicate keys are
	// not decectable this way!
	Engines      map[string]string
	Packages     map[string]PackageConfig
	Dependencies map[string]string
}

type PackageConfig struct {
	Public      bool
	Description string
	Entrypoint  string
	// TODO: License, author, etc.
}
