// Code created from "fixed128.go.tmpl" - don't edit by hand

package fixed

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xmath/num"
)

var (
	// F128d6Max holds the maximum F128d6 value.
	F128d6Max = F128d6{data: num.MaxInt128}
	// F128d6Min holds the minimum F128d6 value.
	F128d6Min                = F128d6{data: num.MinInt128}
	multiplierF128d6BigInt   = new(big.Int).Exp(big.NewInt(10), big.NewInt(6), nil)
	multiplierF128d6BigFloat = new(big.Float).SetPrec(128).SetInt(multiplierF128d6BigInt)
	multiplierF128d6         = num.Int128FromBigInt(multiplierF128d6BigInt)
)

// F128d6 holds a fixed-point value that contains up to 6 decimal places.
// Values are truncated, not rounded. Values can be added and subtracted
// directly. For multiplication and division, the provided Mul() and Div()
// methods should be used.
type F128d6 struct {
	data num.Int128
}

// F128d6FromFloat64 creates a new F128d6 value from a float64.
func F128d6FromFloat64(value float64) F128d6 {
	f, _ := F128d6FromString(new(big.Float).SetPrec(128).SetFloat64(value).Text('f', 7)) //nolint:errcheck
	return f
}

// F128d6FromInt64 creates a new F128d6 value from an int64.
func F128d6FromInt64(value int64) F128d6 {
	return F128d6{data: num.Int128From64(value).Mul(multiplierF128d6)}
}

// F128d6FromString creates a new F128d6 value from a string.
func F128d6FromString(str string) (F128d6, error) {
	if str == "" {
		return F128d6{}, errs.New("empty string is not valid")
	}
	if strings.ContainsRune(str, 'E') {
		// Given a floating-point value with an exponent, which technically
		// isn't valid input, but we'll try to convert it anyway.
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return F128d6{}, err
		}
		return F128d6FromFloat64(f), nil
	}
	parts := strings.SplitN(str, ".", 2)
	var neg bool
	value := new(big.Int)
	fraction := new(big.Int)
	switch parts[0] {
	case "":
	case "-", "-0":
		neg = true
	default:
		if _, ok := value.SetString(parts[0], 10); !ok {
			return F128d6{}, errs.New("invalid value")
		}
		if value.Sign() < 0 {
			neg = true
			value.Neg(value)
		}
		value.Mul(value, multiplierF128d6BigInt)
	}
	if len(parts) > 1 {
		var buffer strings.Builder
		buffer.WriteString("1")
		buffer.WriteString(parts[1])
		for buffer.Len() < 6+1 {
			buffer.WriteString("0")
		}
		frac := buffer.String()
		if len(frac) > 6+1 {
			frac = frac[:6+1]
		}
		if _, ok := fraction.SetString(frac, 10); !ok {
			return F128d6{}, errs.New("invalid value")
		}
		value.Add(value, fraction).Sub(value, multiplierF128d6BigInt)
	}
	if neg {
		value.Neg(value)
	}
	return F128d6{data: num.Int128FromBigInt(value)}, nil
}

// F128d6FromStringForced creates a new F128d6 value from a string.
func F128d6FromStringForced(str string) F128d6 {
	f, _ := F128d6FromString(str) //nolint:errcheck
	return f
}

// Add adds this value to the passed-in value, returning a new value.
func (f F128d6) Add(value F128d6) F128d6 {
	return F128d6{data: f.data.Add(value.data)}
}

// Sub subtracts the passed-in value from this value, returning a new value.
func (f F128d6) Sub(value F128d6) F128d6 {
	return F128d6{data: f.data.Sub(value.data)}
}

// Mul multiplies this value by the passed-in value, returning a new value.
func (f F128d6) Mul(value F128d6) F128d6 {
	return F128d6{data: f.data.Mul(value.data).Div(multiplierF128d6)}
}

// Div divides this value by the passed-in value, returning a new value.
func (f F128d6) Div(value F128d6) F128d6 {
	return F128d6{data: f.data.Mul(multiplierF128d6).Div(value.data)}
}

// Trunc returns a new value which has everything to the right of the decimal
// place truncated.
func (f F128d6) Trunc() F128d6 {
	return F128d6{data: f.data.Div(multiplierF128d6).Mul(multiplierF128d6)}
}

// Int64 returns the truncated equivalent integer to this value.
func (f F128d6) Int64() int64 {
	return f.data.Div(multiplierF128d6).AsInt64()
}

// Float64 returns the floating-point equivalent to this value.
func (f F128d6) Float64() float64 {
	f64, _ := new(big.Float).SetPrec(128).Quo(f.data.AsBigFloat(), multiplierF128d6BigFloat).Float64()
	return f64
}

// Comma returns the same as String(), but with commas for values of 1000 and
// greater.
func (f F128d6) Comma() string {
	var istr string
	integer := f.data.Div(multiplierF128d6)
	if integer.IsInt64() {
		istr = humanize.Comma(integer.AsInt64())
	} else {
		istr = humanize.BigComma(integer.AsBigInt())
	}
	fraction := f.data.Sub(integer.Mul(multiplierF128d6))
	if fraction.IsZero() {
		return istr
	}
	if fraction.Sign() < 0 {
		fraction = fraction.Neg()
	}
	fstr := fraction.Add(multiplierF128d6).String()
	for i := len(fstr) - 1; i > 0; i-- {
		if fstr[i] != '0' {
			fstr = fstr[1 : i+1]
			break
		}
	}
	var neg string
	if integer.IsZero() && f.data.Sign() < 0 {
		neg = "-"
	} else {
		neg = ""
	}
	return fmt.Sprintf("%s%s.%s", neg, istr, fstr)
}

func (f F128d6) String() string {
	integer := f.data.Div(multiplierF128d6)
	istr := integer.String()
	fraction := f.data.Sub(integer.Mul(multiplierF128d6))
	if fraction.IsZero() {
		return istr
	}
	if fraction.Sign() < 0 {
		fraction = fraction.Neg()
	}
	fstr := fraction.Add(multiplierF128d6).String()
	for i := len(fstr) - 1; i > 0; i-- {
		if fstr[i] != '0' {
			fstr = fstr[1 : i+1]
			break
		}
	}
	var neg string
	if integer.IsZero() && f.data.Sign() < 0 {
		neg = "-"
	} else {
		neg = ""
	}
	return fmt.Sprintf("%s%s.%s", neg, istr, fstr)
}

// MarshalText implements the encoding.TextMarshaler interface.
func (f *F128d6) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (f *F128d6) UnmarshalText(text []byte) error {
	f1, err := F128d6FromString(string(text))
	if err != nil {
		return err
	}
	*f = f1
	return nil
}

// MarshalJSON implements the json.Marshaler interface. Note that this
// intentionally generates a string to ensure the correct value is retained.
func (f *F128d6) MarshalJSON() ([]byte, error) {
	return []byte(`"` + f.String() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (f *F128d6) UnmarshalJSON(data []byte) error {
	if data[0] == '"' {
		if data[len(data)-1] != '"' {
			return fmt.Errorf("invalid JSON %q", string(data))
		}
		data = data[1 : len(data)-1]
	}
	v, err := F128d6FromString(string(data))
	if err != nil {
		return err
	}
	*f = v
	return nil
}

// MarshalYAML implements the yaml.Marshaler interface. Note that this
// intentionally generates a string to ensure the correct value is retained.
func (f F128d6) MarshalYAML() (interface{}, error) {
	return f.String(), nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (f *F128d6) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return errs.Wrap(err)
	}
	f1, err := F128d6FromString(str)
	if err != nil {
		return err
	}
	*f = f1
	return nil
}