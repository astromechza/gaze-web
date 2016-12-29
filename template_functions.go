package main

func templateFuncAdd(i int64, j int64) int64 {
	return i + j
}

func templateFuncIRange(count int64) []int64 {
	output := make([]int64, count)
	for i, _ := range output {
		output[i] = int64(i)
	}
	return output
}

func buildTemplateFuncsMap() map[string]interface{} {
	output := make(map[string]interface{})
	output["add"] = templateFuncAdd
	output["irange"] = templateFuncIRange
	return output
}
