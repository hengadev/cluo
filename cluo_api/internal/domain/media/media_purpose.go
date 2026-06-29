package domain

type RecordingPurpose string

const (
	RecordingPurposeGeneral          RecordingPurpose = "general"
	RecordingPurposeWitnessInterview RecordingPurpose = "witness_interview"
)

func (p RecordingPurpose) IsValid() bool {
	switch p {
	case RecordingPurposeGeneral, RecordingPurposeWitnessInterview:
		return true
	default:
		return false
	}
}
