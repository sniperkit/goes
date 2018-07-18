package model

import (
	"fmt"

	//// external
	"github.com/google/go-github/github"

	jsoniter "github.com/json-iterator/go"
	// bleve search
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzers/keyword_analyzer"
	"github.com/blevesearch/bleve/analysis/analyzers/simple_analyzer"
	"github.com/blevesearch/bleve/analysis/language/en"
	// logs
	"github.com/sirupsen/logrus"
)

var textFieldAnalyzer = "standard"

func init() {
	textFieldAnalyzer = "en"
}

// InitIndex initializes the search index at the specified path
func InitIndex(filepath string) (bleve.Index, error) {
	index, err := bleve.Open(filepath)
	// Doesn't yet exist (or error opening) so create a new one
	if err != nil {
		index, err = bleve.New(filepath, buildIndexMapping())
		if err != nil {
			log.WithError(err).WithFields(logrus.Fields{"section:": "model", "typology": "index", "step": "InitIndex"}).Warnf("%#s", err)
			return nil, err
		}
	}
	return index, nil
}

func OpenIndex(path string) bleve.Index {
	index, err := bleve.Open(path)
	if err == bleve.ErrorIndexPathDoesNotExist {
		//log.Printf("Creating new index...")
		log.WithFields(logrus.Fields{"section:": "model", "typology": "index", "step": "OpenIndex"}).Info("Creating new index...")
		// create a mapping
		indexMapping := buildIndexMapping()
		index, err = bleve.New(path, indexMapping)
		if err != nil {
			// log.Fatal(err)
			log.WithError(err).WithFields(logrus.Fields{"section:": "model", "typology": "index", "step": "OpenIndex"}).Fatal(err)
		}
	} else if err == nil {
		//log.Printf("Opening existing index...")
		log.WithError(err).WithFields(logrus.Fields{"section:": "model", "typology": "index", "step": "OpenIndex"}).Warn("Opening existing index...")
	} else {
		log.WithError(err).WithFields(logrus.Fields{"section:": "model", "typology": "index", "step": "OpenIndex"}).Fatal(err)
		//log.Fatal(err)
	}
	return index
}

/*
func ProcessUpdate(index bleve.Index, repo *github.Repository, path string) {
    log.Printf("updated: %s", path)
    rp := utils.relativePath(path)
    wiki, err := NewWikiFromFile(path)
    if err != nil {
        log.WithError(err).WithFields(logrus.Fields{"section:": "model", "typology": "index", "step": "ProcessUpdate"}).Warn(err)
        //log.Print(err)
    } else {
        doGitStuff(repo, rp, wiki)
        index.Index(rp, wiki)
    }
}

func ProcessDelete(index bleve.Index, repo *github.Repository, path string) {
    log.Printf("delete: %s", path)
    rp := utils.relativePath(path)
    err := index.Delete(rp)
    if err != nil {
        log.WithError(err).WithFields(logrus.Fields{"section:": "model", "typology": "index", "step": "ProcessUpdate"}).Warn(err)
        //log.Print(err)
    }
}
*/

func buildIndexMapping() *bleve.IndexMapping {

	simpleTextFieldMapping := bleve.NewTextFieldMapping()
	simpleTextFieldMapping.Analyzer = simple_analyzer.Name

	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = keyword_analyzer.Name

	starMapping := bleve.NewDocumentMapping()
	starMapping.AddFieldMappingsAt("Name", simpleTextFieldMapping)
	starMapping.AddFieldMappingsAt("FullName", simpleTextFieldMapping)
	starMapping.AddFieldMappingsAt("Description", englishTextFieldMapping)
	starMapping.AddFieldMappingsAt("Language", keywordFieldMapping)
	starMapping.AddFieldMappingsAt("Tags.Name", keywordFieldMapping)
	starMapping.AddFieldMappingsAt("Topics.Name", keywordFieldMapping)
	starMapping.AddFieldMappingsAt("Languages.Name", keywordFieldMapping)
	// starMapping.AddFieldMappingsAt("Readmes.Content", keywordFieldMapping)
	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("Star", starMapping)

	// indexMapping.AddDocumentMapping("Repo", starMapping)
	return indexMapping

}

func buildIndexMappingFromWiki() *bleve.IndexMapping {

	enTextFieldMapping := bleve.NewTextFieldMapping()
	enTextFieldMapping.Analyzer = textFieldAnalyzer

	storeFieldOnlyMapping := bleve.NewTextFieldMapping()
	storeFieldOnlyMapping.Index = false
	storeFieldOnlyMapping.IncludeTermVectors = false
	storeFieldOnlyMapping.IncludeInAll = false

	dateTimeMapping := bleve.NewDateTimeFieldMapping()

	wikiMapping := bleve.NewDocumentMapping()
	wikiMapping.AddFieldMappingsAt("name", enTextFieldMapping)
	wikiMapping.AddFieldMappingsAt("body", enTextFieldMapping)
	wikiMapping.AddFieldMappingsAt("modified_by", enTextFieldMapping)
	wikiMapping.AddFieldMappingsAt("modified_by_name", enTextFieldMapping)
	wikiMapping.AddFieldMappingsAt("modified_by_email", enTextFieldMapping)
	wikiMapping.AddFieldMappingsAt("modified_by_avatar", storeFieldOnlyMapping)
	wikiMapping.AddFieldMappingsAt("modified", dateTimeMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("wiki", wikiMapping)

	indexMapping.DefaultAnalyzer = textFieldAnalyzer

	return indexMapping
}

// https://github.com/dastergon/strgz/blob/master/lib/bleve.go
func ShowResults(results *bleve.SearchResult, index bleve.Index) {
	if len(results.Hits) < 1 {
		fmt.Println(results)
	}
	for _, val := range results.Hits {
		id := val.ID
		doc, err := index.Document(id)
		if err != nil {
			log.WithError(err).WithFields(logrus.Fields{"section:": "model", "typology": "index", "step": "ShowResults"}).Warnf("%#s", err)
			fmt.Println(err)
		}
		for _, field := range doc.Fields {
			repo := github.Repository{}
			jsoniter.Unmarshal(field.Value(), &repo)
			log.WithFields(logrus.Fields{"section:": "model", "typology": "index", "step": "ShowResults"}).Infof("%s - %s (%s)\n\t%s\n", *repo.Name, *repo.Description, *repo.Language, *repo.HTMLURL)
			fmt.Printf("%s - %s (%s)\n\t%s\n", *repo.Name, *repo.Description, *repo.Language, *repo.HTMLURL)
		}
	}
}

/*
func startWatching(path string, index bleve.Index, repo *github.Repository) *fsnotify.Watcher {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }

    // start a go routine to process events
    go func() {
        idleTimer := time.NewTimer(10 * time.Second)
        queuedEvents := make([]fsnotify.Event, 0)
        for {
            select {
            case ev := <-watcher.Events:
                queuedEvents = append(queuedEvents, ev)
                idleTimer.Reset(10 * time.Second)
            case err := <-watcher.Errors:
                log.Fatal(err)
            case <-idleTimer.C:
                for _, ev := range queuedEvents {
                    if pathMatch(ev.Name) {
                        switch ev.Op {
                        case fsnotify.Remove, fsnotify.Rename:
                            // delete the path
                            processDelete(index, repo, ev.Name)
                        case fsnotify.Create, fsnotify.Write:
                            // update the path
                            processUpdate(index, repo, ev.Name)
                        default:
                            // ignore
                        }
                    }
                }
                queuedEvents = make([]fsnotify.Event, 0)
                idleTimer.Reset(10 * time.Second)
            }
        }
    }()

    // now actually watch the path requested
    err = watcher.Add(path)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("watching '%s' for changes...", path)

    return watcher
}

func walkForIndexing(path string, index bleve.Index, repo *github.Repository) {
    dirEntries, err := ioutil.ReadDir(path)
    if err != nil {
        log.Fatal(err)
    }
    for _, dirEntry := range dirEntries {
        dirEntryPath := path + string(os.PathSeparator) + dirEntry.Name()
        if dirEntry.IsDir() {
            walkForIndexing(dirEntryPath, index, repo)
        } else if pathMatch(dirEntry.Name()) {
            processUpdate(index, repo, dirEntryPath)
        }
    }
}
*/
