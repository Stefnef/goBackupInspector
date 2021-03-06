package summary

import (
	"reflect"
	"time"
)

type Difference int

const (
	Equal Difference = iota
	LeftDir
	RightDir
	FilesNotInDir
	DirectoriesNotInDir
	ComparedFiles
	UnequalFiles
	IgnoredFiles
	WithDifferences
	BackupFiles
)

func (d Difference) i() int {
	return int(d)
}

type FileTuple struct {
	LeftFile  string `json:"left_file"`
	RightFile string `json:"right_file"`
}

type IgnoredElement struct {
	IgnoredElement string `json:"IgnoredElement"`
	CausedRule     string `json:"CausedRule"`
}

type FileDiffSummary struct {
	Date                time.Time           `json:"Date"`
	LeftDir             string              `json:"LeftDir"`
	RightDir            string              `json:"RightDir"`
	FilesNotInDir       map[string][]string `json:"FilesNotInDir"`
	DirectoriesNotInDir map[string][]string `json:"DirectoriesNotInDir"`
	ComparedFiles       []FileTuple         `json:"ComparedFiles"`
	IgnoredElement      []IgnoredElement    `json:"IgnoredElements"` //TODO plural
	UnequalFiles        []FileTuple         `json:"UnequalFiles"`
	WithDifferences     bool                `json:"WithDifferences"`
	BackupFileName      string              `json:"BackupFileName"`
}

func (ft FileTuple) Compare(other FileTuple) int {
	if ft.LeftFile != other.LeftFile {
		return -1
	}
	if ft.RightFile != other.RightFile {
		return 1
	}
	return 0
}

func (ft FileTuple) String() string {
	return "{LF '" + ft.LeftFile + "', RF '" + ft.RightFile + "'}"
}

func Compare(sum, other []FileTuple) int {
	if len(sum) < len(other) {
		return -1
	} else if len(sum) > len(other) {
		return 1
	}
	for idx, ftElement := range sum {
		if ftElement.Compare(other[idx]) != 0 {
			return idx + 1
		}
	}
	return 0
}

func (ignore IgnoredElement) String() string {
	return "{" + ignore.IgnoredElement + " [" + ignore.CausedRule + "]" + "}"
}

func (ignore IgnoredElement) Compare(other IgnoredElement) int {
	if ignore.IgnoredElement != other.IgnoredElement {
		return -1
	}
	if ignore.CausedRule != other.CausedRule {
		return 1
	}
	return 0
}

func CompareIgnoredElements(act, exp []IgnoredElement) int {
	if len(act) != len(exp) {
		return -1
	}

	for idx, actElem := range act {
		if actElem.Compare(exp[idx]) != 0 {
			return idx + 1
		}
	}
	return 0
}

func (sum FileDiffSummary) Compare(other FileDiffSummary) int {
	if sum.LeftDir != other.LeftDir {
		return LeftDir.i()
	}
	if sum.RightDir != other.RightDir {
		return RightDir.i()
	}
	if eq := reflect.DeepEqual(sum.FilesNotInDir, other.FilesNotInDir); !eq {
		return FilesNotInDir.i()
	}
	if eq := reflect.DeepEqual(sum.DirectoriesNotInDir, other.DirectoriesNotInDir); !eq {
		return DirectoriesNotInDir.i()
	}
	if Compare(sum.ComparedFiles, other.ComparedFiles) != 0 {
		return ComparedFiles.i()
	}
	if Compare(sum.UnequalFiles, other.UnequalFiles) != 0 {
		return UnequalFiles.i()
	}
	if CompareIgnoredElements(sum.IgnoredElement, other.IgnoredElement) != 0 {
		return IgnoredFiles.i()
	}
	if sum.WithDifferences != other.WithDifferences {
		return WithDifferences.i()
	}
	if sum.BackupFileName != other.BackupFileName {
		return BackupFiles.i()
	}
	return Equal.i()
}

func (sum FileDiffSummary) HasDifferences() bool {
	if len(sum.DirectoriesNotInDir[sum.LeftDir]) != 0 ||
		len(sum.DirectoriesNotInDir[sum.RightDir]) != 0 {
		return true
	}
	if len(sum.UnequalFiles) != 0 {
		return true
	}
	if len(sum.FilesNotInDir[sum.LeftDir]) != 0 ||
		len(sum.FilesNotInDir[sum.RightDir]) != 0 {
		return true
	}
	return false
}
