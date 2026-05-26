package foam

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"
)

type DictRequest struct {
	Solver        string  `json:"solver" validate:"required,oneof=simpleFoam pisoFoam rhoCentralFoam"`
	DeltaT        float64 `json:"deltaT" validate:"required,gt=0"`
	EndTime       float64 `json:"endTime" validate:"required,gt=0"`
	Length        float64 `json:"length" validate:"required,gt=0"`
	Width         float64 `json:"width" validate:"required,gt=0"`
	Height        float64 `json:"height" validate:"required,gt=0"`
	XCells        int     `json:"xCells" validate:"required,gte=1"`
	YCells        int     `json:"yCells" validate:"required,gte=1"`
	ZCells        int     `json:"zCells" validate:"required,gte=1"`
	SimpleGrading string  `json:"simpleGrading" validate:"required"`
}

type DictBundle struct {
	BlockMeshDict string `json:"blockMeshDict"`
	ControlDict   string `json:"controlDict"`
	FvSchemes     string `json:"fvSchemes"`
}

type templateData struct {
	DictRequest
}

func Generate(req DictRequest) (DictBundle, error) {
	data := templateData{DictRequest: req}

	blockMeshDict, err := executeTemplate(blockMeshTemplate, data)
	if err != nil {
		return DictBundle{}, fmt.Errorf("failed to render blockMeshDict: %w", err)
	}
	controlDict, err := executeTemplate(controlDictTemplate, data)
	if err != nil {
		return DictBundle{}, fmt.Errorf("failed to render controlDict: %w", err)
	}
	fvSchemes, err := executeTemplate(fvSchemesTemplate, data)
	if err != nil {
		return DictBundle{}, fmt.Errorf("failed to render fvSchemes: %w", err)
	}

	return DictBundle{BlockMeshDict: blockMeshDict, ControlDict: controlDict, FvSchemes: fvSchemes}, nil
}

func executeTemplate(tmpl string, data templateData) (string, error) {
	t, err := template.New("dict").Parse(tmpl)
	if err != nil {
		return "", err
	}
	var out bytes.Buffer
	if err := t.Execute(&out, data); err != nil {
		return "", err
	}
	return out.String(), nil
}

const blockMeshTemplate = `FoamFile
{
    version     2.0;
    format      ascii;
    class       dictionary;
    object      blockMeshDict;
}

convertToMeters 1;

vertices
(
    (0 0 0)
    ({{.Length}} 0 0)
    ({{.Length}} {{.Width}} 0)
    (0 {{.Width}} 0)
    (0 0 {{.Height}})
    ({{.Length}} 0 {{.Height}})
    ({{.Length}} {{.Width}} {{.Height}})
    (0 {{.Width}} {{.Height}})
);

blocks
(
    hex (0 1 2 3 4 5 6 7) ({{.XCells}} {{.YCells}} {{.ZCells}}) simpleGrading ({{.SimpleGrading}})
);
`

const controlDictTemplate = `FoamFile
{
    version     2.0;
    format      ascii;
    class       dictionary;
    object      controlDict;
}

application     {{.Solver}};
startFrom       startTime;
startTime       0;
stopAt          endTime;
endTime         {{.EndTime}};
deltaT          {{.DeltaT}};
writeControl    timeStep;
writeInterval   100;
`

const fvSchemesTemplate = `FoamFile
{
    version     2.0;
    format      ascii;
    class       dictionary;
    object      fvSchemes;
}

ddtSchemes
{
    default         Euler;
}

gradSchemes
{
    default         Gauss linear;
}

divSchemes
{
    default         none;
    div(phi,U)      Gauss upwind;
}
`


func (r DictRequest) Validate() error {
	if r.DeltaT <= 0 || r.EndTime <= 0 || r.Length <= 0 || r.Width <= 0 || r.Height <= 0 {
		return errors.New("all dimensions and times must be positive")
	}
	if r.XCells < 1 || r.YCells < 1 || r.ZCells < 1 {
		return errors.New("cell counts must be at least 1")
	}
	switch r.Solver {
	case "simpleFoam", "pisoFoam", "rhoCentralFoam":
	default:
		return errors.New("solver must be one of simpleFoam, pisoFoam, rhoCentralFoam")
	}
	if r.SimpleGrading == "" {
		return errors.New("simpleGrading is required")
	}
	return nil
}
