// Copyright 2025 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package safelogging

import (
	"cmp"
	"go/ast"
	"go/token"
	"maps"
	"slices"
	"strings"

	"github.com/palantir/safe-logging-go/safelogging/internal/filecomments"
)

type storedLogSafetyVal struct {
	logSafety  LogSafetyType
	commentPos token.Pos
}

type CommentBasedLogSafetyTracker struct {
	// map from filename -> line number of identifier(s) -> storedLogSafetyVal for that line as determined by comment
	identFileLineToLogSafetyMap map[string]map[int]storedLogSafetyVal

	// the FileSet used to compute the log safety map
	fset *token.FileSet
}

// LogSafetyForIdentAtPos returns the LogSafetyType for the identifier at the given position based on comment-based
// annotations. Comment-based log safety annotations are defined/specified based on line number, and the provided
// token.Pos is used only to extract the line number. As such, different token.Pos values that refer to the same line
// will return the same LogSafetyType (which is the expected/desired behavior in the case where there are multiple
// identifiers that occur on the same line).
func (t *CommentBasedLogSafetyTracker) LogSafetyForIdentAtPos(pos token.Pos) (LogSafetyType, token.Position) {
	file := t.fset.File(pos)
	if file == nil {
		return LogSafetyTypeUnmarked, token.Position{}
	}
	lineNum := file.Line(pos)
	lineToLogSafetyMap, ok := t.identFileLineToLogSafetyMap[file.Name()]
	if !ok {
		return LogSafetyTypeUnmarked, token.Position{}
	}
	storedLogSafety, ok := lineToLogSafetyMap[lineNum]
	if !ok {
		return LogSafetyTypeUnmarked, token.Position{}
	}
	return storedLogSafety.logSafety, file.Position(storedLogSafety.commentPos)
}

func NewCommentBasedLogSafetyTracker[T any](allIdentifiers map[*ast.Ident]T, fset *token.FileSet, files []*ast.File) *CommentBasedLogSafetyTracker {
	logSafetyTracker := &CommentBasedLogSafetyTracker{
		identFileLineToLogSafetyMap: make(map[string]map[int]storedLogSafetyVal),
		fset:                        fset,
	}

	// sort all identifiers by their position
	sortedKeys := slices.SortedFunc(maps.Keys(allIdentifiers), func(a, b *ast.Ident) int {
		return cmp.Compare(a.Pos(), b.Pos())
	})

	fileTokenToASTFile := make(map[*token.File]*ast.File)
	for _, currFile := range files {
		fileTokenForCurrFile := fset.File(currFile.Pos())
		fileTokenToASTFile[fileTokenForCurrFile] = currFile
		logSafetyTracker.identFileLineToLogSafetyMap[fileTokenForCurrFile.Name()] = make(map[int]storedLogSafetyVal)
	}

	// variables are initialized and populated each time a new file is encountered.
	// Files are encountered sequentially because identifiers are sorted by position.
	var (
		currProcessingFile     *token.File
		lineToLogSafetyTypeMap map[int]storedLogSafetyVal
	)

	for _, ident := range sortedKeys {
		tokenFile := fset.File(ident.Pos())
		if currProcessingFile != tokenFile {
			currProcessingFile = tokenFile

			lineToLogSafetyTypeMap = make(map[int]storedLogSafetyVal)
			lineToCommentMap := filecomments.CreateLineToCommentMap(tokenFile, fileTokenToASTFile[tokenFile], safeLoggingCommentPrefix)
			for _, k := range slices.Sorted(maps.Keys(lineToCommentMap)) {
				v := lineToCommentMap[k]
				commentContent := strings.TrimPrefix(v.Text, safeLoggingCommentPrefix)

				// special case: to support annotations and analysistest expectations, if a line that matches
				// safeLoggingCommentPrefix contains another section that starts with "//", only consider the part
				// before the "//" as content
				if nextSlashIdx := strings.Index(commentContent, "//"); nextSlashIdx != -1 {
					commentContent = commentContent[:nextSlashIdx]
				}

				logSafety := toLogSafetyType(strings.TrimSpace(commentContent))
				if logSafety == LogSafetyTypeUnmarked {
					continue
				}
				lineToLogSafetyTypeMap[k] = storedLogSafetyVal{
					logSafety:  logSafety,
					commentPos: v.Slash,
				}
			}
		}

		identLineNum := tokenFile.Line(ident.Pos())
		lineNumberUsed := identLineNum

		// check if log safety specified on same line as identifier
		logSafetyVal, ok := lineToLogSafetyTypeMap[lineNumberUsed]

		// if not, check if log safety specified on the line before
		if !ok {
			lineNumberUsed--
			logSafetyVal, ok = lineToLogSafetyTypeMap[lineNumberUsed]
		}

		// log safety not specified on the line of the identifier or line before: nothing to do
		if !ok {
			continue
		}

		// log safety value was located: record in map and remove the line from consideration for future identifiers
		logSafetyTracker.identFileLineToLogSafetyMap[currProcessingFile.Name()][identLineNum] = logSafetyVal
		delete(lineToLogSafetyTypeMap, lineNumberUsed)
	}
	return logSafetyTracker
}
