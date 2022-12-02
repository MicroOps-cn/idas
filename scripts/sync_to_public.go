/*
 Copyright © 2022 MicroOps-cn.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"text/template"
)

type syncer struct {
	DstType       string
	DstFile       string
	SourcePkgPath string
	SourceType    string
}
type module struct {
	Types []syncer
	Path  string
	Name  string
}

var (
	syncTypes    = map[string]map[string][]syncer{}
	expMatchType = regexp.MustCompile(`^type (\S+) \S+$`)
	outDirPrefix = "public_out-"
	tmpl         *template.Template
)

var minLineLength = len("//@sync-to-public")

func init() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("Working directory resolution failed.")
	}
	_ = os.Chdir(path.Dir(path.Dir(file)))
	tmpl = template.Must(template.ParseFiles("scripts/sync_to_public.tmpl"))
}

func main() {
	outDir, err := os.MkdirTemp(".", outDirPrefix)
	if err != nil {
		panic(fmt.Errorf("failed to create dir %s: %s", outDir, err))
	}
	defer os.RemoveAll(outDir)
	err = filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			switch path {
			case ".git", "api", "dist", "public", "public_out", "gogo_out", "":
				return filepath.SkipDir
			}
		} else {
			if filepath.Ext(path) == ".go" {
				f, err := os.Open(path)
				if err != nil {
					return fmt.Errorf("[ERROR]failed to open file %s: %s", path, err)
				}
				scanner := bufio.NewScanner(f)
				for scanner.Scan() {
					buf := strings.TrimSpace(string(scanner.Bytes()))
					if len(buf) >= minLineLength && buf[0:2] == "//" && strings.HasPrefix(strings.TrimSpace(buf[2:]), "@sync-to-public:") {
						syncOptionStr := strings.TrimPrefix(strings.TrimSpace(buf[2:]), "@sync-to-public:")
						syncOptionSplit := strings.Split(syncOptionStr, ":")
						syncOption := syncer{
							SourcePkgPath: filepath.Dir(path),
						}
						if len(syncOptionSplit) == 2 {
							syncOption.DstFile = syncOptionSplit[0]
							syncOption.DstType = syncOptionSplit[1]
						} else {
							return fmt.Errorf("unknown sync option: %s", path)
						}
						scanner.Scan()
						buf = string(scanner.Bytes())
						submatch := expMatchType.FindStringSubmatch(buf)
						if len(submatch) == 2 {
							syncOption.SourceType = submatch[1]
							if _, ok := syncTypes[syncOption.DstFile]; !ok {
								syncTypes[syncOption.DstFile] = map[string][]syncer{}
							}

							syncTypes[syncOption.DstFile][syncOption.SourcePkgPath] = append(syncTypes[syncOption.DstFile][syncOption.SourcePkgPath], syncOption)
						} else {
							return fmt.Errorf("[WARN]The file has annotation, but the corresponding type is not obtained %s", path)
						}
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	for disFile, moduleTypes := range syncTypes {
		func() {
			if stat, err := os.Stat(disFile); err != nil {
				if !os.IsNotExist(err) {
					fmt.Printf("[ERROR]Failed to get file state: %s: %s\n", disFile, err)
					return
				}
			} else if stat.IsDir() {
				fmt.Printf("[ERROR]%s is s directory\n", disFile)
				return
			} else if stat.Size() > 5 {
				f, err := os.Open(disFile)
				if err == nil {
					defer f.Close()
					scanner := bufio.NewScanner(f)
					scanner.Scan()
					if !strings.Contains(string(scanner.Bytes()), "DO NOT EDIT") {
						fmt.Printf("[WARN]The file %s may not be automatically generated or has been modified. Ignore this file.\n", disFile)
						return
					} else if err := os.Truncate(disFile, 0); err != nil {
						fmt.Printf("[WARN]failed to truncate file %s. Ignore this file.\n", disFile)
						return
					}

				} else if !os.IsNotExist(err) {
					fmt.Printf("[ERROR]Failed to open file: %s: %s\n", disFile, err)
					return
				}
				f.Close()
			}
			if err = os.WriteFile(disFile, []byte("// Code generated by idas. DO NOT EDIT.\n\n"), 0644); err != nil {
				fmt.Printf("[ERROR]Failed to write file: %s: %s\n", disFile, err)
			}
			filename := path.Join(outDir, path.Base(disFile)+".main.go")
			outFile, err := os.OpenFile(filename, os.O_APPEND|os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
			if err != nil {
				panic(fmt.Errorf("failed to open file %s: %s", outDir, err))
			}
			defer outFile.Close()
			var modules []module
			for modulePath, moduleType := range moduleTypes {
				modules = append(modules, module{
					Path:  modulePath,
					Name:  strings.Replace(modulePath, "/", "_", -1),
					Types: moduleType,
				})
			}
			err = tmpl.Execute(outFile, map[string]interface{}{
				"filename": disFile,
				"modules":  modules,
			})
			if err != nil {
				panic(err)
			}
			outFile.Close()
			command := exec.Command("go", "run", filename)
			err = command.Start()
			if err != nil {
				panic(err)
			}
			err = command.Wait()
			if err != nil {
				panic(err)
			}
		}()
	}
}