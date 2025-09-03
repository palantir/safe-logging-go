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

package filecomments

import (
	"go/ast"
	"go/token"
	"strings"
)

var _ Retriever = (*fileCommentRetriever)(nil)

type Retriever interface {
	// CommentOnLineRelativeToPos returns the comment on the line that is the line on the provided node combined with
	// the provided delta. For example, a delta value of 0 will return the comment on the node's line, while a value of
	// -1 will return the comment on the line before the node and a value of 1 will return the comment on the line after
	// the node. Returns an empty string if there is no comment on the line. The second return value is the token.Position
	// of the comment's slash (the first character of the comment).
	CommentOnLineRelativeToPos(astNodePos token.Pos, delta int) (string, token.Position)
}

type fileCommentRetriever struct {
	requiredCommentPrefix string

	fset *token.FileSet

	fileTokenToASTFile map[*token.File]*ast.File

	// map from *token.File -> line number of file -> comment on that line number
	fileToLineToComment map[*token.File]map[int]*ast.Comment
}

func NewRetriever(requiredCommentPrefix string, fset *token.FileSet, files []*ast.File) Retriever {
	r := &fileCommentRetriever{
		requiredCommentPrefix: requiredCommentPrefix,
		fset:                  fset,
		fileTokenToASTFile:    make(map[*token.File]*ast.File),
		fileToLineToComment:   make(map[*token.File]map[int]*ast.Comment),
	}
	for _, currFile := range files {
		fileTokenForCurrFile := fset.File(currFile.Pos())
		r.fileTokenToASTFile[fileTokenForCurrFile] = currFile
	}
	return r
}

func CreateLineToCommentMap(tokenFile *token.File, astFile *ast.File, commentPrefix string) map[int]*ast.Comment {
	lineToCommentMap := make(map[int]*ast.Comment)
	for _, commentGroup := range astFile.Comments {
		for _, comment := range commentGroup.List {
			commentSlashLine := tokenFile.Line(comment.Slash)
			if !strings.HasPrefix(comment.Text, commentPrefix) {
				continue
			}
			lineToCommentMap[commentSlashLine] = comment
		}
	}
	return lineToCommentMap
}

func (f *fileCommentRetriever) CommentOnLineRelativeToPos(pos token.Pos, delta int) (string, token.Position) {
	// get file and line number of expression
	fileToken := f.fset.File(pos)
	astNodeLineNum := fileToken.Line(pos)

	// get map from lines to relevant comments for the file
	lineToCommentMap, ok := f.fileToLineToComment[fileToken]

	// if map does not exist, populate it for the file (lazily load files)
	if !ok {
		lineToCommentMap = CreateLineToCommentMap(fileToken, f.fileTokenToASTFile[fileToken], f.requiredCommentPrefix)
		f.fileToLineToComment[fileToken] = lineToCommentMap
	}
	lineNum := astNodeLineNum + delta
	comment, ok := lineToCommentMap[lineNum]
	if !ok {
		return "", token.Position{}
	}
	return comment.Text, fileToken.Position(comment.Slash)
}
