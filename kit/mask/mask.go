package mask

import (
	"encoding/json"

	"github.com/showa-93/go-mask"
)

// Struct takes a struct value and a list of field names (optional).
// It masks the values of the specified fields in the JSON with a predefined mask.
// The function returns the masked struct as a byte slice or an error if any.
// We encourage you to use the mask tag for readability instead of the optional params.
func Struct(v any, params ...string) (any, error) {
	masker := mask.NewMasker()
	masker.RegisterMaskStringFunc(mask.MaskTypeFilled, masker.MaskFilledString)
	masker.RegisterMaskStringFunc(mask.MaskTypeFixed, masker.MaskFixedString)
	for _, p := range params {
		masker.RegisterMaskField(p, "filled")
		masker.RegisterMaskField(p, "fixed")
	}

	masked, err := masker.Mask(v)
	if err != nil {
		return nil, err
	}

	return masked, nil
}

// StructToByte takes a struct value and a list of field names (optional).
// It masks the values of the specified fields in the JSON with a predefined mask.
// The function returns the masked struct as a byte slice or an error if any.
// We encourage you to use the mask tag for readability instead of the optional params.
func StructToByte(v any, params ...string) ([]byte, error) {
	masker := mask.NewMasker()
	masker.RegisterMaskStringFunc(mask.MaskTypeFilled, masker.MaskFilledString)
	masker.RegisterMaskStringFunc(mask.MaskTypeFixed, masker.MaskFixedString)
	for _, p := range params {
		masker.RegisterMaskField(p, "filled")
		masker.RegisterMaskField(p, "fixed")
	}

	masked, err := masker.Mask(v)
	if err != nil {
		return nil, err
	}

	mv, err := json.Marshal(masked)
	if err != nil {
		return nil, err
	}

	return mv, nil
}

// MaskJSONBytes takes a JSON byte slice and a list of field names.
// It masks the values of the specified fields in the JSON with a predefined mask.
// The function returns the masked JSON byte slice or an error if any.
// If you have the struct opt for the Struct function instead, and complement it using the mask tag.
func JSONBytes(data []byte, params ...string) ([]byte, error) {
	var m map[string]any

	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	masker := mask.NewMasker()
	masker.RegisterMaskStringFunc(mask.MaskTypeFilled, masker.MaskFilledString)
	masker.RegisterMaskStringFunc(mask.MaskTypeFixed, masker.MaskFixedString)
	for _, p := range params {
		masker.RegisterMaskField(p, "filled")
		masker.RegisterMaskField(p, "fixed")
	}

	masked, err := masker.Mask(m)
	if err != nil {
		return nil, err
	}

	mv, err := json.Marshal(masked)
	if err != nil {
		return nil, err
	}

	return mv, nil
}
