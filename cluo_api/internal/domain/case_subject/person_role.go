package case_subject

type PersonRole string

const (
	RoleVictim         PersonRole = "victim"
	RoleSuspect        PersonRole = "suspect"
	RoleWitness        PersonRole = "witness"
	RoleClaimant       PersonRole = "claimant"
	RoleRepresentative PersonRole = "representative"
)

func (r PersonRole) IsValid() bool {
	switch r {
	case RoleVictim,
		RoleSuspect,
		RoleWitness,
		RoleClaimant,
		RoleRepresentative:
		return true
	default:
		return false
	}
}
