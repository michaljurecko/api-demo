package characterbiz

import "google.golang.org/protobuf/types/known/fieldmaskpb"

func maskSliceToMap(mask *fieldmaskpb.FieldMask) map[string]bool {
	out := make(map[string]bool)
	if mask == nil {
		return out
	}

	for _, v := range mask.GetPaths() {
		out[v] = true
	}

	return out
}
