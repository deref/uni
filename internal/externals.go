package internal

func getExternals(repo *Repository) []string {
	externals := make([]string, len(repo.Dependencies))
	i := 0
	for external := range repo.Dependencies {
		externals[i] = external
		i++
	}
	return externals
}
