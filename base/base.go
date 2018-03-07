package base

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"math/rand"
	"reflect"

	"gonum.org/v1/gonum/mat"
)

type float = float64

// DenseShuffle shuffles the rows of X and Y matrices
func DenseShuffle(X, Y *mat.Dense) {
	nSamples, nFeatures := X.Dims()
	_, nOutputs := Y.Dims()
	Xrowi := make([]float64, nFeatures, nFeatures)
	Yrowi := make([]float64, nOutputs, nOutputs)
	for i := nSamples - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		copy(Xrowi, X.RawRowView(i))
		X.SetRow(i, X.RawRowView(j))
		X.SetRow(j, Xrowi)
		copy(Yrowi, Y.RawRowView(i))
		Y.SetRow(i, Y.RawRowView(j))
		Y.SetRow(j, Yrowi)
	}
}

// DenseSigmoid put emelent-wise sigmoid of X into dst
func DenseSigmoid(dst *mat.Dense, X mat.Matrix) *mat.Dense {
	if dst == nil {
		r, c := X.Dims()
		dst = mat.NewDense(r, c, nil)
	}
	dst.Apply(func(i int, j int, v float64) float64 {
		return 1. / (1. + math.Exp(-v))
	}, X)
	return dst
}

func unused(...interface{}) {}

// MatStr return a string from a mat.Matrix
func MatStr(Xs ...mat.Matrix) string {
	if len(Xs) == 0 {
		return ""
	}
	nSamples, nFeatures := Xs[0].Dims()
	b := bytes.NewBuffer(nil)

	for i := 0; i < nSamples; i++ {
		for imat, X := range Xs {
			_, nFeatures = X.Dims()
			for j := 0; j < nFeatures; j++ {
				io.WriteString(b, fmt.Sprintf("%g", X.At(i, j)))
				if j < nFeatures-1 || imat < len(Xs)-1 {
					io.WriteString(b, "\t")
				} else {
					io.WriteString(b, "\n")
				}
			}
		}
	}
	return b.String()
}

// MatColStr return the string for a matrix column
func MatColStr(X mat.Matrix, j int) string {
	nSamples, _ := X.Dims()
	var t = make([]float64, nSamples)
	mat.Col(t, j, X)
	return fmt.Sprint(t)
}

// MatRowStr returns the string for a matrix row
func MatRowStr(X mat.Matrix, i int) string {
	_, nFeatures := X.Dims()
	var t = make([]float64, nFeatures)
	mat.Row(t, i, X)
	return fmt.Sprint(t)
}

// CopyStruct create an new *struct with copied fields using reflection. it's not a deep copy.
func CopyStruct(m interface{}) interface{} {

	mstruct := reflect.ValueOf(m)
	if mstruct.Kind() == reflect.Ptr {
		mstruct = mstruct.Elem()
	}
	m2 := reflect.New(mstruct.Type())
	for i := 0; i < mstruct.NumField(); i++ {
		c := m2.Elem().Type().Field(i).Name[0]
		if m2.Elem().Field(i).CanSet() && c >= 'A' && c <= 'Z' {
			m2.Elem().Field(i).Set(mstruct.Field(i))
		}
	}
	return m2.Interface()
}
