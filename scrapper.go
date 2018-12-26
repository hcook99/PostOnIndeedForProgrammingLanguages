package main

import(
	"net/http"
	"io/ioutil"
	"strings"
	"sync"
	"os"
)

var languagesInArea = make(map[string]int)

func getJobTags(location string,nextPage int) []string {

	query := "https://www.indeed.com/jobs?q=Software%20Engineer&l="+location+"&limit=50"

	if nextPage>0{
		query = query + "&start=" + string(nextPage)
	}


	resp, _ := http.Get(query)

	bytes, _ := ioutil.ReadAll(resp.Body)

	htmlStore := string(bytes)
	indexOfVarJobs := strings.Index(htmlStore, "var jobmap = {};")
	testIfThere := strings.Index(htmlStore, "'};\n</script>")	
	if indexOfVarJobs == -1{
		return []string{}
	}
	htmlStore = htmlStore[indexOfVarJobs:testIfThere+3]
	arrayOfValues := make([]string, 50)
	indexOfJK := strings.Index(htmlStore, "{jk:'")
	indexEndJK := strings.Index(htmlStore, "',efccid:")
	indexOfDeleteAfter := strings.Index(htmlStore, "'};")
	
	i := 0
	for strings.Contains(htmlStore, "'};") {
		arrayOfValues[i] = htmlStore[indexOfJK+5:indexEndJK]
		i++
		htmlStore = htmlStore[indexOfDeleteAfter+3:]
		indexOfJK = strings.Index(htmlStore, "{jk:'")
		indexEndJK = strings.Index(htmlStore, "',efccid:")
		indexOfDeleteAfter = strings.Index(htmlStore, "'};")
	}

	resp.Body.Close()

	return arrayOfValues
}

func getJobDescriptions(tag string) string {
	websiteUrl := "https://www.indeed.com/viewjob?jk=" + tag
	resp, _ := http.Get(websiteUrl)
	byteT, _ := ioutil.ReadAll(resp.Body)
	bytesToString := string(byteT)

	beginIndex := strings.Index(bytesToString, "<div class=\"jobsearch-JobComponent-description icl-u-xs-mt--md\">")
	endI := strings.LastIndex(bytesToString,"<div class=\"jobsearch-JobDescriptionTab-content\">")
	bytesToString = bytesToString[beginIndex:endI]
	r := strings.NewReplacer("<li>", "", "</li>", "")
	bytesToString = r.Replace(bytesToString)
	resp.Body.Close()

	return bytesToString
}

func listOfEveryLanguageToConsider() []string{
	bytes,_ := os.Open("static/ProgrammingLanguages.txt")
	sliceOfLanguage,_ := ioutil.ReadAll(bytes)
	bytesToString := string(sliceOfLanguage)

	everyProgrammingLanguages := strings.Split(bytesToString, "\n")

	return everyProgrammingLanguages
}

func mapOfProgrammingLanguagesInArea(listOfProgrammingLanguages []string, webPagesToAnalyze []string) map[string]int{
	mapOfLangs := make(map[string]int)


	for _,valuePages := range webPagesToAnalyze {
		for _,valueLang := range listOfProgrammingLanguages{
			if(len(valueLang)==1){
				valueLang = " "+ valueLang+" "
			}
			if(strings.Contains(valuePages, valueLang)){
				mapOfLangs[valueLang] = mapOfLangs[valueLang]+1
				continue
			}
			continue
		}
	}

	//fmt.Println(elapsed)
	return mapOfLangs
}

func getMap(locationOfSearch string) map[string]int{

	replac := strings.NewReplacer(" ", "+")
	locationOfSearch = replac.Replace(locationOfSearch)
	jobTags := []string{}
	var everyLanguage []string
	var jobsDescriptions []string
	mapOfLanguagesInArea := make(map[string]int)

	operationDone := make(chan bool)
	go func(){
		jobTags = append(jobTags, getJobTags(locationOfSearch,0)...)
		jobTags = append(jobTags, getJobTags(locationOfSearch,50)...)
		jobTags = append(jobTags, getJobTags(locationOfSearch,100)...)
		everyLanguage = listOfEveryLanguageToConsider()
		operationDone <- true
	}()
	<-operationDone

	if len(jobTags) == 0{
		return map[string]int{}
	}

	var wg sync.WaitGroup
	wg.Add(150)
	for _,value := range jobTags {
		go func(value string){
			defer wg.Done()
			jobsDescriptions = append(jobsDescriptions, getJobDescriptions(value))
		}(value)
	}
	wg.Wait()



	mapOfLanguagesInArea = mapOfProgrammingLanguagesInArea(everyLanguage,jobsDescriptions)
	return mapOfLanguagesInArea
}