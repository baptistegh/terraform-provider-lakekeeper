package sdk

// Wire values of the Lakekeeper management API discriminators used in the
// provider schemas. The generated client does not export them as constants.
const (
	// AssigneeTypeUser and AssigneeTypeRole are the principal kinds of a permission assignment.
	AssigneeTypeUser = "user"
	AssigneeTypeRole = "role"

	// DeleteProfileTypeHard and DeleteProfileTypeSoft are the tabular delete profile types.
	DeleteProfileTypeHard = "hard"
	DeleteProfileTypeSoft = "soft"

	// defaultGCSUniverseDomain is the universe_domain of public Google Cloud service account keys.
	defaultGCSUniverseDomain = "googleapis.com"
)
