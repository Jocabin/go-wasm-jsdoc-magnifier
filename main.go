package main

import (
	"regexp"
	"strings"
	"syscall/js"
)

func main() {
	js.Global().Set("ExDetectComment", js.FuncOf(ExDetectComment))
	js.Global().Set("ExExtractComment", js.FuncOf(ExExtractComment))
	js.Global().Set("ExExtractAllComments", js.FuncOf(ExExtractAllComments))
	js.Global().Set("ExJoinMapJStoMarkdown", js.FuncOf(ExJoinMapJStoMarkdown))
	<-make(chan bool)
}

func ExDetectComment(this js.Value, args []js.Value) interface{} {
	return DetectComment(args[0].String())
}
func ExExtractComment(this js.Value, args []js.Value) interface{} {
	return ExtractComment(args[0].String())
}
func ExExtractAllComments(this js.Value, args []js.Value) interface{} {
	return ExtractAllComments(args[0].String())
}
func ExMapJStoMarkdown(this js.Value, args []js.Value) interface{} {
	var val []string = strings.Split(args[0].String(), "")
	return MapJStoMarkdown(val)
}
func ExJoinMapJStoMarkdown(this js.Value, args []js.Value) interface{} {
	return JoinMapJStoMarkdown(args[0].String())
}

// functions
func DetectComment(input string) bool {
	if strings.Contains(input, "/*") {
		return false
	}
	return true
}

func ExtractComment(input string) string {
	start_pattern := regexp.MustCompile(`/\/\*`)
	end_pattern := regexp.MustCompile(`/\*\/`)
	start := start_pattern.FindStringIndex(input)[0]
	end := end_pattern.FindStringIndex(input)[0]
	return `\n` + strings.Trim(input[start+2:end], "") + `\n`
}

func ExtractAllComments(input string) []string {
	pattern := regexp.MustCompile(`/\*[\s\S]*?\*/`)
	matches := pattern.FindAllString(input, -1)

	var comments []string
	if matches == nil {
		return comments
	}

	for _, item := range matches {
		cleaned := strings.ReplaceAll(item, "/*", "")
		cleaned = strings.ReplaceAll(cleaned, "*/", "")
		comments = append(comments, cleaned)
	}

	return comments
}

func MapJStoMarkdown(input []string) [][]string {
	var results [][]string

	for _, comment := range input {
		var el []string
		splittedComment := strings.Split(comment, "\n")
		title := strings.Join(splittedComment[1:2], "")
		splittedComment = append(splittedComment[:1], splittedComment[2:]...)
		finalString := "## " + title

		el = append(el, finalString)
		el = append(el, splittedComment...)
		results = append(results, el)
	}

	return results
}

func JoinMapJStoMarkdown(input string) string {
	comments := ExtractAllComments(input)
	processedComments := MapJStoMarkdown(comments)
	// fmt.Println(processedComments)

	var final_array []string

	for _, el := range processedComments {
		final_array = append(final_array, strings.Join(el, "\n"))
	}

	return strings.Join(final_array, "")
}
