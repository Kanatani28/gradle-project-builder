package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"./utils"
)

var (
	C      utils.Config
	tmpDir = "output"
	// 現在日付
	nowDate    = time.Now().Format("20060102")
	currentEnv string
	envList    []string
)

func init() {
	// config.ymlをロードする
	C = utils.LoadConfig()
	log.Printf(`# Settings:
			Date : %s
			BuildTarget : %s
			OutputDir : %s`,
		nowDate, C.BuildTarget, C.OutputDir)

	// ビルド対象の環境のリストを作成する
	envList = strings.Split(C.BuildTarget, ",")
	for i, s := range envList {
		envList[i] = strings.TrimSpace(s)
	}
}

func main() {

	log.Println("Start Main Function .")

	if err := os.MkdirAll(C.OutputDir, os.ModePerm); err != nil {
		log.Fatalln(err)
	}

	createWorkingDir()
	cloneProjects()
	for _, env := range envList {
		currentEnv = env

		setBuildEnv()
		buildProjects()
		moveModules()
	}
	removeProjects()

	log.Println("Finish Main Function .")
}

func createWorkingDir() {
	log.Println("# Making working directories start .")

	os.Chdir(C.OutputDir)

	// 現在日のディレクトリ作成
	os.Mkdir(nowDate, os.ModePerm)
	os.Chdir(nowDate)

	// tmpディレクトリ作成
	os.Mkdir(tmpDir, os.ModePerm)
	os.Chdir(tmpDir)
	for _, env := range envList {
		os.Mkdir(env, os.ModePerm)
		os.Chdir(env)
		for _, project := range C.Projects {
			if project.Ignore {
				continue
			}
			os.Mkdir(project.Name, os.ModePerm)
		}
		os.Chdir("..")
	}

	os.Chdir("..")

	log.Println("# Making working directories finished .")
}

func cloneProjects() {
	log.Println("# Clone repositoryies start .")
	for _, project := range C.Projects {

		log.Printf("## Start cloning %s .\n", project.Name)
		if project.Ignore {
			log.Printf("### %s is ignored .\n", project.Name)
			continue
		}
		// リポジトリからクローンする
		gitClone := exec.Command("git", "clone", project.RepositoryPath)
		gitClone.Run()
		log.Printf("## End cloning %s .\n", project.Name)
	}
	log.Println("# Cloning repositoryies finished .")

}

func buildProjects() {
	log.Println("# Building projects start .")
	for _, project := range C.Projects {
		log.Printf("## Start building %s .\n", project.Name)
		if project.Ignore {
			log.Printf("### %s is ignored .\n", project.Name)
			continue
		}
		os.Chdir(project.Name)

		buildCmd := exec.Command("gradlew", project.BuildType)
		buildCmd.Run()
		os.Chdir("..")
		log.Printf("## End building %s .\n", project.Name)
	}
	log.Println("# Build projects finished .")
}

// ブランチの切り替えと設定ファイルの書き換えを行う
func setBuildEnv() {

	log.Println("# Set build env start .")

	isV3 := ("v3" == currentEnv[len(currentEnv)-2:])

	for _, project := range C.Projects {
		if project.Ignore {
			continue
		}

		branch := "Phase2.5"
		if isV3 && project.Name != "facility-reservation-system" {
			branch = "Phase3"
		}
		log.Printf("##"+project.Name+" : branch is %s \n", branch)

		os.Chdir(project.Name)

		// ブランチ切替
		cmd := exec.Command("git", "checkout", branch)
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}

		if project.BuildType == "war" {
			// application.ymlの書き換え
			p := filepath.Join("src", "main", "resources", "config")
			p, _ = filepath.Abs(p)
			utils.ReplaceAppConfig(p, currentEnv)
		}

		os.Chdir("..")

	}

	log.Println("# Set build env finished .")

}

// ビルドしたモジュールをtmp/環境別ディレクトリへ移動する
func moveModules() {
	log.Println("# Move modules start .")
	for _, project := range C.Projects {
		if project.Ignore {
			continue
		}
		oldLocation := filepath.Join(".", project.Name, "build", "libs", project.OutputName)
		newLocation := filepath.Join(".", tmpDir, currentEnv, project.Name, project.OutputName)

		err := os.Rename(oldLocation, newLocation)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("# Move modules finished .")

}

// Cloneしてきたプロジェクトを削除する
func removeProjects() {
	log.Println("# Remove projects start .")
	for _, project := range C.Projects {
		err := os.RemoveAll(project.Name)
		if err != nil {
			log.Fatalln(err)
		}
	}
	log.Println("# Remove projects finished.")
}
