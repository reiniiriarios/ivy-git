package cloc

import (
	"sort"
)

type AllLanguageData []*LanguageData

type ClocData struct {
	Languages AllLanguageData
	Total     *LanguageData
}

func (ls AllLanguageData) Len() int {
	return len(ls)
}

func (ls AllLanguageData) Swap(i, j int) {
	ls[i], ls[j] = ls[j], ls[i]
}

func (ls AllLanguageData) Less(i, j int) bool {
	if ls[i].Code == ls[j].Code {
		return ls[i].Name < ls[j].Name
	}
	return ls[i].Code > ls[j].Code
}

func Cloc(paths ...string) (ClocData, error) {
	// Get data
	processor := Processor{
		langs: NewDefinedLanguages(),
	}
	result, err := processor.Analyze(paths)
	if err != nil {
		return ClocData{}, err
	}

	data := ClocData{
		Total: result.Total,
	}

	// Calc percentages and cort data
	for _, l := range result.Languages {
		if l.Data.Files > 0 {
			l.Data.CodePercent = float64(l.Data.Code) / float64(result.Total.Code) * 100
			l.Data.TotalPercent = float64(l.Data.Total) / float64(result.Total.Total) * 100
			data.Languages = append(data.Languages, &l.Data)
		}
	}
	sort.Sort(data.Languages)

	return data, nil
}

// Processor is gocloc analyzing processor.
type Processor struct {
	langs *DefinedLanguages
}

// Result defined processing result.
type Result struct {
	Total         *LanguageData
	Languages     map[string]*Language
	MaxPathLength int
}

// Analyze executes gocloc parsing for the directory of the paths argument and returns the result.
func (p *Processor) Analyze(files []string) (*Result, error) {
	total := LanguageData{}
	languages := parseAllFiles(files, p.langs)
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
			cf := AnalyzeFile(file, language)
			cf.Lang = language.Data.Name

			language.Data.Code += cf.Code
			language.Data.Comments += cf.Comments
			language.Data.Blanks += cf.Blanks
			language.Data.Total += cf.Code + cf.Comments + cf.Blanks
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
	}

	return &Result{
		Total:         &total,
		Languages:     languages,
		MaxPathLength: maxPathLen,
	}, nil
}
