package user

import "encoding/json"

// Gender - 성별
type Gender int

const (
	GenderMale Gender = iota
	GenderFemale
	GenderOther
)

func (g *Gender) String() string {
	switch *g {
	case GenderMale:
		return "male"
	case GenderFemale:
		return "female"
	case GenderOther:
		return "other"
	default:
		return "unknown"
	}
}

func (g *Gender) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.String())
}

func (g *Gender) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*g = GenderFromString(s)
	return nil
}

func GenderFromString(s string) Gender {
	switch s {
	case "male":
		return GenderMale
	case "female":
		return GenderFemale
	case "other":
		return GenderOther
	default:
		return GenderMale
	}
}

// TravelPurpose - 여행 목적
type TravelPurpose int

const (
	TravelPurposeTourism     TravelPurpose = iota // 관광
	TravelPurposeBusiness                         // 비즈니스
	TravelPurposeBackpacking                      // 배낭여행
	TravelPurposeFoodTour                         // 맛집탐방
	TravelPurposeCulture                          // 문화체험
	TravelPurposeActivity                         // 액티비티
	TravelPurposeRelaxation                       // 휴양
)

func (tp *TravelPurpose) String() string {
	purposes := []string{
		"tourism", "business", "backpacking", "food_tour",
		"culture", "activity", "relaxation",
	}
	if int(*tp) < len(purposes) {
		return purposes[*tp]
	}
	return "unknown"
}

func (tp *TravelPurpose) MarshalJSON() ([]byte, error) {
	return json.Marshal(tp.String())
}

func (tp *TravelPurpose) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*tp = TravelPurposeFromString(s)
	return nil
}

func TravelPurposeFromString(s string) TravelPurpose {
	purposes := map[string]TravelPurpose{
		"tourism":     TravelPurposeTourism,
		"business":    TravelPurposeBusiness,
		"backpacking": TravelPurposeBackpacking,
		"food_tour":   TravelPurposeFoodTour,
		"culture":     TravelPurposeCulture,
		"activity":    TravelPurposeActivity,
		"relaxation":  TravelPurposeRelaxation,
	}
	if val, ok := purposes[s]; ok {
		return val
	}
	return TravelPurposeTourism
}

// TravelStyle - 여행 스타일
type TravelStyle int

const (
	TravelStylePlanned     TravelStyle = iota // 계획형
	TravelStyleSpontaneous                    // 즉흥형
	TravelStyleLuxury                         // 럭셔리
	TravelStyleBudget                         // 알뜰형
	TravelStyleAdventure                      // 모험형
	TravelStyleLeisurely                      // 여유형
)

func (ts *TravelStyle) String() string {
	styles := []string{
		"planned", "spontaneous", "luxury", "budget", "adventure", "leisurely",
	}
	if int(*ts) < len(styles) {
		return styles[*ts]
	}
	return "unknown"
}

func (ts *TravelStyle) MarshalJSON() ([]byte, error) {
	return json.Marshal(ts.String())
}

func (ts *TravelStyle) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*ts = TravelStyleFromString(s)
	return nil
}

func TravelStyleFromString(s string) TravelStyle {
	styles := map[string]TravelStyle{
		"planned":     TravelStylePlanned,
		"spontaneous": TravelStyleSpontaneous,
		"luxury":      TravelStyleLuxury,
		"budget":      TravelStyleBudget,
		"adventure":   TravelStyleAdventure,
		"leisurely":   TravelStyleLeisurely,
	}
	if val, ok := styles[s]; ok {
		return val
	}
	return TravelStylePlanned
}
