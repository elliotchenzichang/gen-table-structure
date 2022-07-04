package table_gen

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"os/exec"
)

const InformationSchema = "information_schema"

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
}

type GenInfo struct {
	Schema       string
	Table        string
	ExportFolder string
	TemplatePath string
	FileName     string
	PackageName  string
	StructName   string
}

func NewGenerator(config *Config) (g *Generator, err error) {
	if err := checkParams(config); err != nil {
		return nil, err
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.Username, config.Password, config.Host, config.Port, InformationSchema)
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	g = &Generator{db: db}
	return g, nil
}

type Generator struct {
	db *sqlx.DB
}

func (g *Generator) Gen(config *GenInfo) (isSuccess bool, err error) {
	if err := checkGenInfo(config); err != nil {
		return false, err
	}
	tableInfos, err := g.executeQuery(config.Schema, config.Table)
	if err != nil {
		return false, err
	}
	templateMetaDatas, err := convertTableInfoToMeta(tableInfos)
	templateData := &TemplateData{
		PackageName: config.PackageName,
		StructName:  config.StructName,
		Meta:        templateMetaDatas,
	}
	var genPath string
	if config.ExportFolder == "" {
		genPath = config.FileName
	} else {
		genPath = fmt.Sprintf("%s/%s", config.ExportFolder, config.FileName)
	}
	isSuccess, err = genCodeByTemplate(genPath, config.TemplatePath, templateData)
	if err == nil && isSuccess {
		_, _ = exec.Command("go", "fmt", genPath).Output()
	}
	return isSuccess, err
}

//executeQuery executes query sql and return table infos for generate the struct
func (g *Generator) executeQuery(scheme, table string) (infos []*tableInfo, err error) {
	err = g.db.Select(&infos, " SELECT COLUMN_NAME, DATA_TYPE FROM COLUMNS WHERE TABLE_NAME=? and table_schema=?", table, scheme)
	return infos, err
}

func checkParams(config *Config) error {
	if config.Username == "" {
		return errors.New("miss Username")
	}
	if config.Host == "" {
		return errors.New("miss Host")
	}
	if config.Port == 0 {
		return errors.New("miss Port")
	}
	return nil
}

func checkGenInfo(config *GenInfo) error {
	if config.Table == "" {
		return errors.New("miss Table")
	}
	if config.TemplatePath == "" {
		return errors.New("miss TemplatePath")
	}
	if config.StructName == "" {
		return errors.New("miss StructName")
	}
	if config.FileName == "" {
		return errors.New("miss FileName")
	}
	if config.Schema == "" {
		return errors.New("miss Schema")
	}
	return nil
}
