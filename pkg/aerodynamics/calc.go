package aerodynamics

import (
	"errors"
	"math"
)

type DrivAerBoundsRequest struct {
	VehicleLength float64 `json:"vehicleLength" validate:"required,gt=0"`
	BlockageRatio float64 `json:"blockageRatio" validate:"required,gt=0,lt=1"`
}

type DrivAerBoundsResponse struct {
	DomainLength           float64 `json:"domainLength"`
	DomainWidth            float64 `json:"domainWidth"`
	DomainHeight           float64 `json:"domainHeight"`
	InletToNoseDistance    float64 `json:"inletToNoseDistance"`
	OutletFromTailDistance float64 `json:"outletFromTailDistance"`
	SideWallDistance       float64 `json:"sideWallDistance"`
	RoofDistance           float64 `json:"roofDistance"`
}

type TurbulenceResponse struct {
	ReynoldsNumber          float64 `json:"reynoldsNumber"`
	BoundaryLayerThickness  float64 `json:"boundaryLayerThickness"`
	TurbulentK              float64 `json:"k"`
	TurbulentEpsilon        float64 `json:"epsilon"`
	SpecificDissipationRate float64 `json:"omega"`
}

func CalcDrivAerBounds(req DrivAerBoundsRequest) DrivAerBoundsResponse {
	refWidth := req.VehicleLength * 0.42
	refHeight := req.VehicleLength * 0.30
	vehicleArea := refWidth * refHeight
	tunnelArea := vehicleArea / req.BlockageRatio
	domainHeight := math.Sqrt(tunnelArea / 1.6)
	domainWidth := domainHeight * 1.6

	inlet := 2.0 * req.VehicleLength
	outlet := 5.0 * req.VehicleLength
	domainLength := inlet + req.VehicleLength + outlet

	return DrivAerBoundsResponse{
		DomainLength:           domainLength,
		DomainWidth:            domainWidth,
		DomainHeight:           domainHeight,
		InletToNoseDistance:    inlet,
		OutletFromTailDistance: outlet,
		SideWallDistance:       (domainWidth - refWidth) / 2,
		RoofDistance:           domainHeight - refHeight,
	}
}

func CalcTurbulenceParams(inletVelocity, characteristicLength, nu, turbulenceIntensity float64) TurbulenceResponse {
	re := inletVelocity * characteristicLength / nu
	delta := 0.37 * characteristicLength / math.Pow(re, 0.2)
	k := 1.5 * math.Pow(inletVelocity*turbulenceIntensity, 2)
	cmu := 0.09
	lt := 0.07 * characteristicLength
	epsilon := math.Pow(cmu, 0.75) * math.Pow(k, 1.5) / lt
	omega := math.Sqrt(k) / (math.Pow(cmu, 0.25) * lt)

	return TurbulenceResponse{
		ReynoldsNumber:          re,
		BoundaryLayerThickness:  delta,
		TurbulentK:              k,
		TurbulentEpsilon:        epsilon,
		SpecificDissipationRate: omega,
	}
}


func (r DrivAerBoundsRequest) Validate() error {
	if r.VehicleLength <= 0 {
		return errors.New("vehicleLength must be > 0")
	}
	if r.BlockageRatio <= 0 || r.BlockageRatio >= 1 {
		return errors.New("blockageRatio must be between 0 and 1")
	}
	return nil
}
