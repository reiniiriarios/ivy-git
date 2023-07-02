package git

import (
	"errors"
	"sort"
)

type AllLanguageData []*LanguageData

type ClocData struct {
	LastHashParsed string
	Languages      AllLanguageData
	Total          *LanguageData
}

type ClocProcessor struct {
	langs *DefinedLanguages
}

type ClocResult struct {
	Total         *LanguageData
	Languages     map[string]*Language
	MaxPathLength int
}

func (ls AllLanguageData) Len() int {
	return len(ls)
}

func (ls AllLanguageData) Swap(i, j int) {
	ls[i], ls[j] = ls[j], ls[i]
}

func (ls AllLanguageData) Less(i, j int) bool {
	if ls[i].Bytes == ls[j].Bytes {
		return ls[i].Name < ls[j].Name
	}
	return ls[i].Bytes > ls[j].Bytes
}

func (g *Git) Cloc() (ClocData, error) {
	if g.Repo.Directory == "" {
		return ClocData{}, errors.New("no repo selected")
	}
	if g.Repo.Main == "" {
		return ClocData{}, errors.New("no main branch set for repo")
	}

	paths, err := g.LsTreeBranch(g.Repo.Main)
	if err != nil {
		return ClocData{}, err
	}

	// Get data
	processor := ClocProcessor{
		langs: NewDefinedLanguages(),
	}
	gitattributes_translations := g.parseGitAttributes()

	result, err := g.clocAnalyze(&processor, paths, gitattributes_translations)
	if err != nil {
		return ClocData{}, err
	}

	data := ClocData{
		LastHashParsed: g.lastCommitOnMain(),
		Total:          result.Total,
	}

	// Calc percentages and cort data
	for _, l := range result.Languages {
		if l.Data.Files > 0 {
			if result.Total.Code > 0 {
				l.Data.CodePercent = float64(l.Data.Code) / float64(result.Total.Code) * 100
			}
			if result.Total.Bytes > 0 {
				l.Data.TotalPercent = float64(l.Data.Bytes) / float64(result.Total.Bytes) * 100
			}
			data.Languages = append(data.Languages, &l.Data)
		}
	}
	sort.Sort(data.Languages)

	return data, nil
}

// Analyze executes gocloc parsing for the directory of the paths argument and returns the result.
func (g *Git) clocAnalyze(p *ClocProcessor, files []string, translations map[string]string) (*ClocResult, error) {
	total := LanguageData{}
	languages := g.clocAllFiles(files, p.langs, translations)
	maxPathLen := 0
	num := 0
	for _, lang := range languages {
		num += len(lang.Files)
		for _, file := range lang.Files {
			l := len(file)
			if maxPathLen < l {
				maxPathLen = l
			}
		}
	}
	clocFiles := make(map[string]*ClocFile, num)

	for _, language := range languages {
		for _, file := range language.Files {
			cf, _ := g.clocAnalyzeFileOnBranch(file, language)
			cf.Lang = language.Data.Name

			language.Data.Code += cf.Code
			language.Data.Comments += cf.Comments
			language.Data.Blanks += cf.Blanks
			language.Data.Total += cf.Code + cf.Comments + cf.Blanks
			language.Data.Bytes += cf.Bytes
			clocFiles[file] = cf
		}

		language.Data.Files = int64(len(language.Files))
		if language.Data.Files <= 0 {
			continue
		}

		total.Files += language.Data.Files
		total.Blanks += language.Data.Blanks
		total.Comments += language.Data.Comments
		total.Code += language.Data.Code
		total.Total += language.Data.Total
		total.Bytes += language.Data.Bytes
	}

	return &ClocResult{
		Total:         &total,
		Languages:     languages,
		MaxPathLength: maxPathLen,
	}, nil
}
