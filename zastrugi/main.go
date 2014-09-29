package main

import (
    "bytes"
    "flag"
    "fmt"
    "io/ioutil"
    "encoding/json"
    "path/filepath"
    "strings"
    "os"
    "regexp"
    "github.com/coreos/go-etcd/etcd"
    "sync"
)

const ProjectConf = ".zastrugi"
const LDelim = "%%"
const RDelim = "%%"

// a type for our config file
type ConfigFile struct {
    Namespace string `json:namespace`
    SearchPaths []string `json:searchPaths`
}

// a type that represents a token within a file
type Token struct {
    startPos int
    endPos int
    key string
}

func main() {
    // command line flags
    var dataSourcePtr = flag.String("datasource", "", "the URL for the datasource")
    var rootPathPtr = flag.String("rootpath", "/.zastrugi", "the root path in the datasource")
    flag.Parse()

    // load the default json config file, which is '.zastrugi' in the current dir
    fmt.Println("Loading config:", ProjectConf)
    confFile, err := ioutil.ReadFile(ProjectConf)
    if (err != nil) {
        panic(err)
    }
    config := ConfigFile{}
    if err := json.Unmarshal(confFile, &config); err != nil {
        panic(err)
    }

    keyPrefix := *rootPathPtr + "/" + config.Namespace
    etcdClient := etcd.NewClient([]string{*dataSourcePtr})

    // spawn four goroutines for scanning files and replacing tokens
    tasks := make(chan string, 4)
    var wg sync.WaitGroup

    for i := 0; i < 4; i++ {
        wg.Add(1)

        // the goroutine takes filenames off of the tasks channel
        // and processes them
        go func() {
            for fileName := range tasks {
                fmt.Println("processing file:", fileName)
                processFile(etcdClient, keyPrefix, fileName)
            }

            wg.Done()
        }()
    }

    // walk the searchPaths to find file candidates for token replacement
    for _, path := range config.SearchPaths {
        // searchpaths can be globs
        matches, e := filepath.Glob(path)
        if (e != nil) {
            panic(e)
        }

        for _, fileName := range matches {
            // send each filename to the tasks channel
            tasks <- fileName
        }
    }

    // close the channel and wait for the workers to finish
    close(tasks)
    wg.Wait()
}

func processFile(etcdClient *etcd.Client, keyPrefix string, fileName string) {
    // read the file as a UTF-8 string
    target, e := ioutil.ReadFile(fileName)
    if (e != nil) {
        panic(e)
    }

    original := string(target)

    // scan for tokens - this is currently very naive
    var tokens []Token
    for rs, i := []rune(original), 0; i < len(rs); {
        if (strings.HasPrefix(original[i:], LDelim)) {
            i, tokens = scanToken(i, original, tokens)
        } else {
            i += 1
        }
    }

    if (len(tokens) > 0) {
        // print messages to the console if any of the tokens are invalid
        validateTokens(fileName, tokens)

        // search etcd for keys that match our tokens
        replacements := lookUpValues(etcdClient, keyPrefix, fileName, tokens)

        if (len(replacements) > 0) {
            replaceFile(fileName, original, tokens, replacements)
        }
    }
}

func scanToken(index int, input string, result []Token) (int, []Token) {
    lessToken := input[index + len(LDelim):]

    endPos := strings.Index(lessToken, RDelim)

    if (endPos > -1) {
        actualEndPos := index + endPos + len(LDelim) + len(RDelim)
        key := lessToken[:endPos]
        tok := Token{startPos: index, endPos: actualEndPos, key: key}
        result = append(result, tok)
        index += actualEndPos - index
        return index, result
    } else {
        index += len(LDelim)
    }

    return index, result
}

func validateTokens(fileName string, tokens []Token) bool {
    for _, token := range tokens {
        if match, err := regexp.MatchString("^[a-zA-Z0-9_.-]+$", token.key); err != nil || !match {
            fmt.Println(fmt.Sprintf("keys restricted to alphanumerics, underscores, hyphens, and periods (key '%v' found in file %v)", token.key, fileName))
            return false
        }
    }

    return true
}

func lookUpValues(etcdClient *etcd.Client, keyPrefix string, fileName string, tokens []Token) map[string]string {

    mapping := make(map[string]string)

    for _, token := range tokens {
        fullKey := keyPrefix + "/" + token.key
        rawResponse, err := etcdClient.RawGet(fullKey, true, false)

        if (rawResponse == nil) {
            panic(err)
        } else {
            if (rawResponse.StatusCode == 404) {
                fmt.Println(fmt.Sprintf("Key '%v' not found in file '%v'", token.key, fileName))
            } else {
                item, err := rawResponse.Unmarshal()

                if (err != nil) {
                    panic(err)
                }

                mapping[token.key] = item.Node.Value
            }
        }
    }

    return mapping
}

func replaceFile(fileName string, original string, tokens []Token, replacements map[string]string) {
    // build a buffer of bytes for the replacement file
    var resultBuffer bytes.Buffer
    for rs, i := []rune(original), 0; i < len(rs); {
        for _, token := range tokens {
            if i == token.startPos {
                replacement, present := replacements[token.key]

                // if a replacement was not in etcd, we just return the
                // token to the file
                if (present) {
                    resultBuffer.WriteString(replacement)
                    i += token.endPos - i
                    continue
                }
            }
        }

        // copy the current string if there is no token
        resultBuffer.WriteString(string(rs[i]))
        i += 1
    }

    file, err := os.OpenFile(fileName, os.O_RDWR, 0660);
    if (err != nil) {
        panic(err)
    }

    _, err = file.WriteString(resultBuffer.String())

    if (err != nil) {
        panic(err)
    }

    file.Close()

}
